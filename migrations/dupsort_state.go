package migrations

import (
	"encoding/binary"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/rlp"

	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/turbo/trie"
)

var dupSortHashState = Migration{
	Name: "dupsort_hash_state",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := db.(ethdb.BucketsMigrator).BucketExists(dbutils.CurrentStateBucketOld1); err != nil {
			return err
		} else if !exists {
			return OnLoadCommit(db, nil, true)
		}

		if err := db.(ethdb.BucketsMigrator).ClearBuckets(dbutils.CurrentStateBucket); err != nil {
			return err
		}
		extractFunc := func(k []byte, v []byte, next etl.ExtractNextFunc) error {
			return next(k, k, v)
		}

		if err := etl.Transform(
			db,
			dbutils.CurrentStateBucketOld1,
			dbutils.CurrentStateBucket,
			datadir,
			extractFunc,
			etl.IdentityLoadFunc,
			etl.TransformArgs{OnLoadCommit: OnLoadCommit},
		); err != nil {
			return err
		}

		if err := db.(ethdb.BucketsMigrator).DropBuckets(dbutils.CurrentStateBucketOld1); err != nil {
			return err
		}
		return nil
	},
}

var dupSortPlainState = Migration{
	Name: "dupsort_plain_state",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := db.(ethdb.BucketsMigrator).BucketExists(dbutils.PlainStateBucketOld1); err != nil {
			return err
		} else if !exists {
			return OnLoadCommit(db, nil, true)
		}

		if err := db.(ethdb.BucketsMigrator).ClearBuckets(dbutils.PlainStateBucket); err != nil {
			return err
		}
		extractFunc := func(k []byte, v []byte, next etl.ExtractNextFunc) error {
			return next(k, k, v)
		}

		if err := etl.Transform(
			db,
			dbutils.PlainStateBucketOld1,
			dbutils.PlainStateBucket,
			datadir,
			extractFunc,
			etl.IdentityLoadFunc,
			etl.TransformArgs{OnLoadCommit: OnLoadCommit},
		); err != nil {
			return err
		}

		if err := db.(ethdb.BucketsMigrator).DropBuckets(dbutils.PlainStateBucketOld1); err != nil {
			return err
		}
		return nil
	},
}

var dupSortIH = Migration{
	Name: "dupsort_intermediate_trie_hashes",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if err := db.(ethdb.BucketsMigrator).ClearBuckets(dbutils.IntermediateTrieHashBucket); err != nil {
			return err
		}
		buf := etl.NewSortableBuffer(etl.BufferOptimalSize)
		comparator := db.(ethdb.HasTx).Tx().Comparator(dbutils.IntermediateTrieHashBucket)
		buf.SetComparator(comparator)
		collector := etl.NewCollector(datadir, buf)
		hashCollector := func(keyHex []byte, hash []byte) error {
			if len(keyHex) == 0 {
				return nil
			}
			if len(keyHex) > trie.IHDupKeyLen {
				return collector.Collect(keyHex[:trie.IHDupKeyLen], append(keyHex[trie.IHDupKeyLen:], hash...))
			}
			return collector.Collect(keyHex, hash)
		}
		loader := trie.NewFlatDBTrieLoader(dbutils.CurrentStateBucket, dbutils.IntermediateTrieHashBucket)
		if err := loader.Reset(trie.NewRetainList(0), hashCollector /* HashCollector */, false); err != nil {
			return err
		}
		if _, err := loader.CalcTrieRoot(db, nil); err != nil {
			return err
		}
		if err := collector.Load(db, dbutils.IntermediateTrieHashBucket, etl.IdentityLoadFunc, etl.TransformArgs{
			Comparator: comparator,
		}); err != nil {
			return fmt.Errorf("gen ih stage: fail load data to bucket: %w", err)
		}

		// this Migration is empty, sync will regenerate IH bucket values automatically
		// alternative is - to copy whole stage here
		if err := db.(ethdb.BucketsMigrator).DropBuckets(dbutils.IntermediateTrieHashBucketOld1); err != nil {
			return err
		}
		return OnLoadCommit(db, nil, true)
	},
}

var logsIndex = Migration{
	Name: "logs_index_test",
	Up: func(tx ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if err := tx.(ethdb.BucketsMigrator).ClearBuckets(dbutils.Receipts); err != nil {
			return err
		}
		buf := etl.NewSortableBuffer(etl.BufferOptimalSize)
		comparator := tx.(ethdb.HasTx).Tx().Comparator(dbutils.Receipts)
		buf.SetComparator(comparator)
		collector := etl.NewCollector(datadir, buf)

		i := 0
		if err := tx.Walk(dbutils.BlockReceiptsPrefix, nil, 0, func(k, v []byte) (bool, error) {
			//blockHash := k[len(k)-32:]
			blockNum := k[:len(k)-32]

			storageReceipts := []*types.ReceiptForStorage{}
			if err := rlp.DecodeBytes(v, &storageReceipts); err != nil {
				return false, err
			}
			for txIdx, storageReceipt := range storageReceipts {
				txIndex := make([]byte, 4)
				binary.BigEndian.PutUint32(txIndex, uint32(txIdx))

				for logIdx, log := range storageReceipt.Logs {
					i++
					//fmt.Printf("kk: %x %x\n", common.CopyBytes(blockNum), storageReceipt.ContractAddress[:])

					newK := append(common.CopyBytes(blockNum), log.Address[:]...)

					if len(log.Topics) > 0 {
						for _, topic := range log.Topics {
							fmt.Printf("topic amount: %d\n", len(log.Topics))
							logIndex := make([]byte, 4)
							binary.BigEndian.PutUint32(logIndex, uint32(logIdx))

							//if _, ok := topics[topic]; !ok {
							//fmt.Printf("v: %x \n", topic, )
							//}
							newV := append(topic[:], txIndex...)
							newV = append(newV, logIndex...)
							//fmt.Printf("k: %x\n", newK)
							if err := collector.Collect(newK, newV); err != nil {
								return false, err
							}
						}
					} else {
						panic("No topics!!!\n")
						//fmt.Printf()
					}
				}
			}
			i++

			if i > 1000 {
				return false, nil
			}
			return true, nil
		}); err != nil {
			return err
		}

		panic(1)
		return OnLoadCommit(tx, nil, true)
	},
}
