package migrations

import (
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

var testBodyLargeLMDBTx = Migration{
	Name: "test_body_large_lmdb_tx",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := db.(ethdb.NonTransactional).BucketExists(dbutils.BlockBodyPrefixOld1); err != nil {
			return err
		} else if !exists {
			return nil
		}

		if err := db.(ethdb.NonTransactional).ClearBuckets(dbutils.BlockBodyPrefix); err != nil {
			return err
		}

		extractFunc := func(k []byte, v []byte, next etl.ExtractNextFunc) error {
			return next(k, k, v)
		}

		if err := etl.Transform(
			db,
			dbutils.BlockBodyPrefixOld1,
			dbutils.BlockBodyPrefix,
			datadir,
			extractFunc,
			etl.IdentityLoadFunc,
			etl.TransformArgs{OnLoadCommit: OnLoadCommit},
		); err != nil {
			return err
		}

		if err := db.(ethdb.NonTransactional).DropBuckets(dbutils.BlockBodyPrefixOld1); err != nil {
			return err
		}
		return nil
	},
}
