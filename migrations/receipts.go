package migrations

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/core/rawdb"
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
	Up: func(db ethdb.DbWithPendingMutations, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		logEvery := time.NewTicker(30 * time.Second)
		defer logEvery.Stop()

		buf := make([]byte, 0, 100_000)
		if err := db.Walk(dbutils.BlockReceiptsPrefix, nil, 0, func(k, v []byte) (bool, error) {
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

			buf = buf[:0]
			if err := cbor.Marshal(&buf, storageReceipts); err != nil {
				return false, err
			}
			return true, db.Put(dbutils.BlockReceiptsPrefix, common.CopyBytes(k), common.CopyBytes(buf))
		}); err != nil {
			return err
		}

		return OnLoadCommit(db, nil, true)
	},
}

var receiptsTopicNormalForm = Migration{
	Name: "receipt_topic_normal_form_2",
	Up: func(tx ethdb.DbWithPendingMutations, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := tx.(ethdb.BucketsMigrator).BucketExists(dbutils.BlockReceiptsPrefixOld1); err != nil {
			return err
		} else if !exists {
			return OnLoadCommit(tx, nil, true)
		}

		if err := tx.(ethdb.BucketsMigrator).ClearBuckets(dbutils.BlockReceiptsPrefix, dbutils.LogTopic2Id, dbutils.LogId2Topic); err != nil {
			return err
		}

		logEvery := time.NewTicker(30 * time.Second)
		defer logEvery.Stop()

		ids, err := ethdb.Ids(tx)
		if err != nil {
			return err
		}
		ids.Topic = 0 // Important! to reset topicID counter

		type DeprecatedLog struct {
			Address common.Address `codec:"1"`
			Topics  []common.Hash  `codec:"2"`
			Data    []byte         `codec:"3"`
		}

		type DeprecatedReceipt struct {
			PostState         []byte           `codec:"1"`
			Status            uint64           `codec:"2"`
			CumulativeGasUsed uint64           `codec:"3"`
			Logs              []*DeprecatedLog `codec:"4"`
		}

		iBytes := make([]byte, 8)
		if err := tx.Walk(dbutils.BlockReceiptsPrefix, nil, 0, func(k, v []byte) (bool, error) {
			blockHashBytes := k[len(k)-32:]
			blockNum64Bytes := k[:len(k)-32]
			blockNum := binary.BigEndian.Uint64(blockNum64Bytes)
			canonicalHash := rawdb.ReadCanonicalHash(tx, blockNum)
			if !bytes.Equal(blockHashBytes, canonicalHash[:]) {
				return true, nil
			}

			// Decode the receipts by legacy data type
			deprecatedReceipts := []*DeprecatedReceipt{}
			if err := cbor.Unmarshal(&deprecatedReceipts, v); err != nil {
				return false, fmt.Errorf("invalid receipt array RLP: %w, blockNum=%d", err, blockNum)
			}

			// Encode by new data type
			storageReceipts := make(types.Receipts, len(deprecatedReceipts))
			for ri, r := range deprecatedReceipts {
				storageReceipts[ri] = &types.Receipt{
					PostState:         r.PostState,
					Status:            r.Status,
					CumulativeGasUsed: r.CumulativeGasUsed,
					Logs:              make([]*types.Log, len(r.Logs)),
				}

				for li, l := range r.Logs {
					storageReceipts[ri].Logs[li].TopicIds = make([]uint64, len(r.Logs[li].Topics))
					for ti, topic := range l.Topics {
						id, err := tx.Get(dbutils.LogTopic2Id, topic[:])
						if err != nil && !errors.Is(err, ethdb.ErrKeyNotFound) {
							return false, err
						}

						// create topic if not exists with topicID++
						if err != nil && errors.Is(err, ethdb.ErrKeyNotFound) {
							ids.Topic++
							binary.BigEndian.PutUint64(iBytes, ids.Topic)
							storageReceipts[ri].Logs[li].TopicIds[ti] = ids.Topic

							err = tx.Put(dbutils.LogTopic2Id, topic[:], common.CopyBytes(iBytes))
							if err != nil {
								return false, err
							}
							err = tx.Append(dbutils.LogId2Topic, common.CopyBytes(iBytes), topic[:])
							if err != nil {
								return false, err
							}
							continue
						}

						storageReceipts[ri].Logs[li].TopicIds[ti] = binary.BigEndian.Uint64(id)
					}
				}
			}

			newV, err := rlp.EncodeToBytes(storageReceipts)
			if err != nil {
				log.Crit("Failed to encode block receipts", "err", err)
			}

			if err := tx.Append(dbutils.BlockReceiptsPrefix, blockNum64Bytes, newV); err != nil {
				return false, err
			}

			select {
			default:
			case <-logEvery.C:
				sz, _ := tx.(ethdb.HasTx).Tx().BucketSize(dbutils.BlockReceiptsPrefix)
				sz1, _ := tx.(ethdb.HasTx).Tx().BucketSize(dbutils.LogTopic2Id)
				sz2, _ := tx.(ethdb.HasTx).Tx().BucketSize(dbutils.LogId2Topic)
				log.Info("Progress", "blockNum", blockNum,
					dbutils.BlockReceiptsPrefix, common.StorageSize(sz),
					dbutils.LogTopic2Id, common.StorageSize(sz1),
					dbutils.LogId2Topic, common.StorageSize(sz2),
				)
			}
			return true, nil
		}); err != nil {
			return err
		}

		//if err := tx.(ethdb.BucketsMigrator).DropBuckets(dbutils.BlockReceiptsPrefixOld1); err != nil {
		//	return err
		//}
		return OnLoadCommit(tx, nil, true)
	},
}
