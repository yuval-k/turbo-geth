package downloader

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/binary"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core/rawdb"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"io"
	"math/big"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"syscall"
	"testing"
	"time"
)

func TestGenerateTxLookup(t *testing.T) {

}

func GenerateTxLookups(chaindata string) {
	startTime := time.Now()
	db, err := ethdb.NewBoltDatabase(chaindata)
	check(err)
	//nolint: errcheck
	db.DeleteBucket(dbutils.TxLookupPrefix)
	log.Info("Open databased and deleted tx lookup bucket", "duration", time.Since(startTime))
	startTime = time.Now()
	sigs := make(chan os.Signal, 1)
	interruptCh := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		interruptCh <- true
	}()
	var blockNum uint64 = 1
	var finished bool
	for !finished {
		blockNum, finished = generateLoop(db, blockNum, interruptCh)
	}
	log.Info("All done", "duration", time.Since(startTime))
}

func generateLoop(db ethdb.Database, startBlock uint64, interruptCh chan bool) (uint64, bool) {
	startTime := time.Now()
	var lookups []uint64
	var entry [8]byte
	var blockNum = startBlock
	var interrupt bool
	var finished = true
	for !interrupt {
		blockHash := rawdb.ReadCanonicalHash(db, blockNum)
		body := rawdb.ReadBody(db, blockHash, blockNum)
		if body == nil {
			break
		}
		for txIndex, tx := range body.Transactions {
			copy(entry[:2], tx.Hash().Bytes()[:2])
			binary.BigEndian.PutUint32(entry[2:6], uint32(blockNum))
			binary.BigEndian.PutUint16(entry[6:8], uint16(txIndex))
			lookups = append(lookups, binary.BigEndian.Uint64(entry[:]))
		}
		blockNum++
		if blockNum%100000 == 0 {
			log.Info("Processed", "blocks", blockNum, "tx count", len(lookups))
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.Info("Memory", "alloc", int(m.Alloc/1024), "sys", int(m.Sys/1024), "numGC", int(m.NumGC))
		}
		// Check for interrupts
		select {
		case interrupt = <-interruptCh:
			log.Info("interrupted, please wait for cleanup...")
		default:
		}
		if len(lookups) >= 100000000 {
			log.Info("Reached specified number of transactions")
			finished = false
			break
		}
	}
	log.Info("Processed", "blocks", blockNum, "tx count", len(lookups))
	log.Info("Filling up lookup array done", "duration", time.Since(startTime))
	startTime = time.Now()
	sort.Slice(lookups, func(i, j int) bool {
		return lookups[i] < lookups[j]
	})
	log.Info("Sorting lookup array done", "duration", time.Since(startTime))
	if len(lookups) == 0 {
		return blockNum, true
	}
	startTime = time.Now()
	var rangeStartIdx int
	var range2Bytes uint64
	for i, lookup := range lookups {
		// Find the range where lookups entries share the same first two bytes
		if i == 0 {
			rangeStartIdx = 0
			range2Bytes = lookup & 0xffff000000000000
			continue
		}
		twoBytes := lookup & 0xffff000000000000
		if range2Bytes != twoBytes {
			// Range finished
			fillSortRange(db, lookups, entry[:], rangeStartIdx, i)
			rangeStartIdx = i
			range2Bytes = twoBytes
		}
		if i%1000000 == 0 {
			log.Info("Processed", "transactions", i)
		}
	}
	fillSortRange(db, lookups, entry[:], rangeStartIdx, len(lookups))
	log.Info("Second roung of sorting done", "duration", time.Since(startTime))
	startTime = time.Now()
	batch := db.NewBatch()
	var n big.Int
	for i, lookup := range lookups {
		binary.BigEndian.PutUint64(entry[:], lookup)
		blockNumber := uint64(binary.BigEndian.Uint32(entry[2:6]))
		txIndex := int(binary.BigEndian.Uint16(entry[6:8]))
		blockHash := rawdb.ReadCanonicalHash(db, blockNumber)
		body := rawdb.ReadBody(db, blockHash, blockNumber)
		tx := body.Transactions[txIndex]
		n.SetInt64(int64(blockNumber))
		err := batch.Put(dbutils.TxLookupPrefix, tx.Hash().Bytes(), common.CopyBytes(n.Bytes()))
		check(err)
		if i != 0 && i%1000000 == 0 {
			_, err = batch.Commit()
			check(err)
			log.Info("Commited", "transactions", i)
		}
	}
	_, err := batch.Commit()
	check(err)
	log.Info("Commited", "transactions", len(lookups))
	log.Info("Tx committing done", "duration", time.Since(startTime))
	return blockNum, finished
}

