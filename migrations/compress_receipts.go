package migrations

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/core/rawdb"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/ledgerwatch/turbo-geth/rlp"
	"time"
)

var receiptLeadingZeroes = Migration{
	Name: "receipt_leading_zeroes_5",
	Up: func(tx ethdb.DbWithPendingMutations, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := tx.(ethdb.BucketsMigrator).BucketExists(dbutils.BlockReceiptsPrefixOld1); err != nil {
			return err
		} else if !exists {
			return OnLoadCommit(tx, nil, true)
		}

		if err := tx.(ethdb.BucketsMigrator).ClearBuckets(dbutils.BlockReceiptsPrefix); err != nil {
			return err
		}

		logEvery := time.NewTicker(30 * time.Second)
		defer logEvery.Stop()

		ids, err := ethdb.Ids(tx)
		if err != nil {
			return err
		}
		iBytes := make([]byte, 4)
		if err := tx.Walk(dbutils.BlockReceiptsPrefixOld1, nil, 0, func(k, v []byte) (bool, error) {
			if err != nil {
				return false, err
			}

			blockHashBytes := k[len(k)-32:]
			blockNum64Bytes := k[:len(k)-32]
			blockNum := binary.BigEndian.Uint64(blockNum64Bytes)
			canonicalHash := rawdb.ReadCanonicalHash(tx, blockNum)
			if !bytes.Equal(blockHashBytes, canonicalHash[:]) {
				return true, nil
			}

			// Decode the receipts by legacy data type
			storageReceiptsLegacy := []*types.DeprecatedReceiptForStorage1{}
			if err := rlp.DecodeBytes(v, &storageReceiptsLegacy); err != nil {
				return false, fmt.Errorf("invalid receipt array RLP: %w, blockNum=%d", err, blockNum)
			}

			// Encode by new data type
			storageReceipts := make([]*types.ReceiptForStorage, len(storageReceiptsLegacy))
			for i, r := range storageReceiptsLegacy {
				storageReceipts[i] = (*types.ReceiptForStorage)(r)
			}

			for ri, r := range storageReceiptsLegacy {
				for li, l := range r.Logs {
					storageReceipts[ri].Logs[li].TopicIds = make([]uint32, len(storageReceipts[ri].Logs[li].Topics))
					for ti, topic := range l.Topics {
						id, err := tx.Get(dbutils.LogTopic2Id, topic[:])
						if err != nil {
							return false, err
						}
						if len(id) != 0 {
							storageReceipts[ri].Logs[li].TopicIds[ti] = binary.BigEndian.Uint32(id)
						} else {
							ids.Topic++
							binary.BigEndian.PutUint32(iBytes, ids.Topic)
							storageReceipts[ri].Logs[li].TopicIds[ti] = ids.Topic

							err := tx.Put(dbutils.LogId2Topic, topic[:], common.CopyBytes(iBytes))
							if err != nil {
								return false, err
							}
							err = tx.Put(dbutils.LogTopic2Id, common.CopyBytes(iBytes), topic[:])
							if err != nil {
								return false, err
							}
						}
					}
				}
			}

			newV, err := rlp.EncodeToBytes(storageReceipts)
			if err != nil {
				log.Crit("Failed to encode block receipts", "err", err)
			}

			if err := rlp.DecodeBytes(v, &storageReceipts); err != nil {
				return false, fmt.Errorf("invalid receipt array RLP: %w, blockNum=%d", err, blockNum)
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
