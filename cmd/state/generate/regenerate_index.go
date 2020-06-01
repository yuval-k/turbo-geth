package generate

import (
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"time"
)

func RegenerateIndex(chaindata string, csBucket []byte) error {
	db, err := ethdb.NewBoltDatabase(chaindata)
	if err != nil {
		return err
	}
	ig := core.NewIndexGenerator(db)
	ig.ChangeSetBufSize = 256 * 1024 * 1024

	err = ig.DropIndex(dbutils.StorageHistoryBucket)
	if err != nil {
		return err
	}
	startTime := time.Now()
	fmt.Println("Index generation started", startTime)
	err = ig.GenerateIndex(0, csBucket)
	if err != nil {
		return err
	}
	fmt.Println("Index is successfully regenerated", "it took", time.Since(startTime))
	return nil
}

/**
account index
merge 24m28.76168323s
fill 18m15.03340184s
walk 5m53.023704948s
wri 5m46.932535463s
Index is successfully regenerated it took 54m24.421808829s


caltime 11m52.681896245s
fill 649580695310
walk 248915314990
wri 265573819339
Index is successfully regenerated it took 27m26.497011965s


caltime 8m47.152730189s
merge 26m52.724080005s
fill 1056139209166
walk 505363152516
wri 549504010353
Index is successfully regenerated it took 35m40.326613534s

add concurrency to merge
merge 17m37.827462955s
fill 1112782006332
walk 522476159190
wri 606755400206
caltime 9m14.132599549s

128mb buf
caltime 9m2.71814643s
merge 18m22.699848349s
fill 18m14.547520731s
walk 8m17.07871856s
wri 10m1.586380075s
Index is successfully regenerated it took 27m25.871274s

4 потока мерж
caltime 8m41.669206731s
merge 21m36.471219819s
fill 18m12.304514889s
walk 8m9.658442011s
wri 9m9.48262477s

storage
calctime 17.12
merge 51m43.100763556s
fill 22m27.188621768s
walk 11m44.751170304s
wri 19m25.136709603s
Index is successfully regenerated it took 1h8m56.707856368s

storage 4 threads
caltime 16m12.283974528s
merge 52m35.329102088s
fill 22m38.1297512s
walk 10m34.275977833s
wri 17m57.306465989s
Index is successfully regenerated it took 1h8m48.618407089s

*/
