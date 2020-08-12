package migrations

import (
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

var dupsortHashState = Migration{
	Name: "dupsort_hash_state",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := db.(ethdb.NonTransactional).BucketExists(dbutils.CurrentStateBucketOld1); err != nil {
			return err
		} else if !exists {
			return OnLoadCommit(db, nil, true)
		}

		if err := db.(ethdb.NonTransactional).ClearBuckets(dbutils.CurrentStateBucket); err != nil {
			return err
		}
		i := 0
		extractFunc := func(k []byte, v []byte, next etl.ExtractNextFunc) error {
			i++
			if i%1_000_000 == 0 {
				fmt.Printf("len: %d\n", len(k))
			}
			return next(k, k, v)
		}

		if err := etl.Transform(
			db,
			dbutils.CurrentStateBucketOld1,
			dbutils.CurrentStateBucket,
			datadir,
			extractFunc,
			func(k []byte, value []byte, _ etl.State, next etl.LoadNextFunc) error {
				fmt.Printf("len2: %d\n", len(k))
				return next(k, k, value)
			},
			etl.TransformArgs{OnLoadCommit: OnLoadCommit},
		); err != nil {
			return err
		}

		if err := db.(ethdb.NonTransactional).DropBuckets(dbutils.CurrentStateBucketOld1); err != nil {
			return err
		}
		return nil
	},
}