func generateLoop1(db ethdb.Database, startBlock uint64, interruptCh chan bool) (uint64, bool) {
	f, _ := os.OpenFile(".lookups.tmp",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	entry := make([]byte, 36)
	var blockNum = startBlock
	var interrupt bool
	var finished = true
	for !interrupt {
		blockHash := rawdb.ReadCanonicalHash(db, blockNum)
		body := rawdb.ReadBody(db, blockHash, blockNum)
		if body == nil {
			break
		}
		for _, tx := range body.Transactions {
			copy(entry[:32], tx.Hash().Bytes())
			binary.BigEndian.PutUint32(entry[32:], uint32(blockNum))
			_, _ = f.Write(append(entry, '\n'))
		}
		blockNum++
		// Check for interrupts
		select {
		case interrupt = <-interruptCh:
			log.Info("interrupted, please wait for cleanup...")
		default:
		}
	}
	return blockNum, finished
}

func fillSortRange(db rawdb.DatabaseReader, lookups []uint64, entry []byte, start, end int) {
	for j := start; j < end; j++ {
		binary.BigEndian.PutUint64(entry[:], lookups[j])
		blockNum := uint64(binary.BigEndian.Uint32(entry[2:6]))
		txIndex := int(binary.BigEndian.Uint16(entry[6:8]))
		blockHash := rawdb.ReadCanonicalHash(db, blockNum)
		body := rawdb.ReadBody(db, blockHash, blockNum)
		tx := body.Transactions[txIndex]
		copy(entry[:2], tx.Hash().Bytes()[2:4])
		lookups[j] = binary.BigEndian.Uint64(entry[:])
	}
	sort.Slice(lookups[start:end], func(i, j int) bool {
		return lookups[i] < lookups[j]
	})
}

func GenerateTxLookups1(chaindata string, block int) {
	startTime := time.Now()
	db, err := ethdb.NewBoltDatabase(chaindata)
	check(err)
	//nolint: errcheck
	db.DeleteBucket(dbutils.TxLookupPrefix)
	log.Info("Open databased and deleted tx lookup bucket", "duration", time.Since(startTime))
	startTime = time.Now()
	sigs := make(chan os.Signal, 1)
	interruptCh := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		interruptCh <- true
	}()
	var blockNum uint64 = 1
	var interrupt bool
	var txcount int
	var n big.Int
	batch := db.NewBatch()
	for !interrupt {
		blockHash := rawdb.ReadCanonicalHash(db, blockNum)
		body := rawdb.ReadBody(db, blockHash, blockNum)
		if body == nil {
			break
		}
		for _, tx := range body.Transactions {
			txcount++
			n.SetInt64(int64(blockNum))
			err = batch.Put(dbutils.TxLookupPrefix, tx.Hash().Bytes(), common.CopyBytes(n.Bytes()))
			check(err)
			if txcount%100000 == 0 {
				_, err = batch.Commit()
				check(err)
			}
			if txcount%1000000 == 0 {
				log.Info("Commited", "transactions", txcount)
			}
		}
		blockNum++
		if blockNum%100000 == 0 {
			log.Info("Processed", "blocks", blockNum)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.Info("Memory", "alloc", int(m.Alloc/1024), "sys", int(m.Sys/1024), "numGC", int(m.NumGC))
		}
		// Check for interrupts
		select {
		case interrupt = <-interruptCh:
			log.Info("interrupted, please wait for cleanup...")
		default:
		}
		if block != 1 && int(blockNum) > block {
			log.Info("Reached specified block count")
			break
		}
	}
	_, err = batch.Commit()
	check(err)
	log.Info("Commited", "transactions", txcount)
	log.Info("Processed", "blocks", blockNum)
	log.Info("Tx committing done", "duration", time.Since(startTime))
}

