package migrations

import (
	"encoding/binary"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/rlp"
	"github.com/ledgerwatch/turbo-geth/turbo/trie"
	"github.com/spaolacci/murmur3"
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
		if err := tx.(ethdb.BucketsMigrator).ClearBuckets(dbutils.ReceiptsIndex); err != nil {
			return err
		}
		comparator := tx.(ethdb.HasTx).Tx().Comparator(dbutils.ReceiptsIndex)

		murmur := murmur3.New32()
		var NoTopic = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
		var NoTopics = []common.Hash{NoTopic}

		extractFunc := func(k []byte, v []byte, next etl.ExtractNextFunc) error {
			//blockHash := k[len(k)-32:]
			blockNum := k[:len(k)-32]

			storageReceipts := []*types.ReceiptForStorage{}
			if err := rlp.DecodeBytes(v, &storageReceipts); err != nil {
				return err
			}

			for txIdx, storageReceipt := range storageReceipts {
				for logIdx, log := range storageReceipt.Logs {

					topics := log.Topics
					if len(log.Topics) == 0 {
						topics = NoTopics
					}

					logIndex := make([]byte, 4)
					binary.BigEndian.PutUint32(logIndex, uint32(logIdx))

					txIndex := make([]byte, 4)
					binary.BigEndian.PutUint32(txIndex, uint32(txIdx))

					duplicatedTopics := map[common.Hash]struct{}{}
					for _, topic := range topics {
						if _, ok := duplicatedTopics[topic]; ok { // skip topics duplicates
							continue
						}
						duplicatedTopics[topic] = struct{}{}

						if _, err := murmur.Write(topic[:]); err != nil {
							return err
						}

						topicId := make([]byte, 4)
						binary.BigEndian.PutUint32(topicId, murmur.Sum32())
						murmur.Reset()

						newK := log.Address[:]

						newV := common.CopyBytes(blockNum)
						newV = append(newV, topicId...)
						newV = append(newV, txIndex...)
						newV = append(newV, logIndex...)
						if err := next(k, newK, newV); err != nil {
							return err
						}
					}
				}
			}
			return nil
		}

		if err := etl.Transform(
			tx,
			dbutils.BlockReceiptsPrefix,
			dbutils.ReceiptsIndex,
			datadir,
			extractFunc,
			etl.IdentityLoadFunc,
			etl.TransformArgs{
				OnLoadCommit: OnLoadCommit,
				Comparator:   comparator,
			},
		); err != nil {
			return err
		}

		return nil
	},
}
