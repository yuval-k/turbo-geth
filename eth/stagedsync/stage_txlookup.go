package stagedsync

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/snappy"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/debug"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/ledgerwatch/turbo-geth/rlp"
	"math/big"
	"runtime"
	"sync"
	"time"
)

func spawnTxLookup(s *StageState, db ethdb.Database, dataDir string, quitCh chan struct{}) error {
	var blockNum uint64
	//var startKey []byte

	lastProcessedBlockNumber := s.BlockNumber
	if lastProcessedBlockNumber > 0 {
		blockNum = lastProcessedBlockNumber + 1
	}
	syncHeadNumber, err := s.ExecutionAt(db)
	if err == nil {
		//chunks = calculateTxLookupChunks(lastProcessedBlockNumber, syncHeadNumber, runtime.NumCPU()/2+1)
	}

	err = TxLookupTransform(db, blockNum, syncHeadNumber, quitCh, dataDir)
	if err != nil {
		return err
	}

	return s.DoneAndUpdate(db, syncHeadNumber)
}

var endHash common.Hash = common.Hash{
	255,255,255,255,255,
	255,255,255,255,255,
	255,255,255,255,255,
	255,255,255,255,255,
	255,255,255,255,255,
	255,255,255,255,255,
	255,255,
}
func TxLookupTransform(db ethdb.Database, start, end uint64, quitCh chan struct{}, datadir string) error {
	var chunks [][]byte
	const mpSize = 4000000
	mp:=make(map[uint64][]byte, mpSize)
	lastBlock :=start
	startBlock:=start
	var err error
	for {
		t:=time.Now()
		//fmt.Println("header prefix")
		err = db.Walk(dbutils.HeaderPrefix, dbutils.HeaderHashKey(startBlock),0, func(k []byte, v []byte) (bool, error) {
			if !dbutils.CheckCanonicalKey(k) {
				return true, nil
			}
			blocknum := binary.BigEndian.Uint64(k)
			//fmt.Println(blocknum,"canonical", common.Bytes2Hex(v))
			if blocknum > end {
				return false, nil
			}
			lastBlock = blocknum
			if v1,ok:=mp[blocknum]; ok {
				log.Info("duplicate block", "v", common.Bytes2Hex(v),"v1", common.Bytes2Hex(v1), "bn", blocknum)
			}
			mp[blocknum] = v

			if len(mp) >= mpSize {
				return false, nil
			}

			return true, nil
		})
		if err!=nil {
			return err
		}
		log.Warn("Canonical hash collection", "t", time.Since(t))
		if lastBlock>0 && lastBlock>startBlock {
			chunks = calculateTxLookupChunks(startBlock, lastBlock, runtime.NumCPU()/2+1)
		} else {
			chunks = nil
		}
		t=time.Now()
		log.Error("transform", "from", startBlock, "to", lastBlock, "chunks", len(chunks))
		p:=sync.Pool{}
		p2:=sync.Pool{}
		p2.New = func() interface{} {
			return new(big.Int)
		}
		p.New= func() interface{} {
			return new(types.Body)
		}
		err = etl.Transform(db, dbutils.BlockBodyPrefix, dbutils.TxLookupPrefix, datadir, func(k []byte, v []byte, next etl.ExtractNextFunc) error {
			blocknum := binary.BigEndian.Uint64(k)


			if vv:=mp[blocknum]; !bytes.Equal(vv,k[8:]) {
				log.Warn("not exist", "bn",blocknum, "hash", common.Bytes2Hex(k[8:]))
				return nil
			}

			body := p.Get().(*types.Body)
			if err := rlp.Decode(bytes.NewReader(v), body); err != nil {
				return err
			}
			//fmt.Println(blocknum, common.Bytes2Hex(k[8:]), len(body.Transactions))
			bnBig:=p2.Get().(*big.Int)
			blockNumBytes := bnBig.SetUint64(blocknum).Bytes()
			p2.Put(bnBig)

			for _, tx := range body.Transactions {
				//if blocknum >= 46169 && blocknum < 46171 {
				//	fmt.Println(blocknum, common.Bytes2Hex(k[8:]),  common.Bytes2Hex(tx.Hash().Bytes()), common.Bytes2Hex(blockNumBytes))
				//}
				if err := next(k, tx.Hash().Bytes(), blockNumBytes); err != nil {
					return err
				}
			}
			body.Transactions=nil
			p.Put(body)
			return nil
		}, etl.IdentityLoadFunc, etl.TransformArgs{
			Quit:            quitCh,
			ExtractStartKey: dbutils.BlockBodyKey(startBlock, common.Hash{}),
			ExtractEndKey:   dbutils.BlockBodyKey(lastBlock, endHash),
			Chunks:          chunks,
		})
		if err!=nil {
			return err
		}
		
		log.Warn("transform end", "t", time.Since(t))
		if lastBlock >= end {
			return nil
		}
		mp=make(map[uint64][]byte, mpSize)
		startBlock=lastBlock
	}


	/*
	INFO [06-19|22:02:41.212] Extraction finished                      it took=4m7.809185727s
	INFO [06-19|22:05:08.842] Collection finished                      it took=2m27.630319434s
	WARN [06-19|22:05:08.842] transform end                            t=6m35.439607754s
	WARN [06-19|22:06:44.305] Canonical hash collection                t=1m35.279625775s

		INFO [06-19|22:37:56.709] Extraction finished                      it took=31m12.404596901s
	INFO [06-19|22:56:57.609] Collection finished                      it took=19m0.899766227s
	WARN [06-19|22:56:57.649] transform end                            t=50m13.344053568s
	WARN [06-19|22:57:56.420] Canonical hash collection                t=58.554089693s

		INFO [06-19|23:15:15.058] Extraction finished                      it took=17m18.557741737s
	INFO [06-19|23:32:23.840] Collection finished                      it took=17m8.769042804s
	WARN [06-19|23:32:23.889] transform end                            t=34m27.445912361s
	INFO [06-19|23:32:23.916] TxLookup index is successfully regenerated it took=1h35m19.272436468s



	INFO [06-20|13:53:41.956] Extraction finished                      it took=7m45.417649175s
	INFO [06-20|13:54:42.572] Collection finished                      it took=1m0.616073575s
	WARN [06-20|13:54:42.572] transform end                            t=8m46.033842572s
	WARN [06-20|13:56:26.438] Canonical hash collection                t=1m43.832558145s
	INFO [06-20|14:11:33.291] Extraction finished                      it took=15m6.853041186s
	INFO [06-20|14:25:39.838] Collection finished                      it took=14m6.546253136s
	WARN [06-20|14:25:39.838] transform end                            t=29m13.39942723s
	WARN [06-20|14:26:40.231] Canonical hash collection                t=1m0.231635722s
	INFO [06-20|14:36:07.754] Extraction finished                      it took=9m27.522968432s

	INFO [06-20|14:42:38.312] Collection finished                      it took=6m30.558488024s
	WARN [06-20|14:42:38.312] transform end                            t=15m58.081557737s
	INFO [06-20|14:42:38.312] TxLookup index is successfully regenerated it took=58m12.29217694s
	INFO [06-20|14:42:38.488] bolt database closed                     bolt_db=/media/b00ris/nvme/tgstaged2/geth/chaindata

	*/
	//etl.Transform(db, dbutils.HeaderPrefix, dbutils.TxLookupPrefix, datadir, func(k []byte, v []byte, next etl.ExtractNextFunc) error {
	//	if !dbutils.CheckCanonicalKey(k) {
	//		return nil
	//	}
	//	blocknum := binary.BigEndian.Uint64(k)
	//	blockHash := common.BytesToHash(v)
	//	body := rawdb.ReadBody(db, blockHash, blocknum)
	//	if body == nil {
	//		return fmt.Errorf("tx lookup generation, empty block body %d, hash %x", blocknum, v)
	//	}
	//
	//	blockNumBytes := new(big.Int).SetUint64(blocknum).Bytes()
	//	for _, tx := range body.Transactions {
	//		if err := next(k, tx.Hash().Bytes(), blockNumBytes); err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//}, etl.IdentityLoadFunc, etl.TransformArgs{
	//	Quit:            quitCh,
	//	ExtractStartKey: startKey,
	//	ExtractEndKey:   endKey,
	//	Chunks:          chunks,
	//})
	return nil
}

