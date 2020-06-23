package verify

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core/rawdb"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
)

func ValidateTxLookups(chaindata string) error {
	db := ethdb.MustOpen(chaindata)

	ch := make(chan os.Signal, 1)
	quitCh := make(chan struct{})
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		close(quitCh)
	}()
	t := time.Now()
	defer func() {
		log.Info("Validation ended", "it took", time.Since(t))
	}()
	var errorsNum uint64
	var blockNum uint64
	iterations := 0
	var interrupt bool
	// Validation Process
	blockBytes := big.NewInt(0)
	for !interrupt {
		if err := common.Stopped(quitCh); err != nil {
			return err
		}
		//fmt.Println(blockNum)
		blockHash := rawdb.ReadCanonicalHash(db, blockNum)
		body := rawdb.ReadBody(db, blockHash, blockNum)

		if body == nil {
			log.Error("Empty body", "blocknum", blockNum)
			break
		}
		//if blockNum > 100000 {
		//	break
		//}
		blockBytes.SetUint64(blockNum)
		bn := blockBytes.Bytes()

		for _, tx := range body.Transactions {
			val, err := db.Get(dbutils.TxLookupPrefix, tx.Hash().Bytes())
			iterations++
			if iterations%10000 == 0 {
				log.Info("Validated", "entries", iterations, "number", blockNum, "errors", atomic.LoadUint64(&errorsNum))

			}
			if !bytes.Equal(val, bn) {
				//fmt.Println(blockHash.String(), tx.Hash().String(), blockNum, big.NewInt(0).SetBytes(val).Uint64(), common.Bytes2Hex(val))
				if err != nil {
					log.Error("err!=nil","blockHash",blockHash.String(), "bn", blockNum, "bn",big.NewInt(0).SetBytes(val).Uint64())
					panic(err)
				}
				atomic.AddUint64(&errorsNum, 1)
				panic(fmt.Sprintf("Validation process failed(%d). Expected %b, got %b", iterations, bn, val))
			}
		}
		blockNum++
	}
	fmt.Println("errors num", errorsNum)
	return nil
}
