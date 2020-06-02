package downloader

import (
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

func spawnTxLookup(db ethdb.Database, dataDir string, quitCh chan struct{}) error {
	//return etl.Transform(
	//	db,
	//	dbutils.PlainContractCodeBucket,
	//	dbutils.ContractCodeBucket,
	//	dataDir,
	//	keyTransformExtractFunc(transformContractCodeKey),
	//	etl.IdentityLoadFunc,
	//	etl.TransformArgs{Quit: quit},
	//)
}

func unwindTxLookup(unwindPoint uint64, db ethdb.Database, quitCh chan struct{}) error {
	return nil
}