func unwindTxLookup(u *UnwindState, db ethdb.Database, quitCh chan struct{}) error {
	var txsToRemove [][]byte
	// Remove lookup entries for all blocks above unwindPoint
	if err := db.Walk(dbutils.BlockBodyPrefix, dbutils.EncodeBlockNumber(u.UnwindPoint+1), 0, func(k, v []byte) (b bool, e error) {
		if err := common.Stopped(quitCh); err != nil {
			return false, err
		}
		data := v
		if debug.IsBlockCompressionEnabled() && len(data) > 0 {
			var err1 error
			data, err1 = snappy.Decode(nil, v)
			if err1 != nil {
				return false, fmt.Errorf("unwindTxLookup, snappy err: %w", err1)
			}
		}
		body := new(types.Body)
		if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
			return false, fmt.Errorf("unwindTxLookup, rlp decode err: %w", err)
		}
		for _, tx := range body.Transactions {
			txsToRemove = append(txsToRemove, tx.Hash().Bytes())
		}

		return true, nil
	}); err != nil {
		return err
	}
	// TODO: Do it in a batcn and update the progress
	for _, v := range txsToRemove {
		if err := db.Delete(dbutils.TxLookupPrefix, v); err != nil {
			return err
		}
	}
	if err := u.Done(db); err != nil {
		return fmt.Errorf("unwind TxLookup: %w", err)
	}
	return nil
}

func calculateTxLookupChunks(startBlock, endBlock uint64, numOfChunks int) [][]byte {
	if endBlock < startBlock+1000000 || numOfChunks < 2 {
		return nil
	}

	chunkSize := (endBlock - startBlock) / uint64(numOfChunks)
	var chunks = make([][]byte, numOfChunks-1)
	for i := uint64(1); i < uint64(numOfChunks); i++ {
		chunks[i-1] = dbutils.BlockBodyKey(startBlock + i * chunkSize, endHash)
	}
	return chunks
}