func GenerateTxLookups2(chaindata string) {
	startTime := time.Now()
	db, err := ethdb.NewBoltDatabase(chaindata)
	check(err)
	//nolint: errcheck
	db.DeleteBucket(dbutils.TxLookupPrefix)
	log.Info("Open databased and deleted tx lookup bucket", "duration", time.Since(startTime))
	startTime = time.Now()
	sigs := make(chan os.Signal, 1)
	interruptCh := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		interruptCh <- true
	}()
	var blockNum uint64 = 1
	generateTxLookups2(db, blockNum, interruptCh)
	log.Info("All done", "duration", time.Since(startTime))
}

type LookupFile struct {
	reader io.Reader
	file   *os.File
	buffer []byte
	pos    uint64
}

type Entries []byte

type HeapElem struct {
	val   []byte
	index int
}

type Heap []HeapElem

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	return bytes.Compare(h[i].val, h[j].val) < 0
}
func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(HeapElem))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (a Entries) Len() int {
	return len(a) / 35
}

func (a Entries) Less(i, j int) bool {
	return bytes.Compare(a[35*i:35*i+35], a[35*j:35*j+35]) < 0
}

func (a Entries) Swap(i, j int) {
	tmp := common.CopyBytes(a[35*i : 35*i+35])
	copy(a[35*i:35*i+35], a[35*j:35*j+35])
	copy(a[35*j:35*j+35], tmp)
}

func insertInFileForLookups2(file *os.File, entries Entries, it uint64) {
	sorted := entries[:35*it]
	sort.Sort(sorted)
	_, err := file.Write(sorted)
	check(err)
	log.Info("File Insertion Occured")
}

