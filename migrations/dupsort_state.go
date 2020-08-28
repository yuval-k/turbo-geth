package migrations

import (
	"encoding/binary"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

var dupSortHashState = Migration{
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

		if err := db.(ethdb.NonTransactional).DropBuckets(dbutils.CurrentStateBucketOld1); err != nil {
			return err
		}
		return nil
	},
}

var dupSortPlainState = Migration{
	Name: "dupsort_plain_state",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := db.(ethdb.NonTransactional).BucketExists(dbutils.PlainStateBucketOld1); err != nil {
			return err
		} else if !exists {
			return OnLoadCommit(db, nil, true)
		}

		if err := db.(ethdb.NonTransactional).ClearBuckets(dbutils.PlainStateBucket); err != nil {
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

		if err := db.(ethdb.NonTransactional).DropBuckets(dbutils.PlainStateBucketOld1); err != nil {
			return err
		}
		return nil
	},
}

var dupSortHistoryAcc = Migration{
	Name: "dupsort_history_acc_test2",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if err := db.(ethdb.NonTransactional).ClearBuckets(dbutils.AccountsHistoryBucket2); err != nil {
			return err
		}
		i := 0
		extractFunc := func(k []byte, v []byte, next etl.ExtractNextFunc) error {
			i++
			fmt.Printf("1: %x, %x\n", k, v)
			fmt.Printf("2: %x, %d\n", k[:common.AddressLength], binary.BigEndian.Uint64(k[common.AddressLength:]))

			a, b, _ := dbutils.WrapHistoryIndex(v).Decode()
			fmt.Printf("3: %d, %t\n", a, b)
			if i > 100 {
				panic(1)
			}
			return next(k, k, v)
		}

		if err := etl.Transform(
			db,
			dbutils.AccountsHistoryBucket,
			dbutils.AccountsHistoryBucket2,
			datadir,
			extractFunc,
			etl.IdentityLoadFunc,
			etl.TransformArgs{OnLoadCommit: OnLoadCommit},
		); err != nil {
			return err
		}
		return nil
	},
}
