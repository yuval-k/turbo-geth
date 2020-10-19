package migrations

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"time"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/ethdb/cbor"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/ledgerwatch/turbo-geth/rlp"
)

var receiptsCborEncode = Migration{
	Name: "receipts_cbor_encode",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		logEvery := time.NewTicker(30 * time.Second)
		defer logEvery.Stop()
		collector, err1 := etl.NewCollectorFromFiles(datadir)
		if err1 != nil {
			return err1
		}
		if collector == nil {
			collector = etl.NewCriticalCollector(datadir, etl.NewSortableBuffer(etl.BufferOptimalSize))
			buf := bytes.NewBuffer(make([]byte, 0, 100_000))
			if err1 = db.Walk(dbutils.BlockReceiptsPrefix, nil, 0, func(k, v []byte) (bool, error) {
				blockNum := binary.BigEndian.Uint64(k[:8])
				select {
				default:
				case <-logEvery.C:
					var m runtime.MemStats
					runtime.ReadMemStats(&m)
					log.Info("Migration progress", "blockNum", blockNum, "alloc", common.StorageSize(m.Alloc), "sys", common.StorageSize(m.Sys))
				}

				// Convert the receipts from their storage form to their internal representation
				storageReceipts := []*types.ReceiptForStorage{}
				if err := rlp.DecodeBytes(v, &storageReceipts); err != nil {
					return false, fmt.Errorf("invalid receipt array RLP: %w, k=%x", err, k)
				}

				buf.Reset()
				if err := cbor.Marshal(buf, storageReceipts); err != nil {
					return false, err
				}
				if err := collector.Collect(k, buf.Bytes()); err != nil {
					return false, fmt.Errorf("collecting key %x: %w", k, err)
				}
				return true, nil
			}); err1 != nil {
				return err1
			}
		}
		if err := db.(ethdb.BucketsMigrator).ClearBuckets(dbutils.BlockReceiptsPrefix); err != nil {
			return fmt.Errorf("clearing the receipt bucket: %w", err)
		}
		// Commit clearing of the bucket - freelist should now be written to the database
		if err := OnLoadCommit(db, nil, false); err != nil {
			return fmt.Errorf("committing the removal of receipt table")
		}
		// Commit again
		if err := OnLoadCommit(db, nil, false); err != nil {
			return fmt.Errorf("committing again to create a stable view the removal of receipt table")
		}
		// Now transaction would have been re-opened, and we should be re-using the space
		if err := collector.Load(db, dbutils.BlockReceiptsPrefix, etl.IdentityLoadFunc, etl.TransformArgs{OnLoadCommit: OnLoadCommit}); err != nil {
			return fmt.Errorf("loading the transformed data back into the receipts table: %w", err)
		}
		return nil
	},
}

var receiptsOnePerTxEncode = Migration{
	Name: "receipts_one_per_tx2",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		logEvery := time.NewTicker(30 * time.Second)
		defer logEvery.Stop()

		type LegacyReceipt struct {
			// Consensus fields: These fields are defined by the Yellow Paper
			PostState         []byte       `codec:"1"`
			Status            uint64       `codec:"2"`
			CumulativeGasUsed uint64       `codec:"3"`
			Logs              []*types.Log `codec:"4"`
		}

		buf := bytes.NewBuffer(make([]byte, 0, 100_000))
		reader := bytes.NewReader(nil)
		i := 0

		collectorReceipts, err1 := etl.NewCollectorFromFiles(datadir)
		if err1 != nil {
			return err1
		}
		if collectorReceipts == nil && false {
			goto LoadPart
		}

		collectorReceipts = etl.NewCriticalCollector(datadir, etl.NewSortableBuffer(etl.BufferOptimalSize))
		if err := db.Walk(dbutils.BlockReceiptsPrefix, nil, 0, func(k, v []byte) (bool, error) {
			i++
			blockNum := binary.BigEndian.Uint64(k[:8])
			select {
			default:
			case <-logEvery.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				log.Info("Migration progress", "blockNum", blockNum, "alloc", common.StorageSize(m.Alloc), "sys", common.StorageSize(m.Sys))
			}

			// Convert the receipts from their storage form to their internal representation
			legacyReceipts := []*LegacyReceipt{}

			reader.Reset(v)
			if err := cbor.Unmarshal(&legacyReceipts, reader); err != nil {
				return false, err
			}

			// Convert the receipts from their storage form to their internal representation
			receipts := make(types.Receipts, len(legacyReceipts))
			for i := range legacyReceipts {
				receipts[i] = &types.Receipt{}
				receipts[i].PostState = legacyReceipts[i].PostState
				receipts[i].Status = legacyReceipts[i].Status
				receipts[i].CumulativeGasUsed = legacyReceipts[i].CumulativeGasUsed
				receipts[i].Logs = legacyReceipts[i].Logs
			}
			for txId, r := range receipts {
				newK := make([]byte, 8+4)
				copy(newK, k[:8])
				binary.BigEndian.PutUint32(newK[8:], uint32(txId))

				buf.Reset()
				if err := cbor.Marshal(buf, r); err != nil {
					return false, err
				}
				if err := collectorReceipts.Collect(newK, buf.Bytes()); err != nil {
					return false, fmt.Errorf("collecting key %x: %w", k, err)
				}
			}
			//if i > 100000 {
			//	panic(1)
			//}
			return true, nil
		}); err != nil {
			return err
		}

	LoadPart:
		//panic(2)
		if err := db.(ethdb.BucketsMigrator).ClearBuckets(dbutils.BlockReceiptsPrefix); err != nil {
			return fmt.Errorf("clearing the receipt bucket: %w", err)
		}
		// Commit clearing of the bucket - freelist should now be written to the database
		if err := OnLoadCommit(db, nil, true); err != nil {
			return fmt.Errorf("committing the removal of receipt table")
		}
		// Commit again
		if err := OnLoadCommit(db, nil, true); err != nil {
			return fmt.Errorf("committing again to create a stable view the removal of receipt table")
		}
		// Now transaction would have been re-opened, and we should be re-using the space
		if err := collectorReceipts.Load(db, dbutils.BlockReceiptsPrefix, etl.IdentityLoadFunc, etl.TransformArgs{OnLoadCommit: OnLoadCommit}); err != nil {
			return fmt.Errorf("loading the transformed data back into the receipts table: %w", err)
		}
		return nil
	},
}