func generateTxLookups2(db *ethdb.BoltDatabase, startBlock uint64, interruptCh chan bool) {
	var bufferLen int = 143360 // 35 * 4096
	var count uint64 = 5000000
	var entries Entries = make([]byte, count*35)
	bn := make([]byte, 3)
	var lookups []LookupFile
	var iterations uint64
	var blockNum = startBlock
	var interrupt bool
	filename := fmt.Sprintf(".lookups_%d.tmp", len(lookups))
	fileTmp, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	for !interrupt {
		blockHash := rawdb.ReadCanonicalHash(db, blockNum)
		body := rawdb.ReadBody(db, blockHash, blockNum)

		if body == nil {
			log.Info("Now Inserting to file")
			insertInFileForLookups2(fileTmp, entries, iterations)
			lookups = append(lookups, LookupFile{nil, nil, nil, 35})
			iterations = 0
			entries = nil
			fileTmp.Close()
			break
		}

		select {
		case interrupt = <-interruptCh:
			log.Info("interrupted, please wait for cleanup...")
		default:
		}
		bn[0] = byte(blockNum >> 16)
		bn[1] = byte(blockNum >> 8)
		bn[2] = byte(blockNum)

		for _, tx := range body.Transactions {
			copy(entries[35*iterations:], tx.Hash().Bytes())
			copy(entries[35*iterations+32:], bn)
			iterations++
			if iterations == count {
				log.Info("Now Inserting to file")
				insertInFileForLookups2(fileTmp, entries, iterations)
				lookups = append(lookups, LookupFile{nil, nil, nil, 35})
				iterations = 0
				fileTmp.Close()
				filename = fmt.Sprintf(".lookups_%d.tmp", len(lookups))
				fileTmp, _ = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
			}
		}
		blockNum++
		if blockNum%100000 == 0 {
			log.Info("Processed", "blocks", blockNum, "iterations", iterations)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.Info("Memory", "alloc", int(m.Alloc/1024), "sys", int(m.Sys/1024), "numGC", int(m.NumGC))
		}
	}
	batch := db.NewBatch()
	h := &Heap{}
	heap.Init(h)
	for i := range lookups {
		file, err := os.Open(fmt.Sprintf(".lookups_%d.tmp", i))
		check(err)
		lookups[i].file = file
		lookups[i].reader = bufio.NewReader(file)
		lookups[i].buffer = make([]byte, bufferLen)
		check(err)
		n, err := lookups[i].file.Read(lookups[i].buffer)
		if n != bufferLen {
			lookups[i].buffer = lookups[i].buffer[:n]
		}
		heap.Push(h, HeapElem{lookups[i].buffer[:35], i})
	}

	for !interrupt && len(*h) != 0 {
		val := (heap.Pop(h)).(HeapElem)
		if lookups[val.index].pos == uint64(bufferLen) {
			if val.val[32] != 0 {
				err := batch.Put(dbutils.TxLookupPrefix, val.val[:32], common.CopyBytes(val.val[32:]))
				check(err)
			} else {
				err := batch.Put(dbutils.TxLookupPrefix, val.val[:32], common.CopyBytes(val.val[33:]))
				check(err)
			}
			n, _ := lookups[val.index].reader.Read(lookups[val.index].buffer)
			iterations++
			if n == 0 {
				err := lookups[val.index].file.Close()
				check(err)
				os.Remove(fmt.Sprintf(".lookups_%d.tmp", val.index))
			} else {
				if n != bufferLen {
					lookups[val.index].buffer = lookups[val.index].buffer[:n]
				}
				lookups[val.index].pos = 35
				heap.Push(h, HeapElem{lookups[val.index].buffer[:35], val.index})
			}
			continue
		}

		heap.Push(h, HeapElem{lookups[val.index].buffer[lookups[val.index].pos : lookups[val.index].pos+35], val.index})
		lookups[val.index].pos += 35
		iterations++
		if val.val[32] != 0 {
			err := batch.Put(dbutils.TxLookupPrefix, val.val[:32], common.CopyBytes(val.val[32:]))
			check(err)
		} else {
			err := batch.Put(dbutils.TxLookupPrefix, val.val[:32], common.CopyBytes(val.val[33:]))
			check(err)
		}

		if iterations%1000000 == 0 {
			batch.Commit()
			log.Info("Commit Occured", "progress", iterations)
		}
		select {
		case interrupt = <-interruptCh:
			log.Info("interrupted, please wait for cleanup...")
		default:
		}
	}
	batch.Commit()
	batch.Close()
}

func ValidateTxLookups2(chaindata string) {
	startTime := time.Now()
	db, err := ethdb.NewBoltDatabase(chaindata)
	check(err)
	//nolint: errcheck
	startTime = time.Now()
	sigs := make(chan os.Signal, 1)
	interruptCh := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		interruptCh <- true
	}()
	var blockNum uint64 = 1
	validateTxLookups2(db, blockNum, interruptCh)
	log.Info("All done", "duration", time.Since(startTime))
}

func validateTxLookups2(db *ethdb.BoltDatabase, startBlock uint64, interruptCh chan bool) {
	blockNum := startBlock
	iterations := 0
	var interrupt bool
	// Validation Process
	blockBytes := big.NewInt(0)
	for !interrupt {
		blockHash := rawdb.ReadCanonicalHash(db, blockNum)
		body := rawdb.ReadBody(db, blockHash, blockNum)

		if body == nil {
			break
		}

		select {
		case interrupt = <-interruptCh:
			log.Info("interrupted, please wait for cleanup...")
		default:
		}
		blockBytes.SetUint64(blockNum)
		bn := blockBytes.Bytes()

		for _, tx := range body.Transactions {
			val, err := db.Get(dbutils.TxLookupPrefix, tx.Hash().Bytes())
			iterations++
			if iterations%100000 == 0 {
				log.Info("Validated", "entries", iterations, "number", blockNum)
			}
			if bytes.Compare(val, bn) != 0 {
				check(err)
				panic(fmt.Sprintf("Validation process failed(%d). Expected %b, got %b", iterations, bn, val))
			}
		}
		blockNum++
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}