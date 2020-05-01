package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davecgh/go-spew/spew"

	static "github.com/ledgerwatch/turbo-geth/cmd/jumps"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core/asm"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

const dbPath = "/mnt/sdb/contract_codes_cleanupBeggining_trimMetadata_cleanupErrors" // _goodJumps

func main() {
	/*
	// uncomment to prepare the contracts database
		p := &processor{from: "/mnt/sdb/contract_codes"}
		p.cleanupBeggining().
			trimMetadata().
			cleanupErrors().
			jumpsPaths()
	*/
	p := &processor{from: dbPath}
	p.jumpsPaths()
}

func (p *processor) goodJumps() {
	jumps := 0
	i := 0
	unknownJumpDest := 0
	nilJumpDest := 0
	nilJumpDestOPs := make(map[vm.OpCode]int)

	errJumpDest := 0
	okJumpDest := 0
	nilPrev := 0

	db, to, dbCleaned, toErr, dbErr := openDBs(p.from)
	p.to = to
	from := p.from

	fmt.Printf("goodJumps started. From %q to %q\n", p.from, p.to)
	defer func() {
		db.Close()
		dbCleaned.Close()
		dbErr.Close()

		fmt.Printf("total contracts %q: %d\n", from, i)
		fmt.Printf("correct contracts %q: %d\n", to, okJumpDest)
		fmt.Printf("incorrect contracts %q: %d\n", toErr, errJumpDest+unknownJumpDest+nilJumpDest)
		fmt.Printf("goodJumps finished\n\n\n")

		p.from = p.to
		p.to = ""
	}()

	var it *asm.InstructionIterator
	storageCleaned := NewToPut(dbutils.CodeBucket, BunchSize, dbCleaned)
	storageErr := NewToPut(dbutils.CodeBucket, BunchSize, dbErr)
	err := db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		i++
		if i%10000 == 0 {
			fmt.Println("done", i)
		}

		jumpDests, latest, _, _ := p.getJumpDests(v)
		it = asm.NewInstructionIterator(v)

		okJumps := true
	codeLoop:
		for it.Next() {
			op := it.Op()

			if op == vm.JUMP || op == vm.JUMPI {
				jumps++

				n := 1
				var prev *asm.Command
				for {
					prev = it.Prev(n)
					if prev == nil {
						nilPrev++
						okJumps = false
						continue codeLoop
					}
					if len(prev.Arg) != 0 {
						break
					}
					n++
				}

				argStr := fmt.Sprintf("%x", prev.Arg)
				arg, err := strconv.ParseUint(argStr, 16, 64)
				if err != nil {
					errJumpDest++
					okJumps = false
					//fmt.Println("error", spew.Sdump(prev.Arg), err)
					continue
				}

				if _, ok := jumpDests[arg]; !ok && arg > 2 && arg < latest {
					unknownJumpDest++
					//fmt.Println("dest not found")
					okJumps = false
					continue
				}

				okJumpDest++
			}
		}

		if okJumps {
			storageCleaned.Add(k, v)
		} else {
			storageErr.Add(k, v)
		}
		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}
	storageCleaned.Flush()
	storageErr.Flush()
	fmt.Println("done", i)

	fmt.Println("different contracts", i)
	fmt.Println("different jumps", jumps)
	fmt.Println("errJumpDest", errJumpDest)
	fmt.Println("nilJumpDest", nilJumpDest)
	fmt.Println("unknownJumpDest", unknownJumpDest)
	fmt.Println("nilPrev", nilPrev)
	fmt.Println("okJumpDest", okJumpDest)

	spew.Dump(nilJumpDestOPs)
}

/*
func (p *processor) jumpsPathsBackwards() {
	fmt.Printf("jumpsPaths started. From %q to %q\n", p.from, p.to)
	defer fmt.Println("jumpsPaths finished")

	db, err := ethdb.NewBoltDatabase(p.from)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jumps := 0
	i := 0
	unknownJumpDest := 0
	nilJumpDest := 0

	errJumpDest := 0
	okJumpDest := 0
	nilPrev := 0

	executionPaths := map[common.Hash][][]asm.Command{}
	_ = executionPaths

	destsPerContract := make(map[int]int, 500_000)

	var codeAddress common.Hash
	res := struct{
		contracts map[bool]int
		ifsTotal map[bool]int
	} {
		contracts: make(map[bool]int),
		ifsTotal: make(map[bool]int),
	}

	isStatic := true
	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		i++
		if i%10000 == 0 {
			fmt.Println("done", i)
		}

		codeAddress = common.BytesToHash(k)
		var executionPath []asm.Command

		jumpDests, latest, funcsJumps, err := p.getJumpDests(v)
		if err != nil {
			return true, nil
		}
		_ = funcsJumps

		destsPerContract[len(jumpDests)]++

		it := asm.NewInstructionIterator(v)
		_ = executionPath
		_ = codeAddress
		_ = latest

		for it.Next() {
			if it.Op() == vm.JUMP || it.Op() == vm.JUMPI {
				history, found := it.PreviousUntilStackValue(1)
				_ = history
				res.ifsTotal[found]++

				if !found {
					isStatic = false

					spew.Dump(history)
					spew.Dump(asm.PrintDisassembledBytes(v))
					return false, nil
				}
			}
		}

		res.contracts[isStatic]++

		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}
	fmt.Println("done", i)

	fmt.Println("different contracts", i)
	fmt.Println("different jumps", jumps)
	fmt.Println("errJumpDest", errJumpDest)
	fmt.Println("nilJumpDest", nilJumpDest)
	fmt.Println("unknownJumpDest", unknownJumpDest)
	fmt.Println("nilPrev", nilPrev)
	fmt.Println("okJumpDest", okJumpDest)

	keys := make([]int, len(destsPerContract))
	j := 0
	for k := range destsPerContract {
		keys[j] = k
		j++
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("Contracts per contract %d: total contracts %d\n", k, destsPerContract[k])
	}

	spew.Dump(res)
}
*/

func (p *processor) jumpsPaths() {
	fmt.Printf("jumpsPaths started. From %q to %q\n", p.from, p.to)
	defer fmt.Println("jumpsPaths finished")

	db, err := ethdb.NewBoltDatabase(p.from)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	i := new(uint64)

	res := struct {
		*sync.RWMutex

		static    map[common.Hash]struct{}
		notstatic map[common.Hash]struct{}

		errJump         map[common.Hash]struct{}
		novalueStatic   map[common.Hash]struct{}
		errTooManyJumps map[common.Hash]struct{}
		timeout         map[common.Hash]struct{}
	}{
		RWMutex:   new(sync.RWMutex),
		static:    make(map[common.Hash]struct{}, 200_000),
		notstatic: make(map[common.Hash]struct{}, 2_000),

		errJump:       make(map[common.Hash]struct{}, 1_000),
		novalueStatic: make(map[common.Hash]struct{}, 1_000),
		timeout:       make(map[common.Hash]struct{}, 1_000),
	}

	type job struct {
		fn   func(ctx context.Context, k, v []byte) error
		k, v []byte
	}
	ch := make(chan job, 10000)

	var numWorkers = 2*runtime.NumCPU()
	wg := sync.WaitGroup{}

	for n := 0; n < numWorkers; n++ {
		wg.Add(1)

		go func() {
			for jb := range ch {

				done := atomic.AddUint64(i, 1)
				if done%100 == 0 {
					fmt.Println("done", done)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				jb.fn(ctx, jb.k, jb.v)
				cancel()
			}
			defer wg.Done()
		}()
	}

	testContracts := [][]byte{
		{0, 0, 52, 48, 27, 92, 143, 64, 73, 215, 89, 9, 44, 209, 25, 4, 245, 137, 215, 62, 98, 131, 188, 46, 161, 226, 125, 62, 173, 184, 74, 132},        // 5
		{0, 0, 68, 62, 128, 192, 9, 218, 248, 40, 202, 58, 152, 195, 40, 42, 3, 188, 148, 228, 49, 3, 249, 220, 123, 120, 55, 216, 123, 62, 249, 32},      // 201
		{15, 83, 113, 94, 36, 52, 138, 104, 41, 166, 140, 58, 101, 223, 8, 20, 71, 149, 163, 216, 193, 228, 127, 125, 15, 211, 2, 134, 34, 228, 141, 182}, // 6

		{0, 1, 140, 92, 146, 30, 249, 244, 87, 138, 251, 103, 93, 32, 165, 225, 9, 180, 66, 140, 16, 23, 95, 162, 74, 60, 191, 202, 149, 100, 36, 208},      // 90
		{0, 0, 165, 215, 192, 253, 46, 213, 109, 128, 46, 84, 136, 231, 159, 200, 139, 170, 38, 65, 251, 138, 79, 165, 16, 132, 186, 47, 42, 106, 245, 124}, // 322
		{0, 50, 69, 225, 158, 184, 55, 168, 227, 196, 43, 30, 246, 157, 69, 46, 88, 234, 22, 162, 72, 82, 76, 245, 152, 0, 231, 73, 43, 148, 207, 240},      //33
	}
	_ = testContracts

	const debug = false

	//err = db.Walk(dbutils.CodeBucket, testContracts[5], 32, func(key, value []byte) (bool, error) {
	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(key, value []byte) (bool, error) {
		ch <- job{
			fn: func(ctx context.Context, k, v []byte) error {
				codeAddress := common.BytesToHash(k)
				runner := static.NewContractRunner(v, debug)

				pathsDone, err := runner.Run(ctx)
				//fmt.Println(codeAddress.String(), "end")

				res.Lock()
				defer res.Unlock()

				if errors.Is(err, static.ErrTimeout) {
					res.timeout[codeAddress] = struct{}{}

					if debug {
						fmt.Printf("timeout error. code address %s. jumpi %d. paths done %d. key %v\ncode:\n%v\n\n",
							codeAddress.String(), len(runner.Jumpi), pathsDone, key, value)
					}
					return nil
				}

				if errors.Is(err, static.ErrNonStatic) {
					res.notstatic[codeAddress] = struct{}{}
					return nil
				}

				if errors.Is(err, static.ErrNoValueStatic) {
					res.novalueStatic[codeAddress] = struct{}{}
					return nil
				}

				// good cases
				if errors.Is(err, static.ErrInvalidJump) {
					res.errJump[codeAddress] = struct{}{}
					res.static[codeAddress] = struct{}{}
					return nil
				}

				if errors.Is(err, static.ErrUnknownOpcode) {
					// nothing to do
				}

				if errors.Is(err, static.ErrMaxJumps) {
					res.errTooManyJumps[codeAddress] = struct{}{}
				}

				res.static[codeAddress] = struct{}{}
				return nil
			},
			k: key,
			v: value,
		}

		return true, nil
	})
	close(ch)

	if err != nil {
		log.Println("err", err)
	}

	wg.Wait()

	fmt.Println("done contracts", i)

	fmt.Println("static contracts", len(res.static))
	fmt.Println("notstatic contracts", len(res.notstatic))
	fmt.Println("errJump contracts", len(res.errJump))
	fmt.Println("novalueStatic contracts", len(res.novalueStatic))
	fmt.Println("errTooManyJumps contracts", len(res.errTooManyJumps))
	fmt.Println("timeout contracts", len(res.timeout))

	//spew.Dump(res)
}

func (p *processor) getJumpDests(code []byte) (map[uint64]struct{}, uint64, map[uint64]struct{}, error) {
	jumpDest := make(map[uint64]struct{}, 10)
	funcJumpDest := make(map[uint64]struct{}, 10)
	it := asm.NewInstructionIterator(code)

	funcsEnded := false
	latest := uint64(0)

	var prevOp vm.OpCode
	var prevArgs []byte
	for it.Next() {
		op := it.Op()
		pc := it.PC()

		if op == vm.JUMPDEST {
			pos := pc
			// if pos > 2 {
			jumpDest[pos] = struct{}{}
			// }
		}

		if op == vm.REVERT || op == vm.JUMPDEST {
			funcsEnded = true
		}
		if !funcsEnded && op == vm.JUMPI {
			if prevOp.IsPush() {
				argStr := fmt.Sprintf("%x", prevArgs)
				arg, err := strconv.ParseUint(argStr, 16, 64)
				if err != nil {
					return nil, 0, nil, err
				}
				funcJumpDest[arg] = struct{}{}
			}
		}

		latest = pc

		prevOp = it.Op()
		prevArgs = it.Arg()
	}

	// default jumps to cancel jump
	jumpDest[0] = struct{}{}
	jumpDest[1] = struct{}{}
	jumpDest[2] = struct{}{}
	jumpDest[3] = struct{}{}
	jumpDest[4] = struct{}{}
	jumpDest[5] = struct{}{}
	jumpDest[6] = struct{}{}
	jumpDest[7] = struct{}{}
	/*
		00000: PUSH1 0x80
		00002: PUSH1 0x40
		00004: MSTORE
		00005: CALLVALUE
		00006: DUP1
		00007: ISZERO
	*/

	return jumpDest, latest, funcJumpDest, nil
}

func jumps() {
	db, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes_cleaned_begginings")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jumpPatterns := make(map[int]map[string]int)

	jumps := 0
	i := 0
	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		defer func() { i++ }()

		it := asm.NewInstructionIterator(v)

		for it.Next() {
			op := it.Op()

			if op == vm.JUMP || op == vm.JUMPI {
				//jumpsStats[asm.Commands(prev).String()]++
			}
		}

		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}

	//spew.Dump(jumpPatterns, len(jumpPatterns), i)

	fmt.Println("different contracts", i)
	fmt.Println("different jumps", jumps)
	fmt.Println()

	type pair struct {
		k string
		v int
	}

	for n, jumpStats := range jumpPatterns {
		f, err := os.OpenFile("./patterns_"+strconv.Itoa(n), os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		var jumpPatternsArr []pair
		for k, v := range jumpStats {
			jumpPatternsArr = append(jumpPatternsArr, pair{k, v})
		}

		sort.Slice(jumpPatternsArr, func(i, j int) bool {
			return jumpPatternsArr[i].v > jumpPatternsArr[j].v
		})

		for _, kv := range jumpPatternsArr {
			fmt.Fprintf(f, "%s,%d,%.4f\n", kv.k, kv.v, float64(kv.v)/float64(jumps)*100)
		}

		f.Close()
	}
}

type processor struct {
	from string
	to   string
}

func (p *processor) cleanupBeggining() *processor {
	i := 0
	correct := 0
	errCount := 0

	db, to, dbCleaned, toErr, dbErr := openDBs(p.from)
	p.to = to

	fmt.Printf("cleanupBeggining started. From %q to %q\n", p.from, p.to)

	defer func() {
		db.Close()
		dbCleaned.Close()
		dbErr.Close()

		fmt.Printf("total contracts %q: %d\n", p.from, i)
		fmt.Printf("correct contracts %q: %d\n", p.to, correct)
		fmt.Printf("incorrect contracts %q: %d\n", toErr, errCount)
		fmt.Printf("cleanupBeggining finished\n\n\n")

		p.from = p.to
		p.to = ""
	}()

	var (
		code   []byte
		it     *asm.InstructionIterator
		prefix string
		ok     bool
	)
	badPrefixes := make(map[string]int)
	storageCleaned := NewToPut(dbutils.CodeBucket, BunchSize, dbCleaned)
	storageErr := NewToPut(dbutils.CodeBucket, BunchSize, dbErr)
	err := db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		i++
		if i%10000 == 0 {
			fmt.Printf("done %d, fail %d, succ %d\n", i, errCount, correct)
		}

		code = v[:trimLen(v, 5)]
		it = asm.NewInstructionIterator(code)

		if prefix, ok = hasSolidityPrefix(it); !ok {
			errCount++
			badPrefixes[prefix]++
			storageErr.Add(k, v)
			return true, nil
		}

		correct++
		storageCleaned.Add(k, v)
		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}
	storageCleaned.Flush()
	storageErr.Flush()
	fmt.Printf("done %d, fail %d, succ %d\n", i, errCount, correct)

	keys := make([]string, 0, len(badPrefixes))
	for pref := range badPrefixes {
		keys = append(keys, pref)
	}

	sort.Slice(keys, func(i, j int) bool {
		return badPrefixes[keys[i]] > badPrefixes[keys[j]]
	})

	fmt.Println("\nBad prefixes:")
	for _, pref := range keys {
		fmt.Printf("%s: %d\n", pref, badPrefixes[pref])
	}
	fmt.Printf("\n\n\n")

	return p
}

func (p *processor) cleanupErrors() *processor {
	i := 0
	correct := 0
	errCount := 0

	db, to, dbCleaned, toErr, dbErr := openDBs(p.from)
	p.to = to
	from := p.from

	fmt.Printf("cleanupErrors started. From %q to %q\n", p.from, p.to)

	defer func() {
		db.Close()
		dbCleaned.Close()
		dbErr.Close()

		fmt.Printf("total contracts %q: %d\n", from, i)
		fmt.Printf("correct contracts %q: %d\n", to, correct)
		fmt.Printf("incorrect contracts %q: %d\n", toErr, errCount)
		fmt.Printf("cleanupErrors finished\n\n\n")

		p.from = p.to
		p.to = ""
	}()

	var it *asm.InstructionIterator
	storageCleaned := NewToPut(dbutils.CodeBucket, BunchSize, dbCleaned)
	storageErr := NewToPut(dbutils.CodeBucket, BunchSize, dbErr)
	err := db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		i++
		if i%10000 == 0 {
			fmt.Printf("done %d, fail %d, succ %d\n", i, errCount, correct)
		}

		it = asm.NewInstructionIterator(v)

		for it.Next() {
			if it.Error() != nil {
				errCount++
				storageErr.Add(k, v)
				return true, nil
			}
		}

		correct++
		storageCleaned.Add(k, v)
		return true, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	storageCleaned.Flush()
	storageErr.Flush()
	fmt.Printf("done %d, fail %d, succ %d\n", i, errCount, correct)

	return p
}

func (p *processor) trimMetadata() *processor {
	i := 0
	correct := 0
	errCount := 0

	db, to, dbCleaned, toErr, dbErr := openDBs(p.from)
	p.to = to
	from := p.from

	fmt.Printf("trimMetadata started. From %q to %q\n", p.from, p.to)

	defer func() {
		db.Close()
		dbCleaned.Close()
		dbErr.Close()

		fmt.Printf("total contracts %q: %d\n", from, i)
		fmt.Printf("correct contracts %q: %d\n", to, correct)
		fmt.Printf("incorrect contracts %q: %d\n", toErr, errCount)
		fmt.Printf("trimMetadata finished\n\n\n")

		p.from = p.to
		p.to = ""
	}()

	incorrectCodes := make(map[vm.OpCode]int)

	patterns := [][]byte{
		{0xa1, 0x65, 0x62, 0x7a, 0x7a},
		{0xfe, 0xa2, 0x65, 0x62, 0x7a, 0x7a},
	}

	stats := map[string]uint64{
		"trimmedOnce": 0,
		"all":         0,
	}

	var it *asm.InstructionIterator
	storageCleaned := NewToPut(dbutils.CodeBucket, BunchSize, dbCleaned)
	storageErr := NewToPut(dbutils.CodeBucket, BunchSize, dbErr)
	err := db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		i++
		if i%10000 == 0 {
			fmt.Printf("done %d, fail %d, succ %d\n", i, errCount, correct)
		}

		stats["all"]++
		trimmedCode := v
		for _, p := range patterns {
			if idx := bytes.LastIndex(v, p); idx != -1 {
				trimmedCode = v[:idx]
				stats["trimmedOnce"]++
				break
			}
		}

		it = asm.NewInstructionIterator(trimmedCode)
		var wasUnknown bool
		for it.Next() {
			op := it.Op()

			if op.IsUnknown() {
				wasUnknown = true
				incorrectCodes[op]++
			}
		}

		if wasUnknown {
			errCount++
			storageErr.Add(k, trimmedCode)

			/*
				res, err := asm.DisassembledBytes(trimmedCode)
				if err != nil {
					res, innerErr := asm.DisassembledBytes(v)
					if innerErr == nil {
						toFile := []byte(res)
						ioutil.WriteFile("/mnt/sdb/go/src/github.com/ledgerwatch/turbo-geth/cmd/jumps/prints/errors/err_TRIMMED_contract_"+strconv.Itoa(i)+".print",
							[]byte(fmt.Sprintf("Trimmed:\n%v\n\n\nOriginal:\n%v", trimmedCode, toFile)), 0666)
					} else {
						ioutil.WriteFile("/mnt/sdb/go/src/github.com/ledgerwatch/turbo-geth/cmd/jumps/prints/errors/err_DISASM_contract_"+strconv.Itoa(i)+".print",
							[]byte(fmt.Sprintf("%s\n\n\n%v", err.Error(), res)), 0666)
					}
					return true, nil
				}

				ioutil.WriteFile("/mnt/sdb/go/src/github.com/ledgerwatch/turbo-geth/cmd/jumps/prints/errors/err_contract_"+strconv.Itoa(i)+".print", []byte(res), 0666)
			*/
			return true, nil
		}

		correct++
		storageCleaned.Add(k, trimmedCode)
		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}
	storageCleaned.Flush()
	storageErr.Flush()
	fmt.Printf("done %d, fail %d, succ %d\n", i, errCount, correct)

	fmt.Println("trimmed metadata", spew.Sdump(stats))

	type pair struct {
		k vm.OpCode
		v int
	}

	var incorrectCodesArr []pair
	for k, v := range incorrectCodes {
		incorrectCodesArr = append(incorrectCodesArr, pair{k, v})
	}

	sort.Slice(incorrectCodesArr, func(i, j int) bool {
		return incorrectCodesArr[i].v > incorrectCodesArr[j].v
	})

	for _, kv := range incorrectCodesArr {
		fmt.Printf("%d: %d\n", kv.k, kv.v)
	}

	return p
}

func addFileSuffix(from string) string {
	suf := currentFunction()
	return addSuffix(from, suf)
}

func currentFunction() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	var name string
	f := strings.Split(frame.Function, ".")
	if len(f) > 0 {
		name = f[2]
	}
	return name
}

func addSuffix(file, suffix string) string {
	f := strings.Split(file, ".")
	if len(f) > 0 {
		file = f[0]
	}

	var ext string
	if len(f) > 1 {
		ext = f[1]
	}
	return fmt.Sprintf("%s_%s%s", file, suffix, ext)
}

func openDBs(from string) (*ethdb.BoltDatabase, string, *ethdb.BoltDatabase, string, *ethdb.BoltDatabase) {
	db, err := ethdb.NewBoltDatabase(from)
	if err != nil {
		log.Fatal(err)
	}

	to := addFileSuffix(from)
	dbCleaned, err := ethdb.NewBoltDatabase(to)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	toErr := addSuffix(to, "ERROR")
	dbErr, err := ethdb.NewBoltDatabase(toErr)
	if err != nil {
		db.Close()
		dbCleaned.Close()
		log.Fatal(err)
	}

	return db, to, dbCleaned, toErr, dbErr
}

func trimLen(v []byte, length int) int {
	if len(v) < length {
		return len(v)
	}
	return length
}

func begginings() {
	db, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	begginings := make(map[[3]string]int)
	beggining := [3]string{}

	i := 0
	var it *asm.InstructionIterator
	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		it = asm.NewInstructionIterator(v)

		n := 0
		for it.Next() {
			if n < 3 {
				beggining[n] = it.Op().String()
			}
			if n == 2 {
				begginings[beggining]++
			}
			n++
		}

		i++
		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}

	type pair struct {
		k [3]string
		v int
	}

	var begginingsArr []pair
	for k, v := range begginings {
		begginingsArr = append(begginingsArr, pair{k, v})
	}

	sort.Slice(begginingsArr, func(i, j int) bool {
		return begginingsArr[i].v > begginingsArr[j].v
	})

	for _, kv := range begginingsArr {
		fmt.Printf("%s, %d\n", kv.k, kv.v)
	}
}

func jumpsPrefixes() {
	db, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes_cleaned_begginings")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jumpPatterns := make(map[int]map[string]int)

	const maxContext = 20
	const beforeJump = false
	const beforeJumpDest = true

	jumps := 0
	i := 0
	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		/*
			if i <= 3 {
				return false, nil
			}
		*/
		defer func() { i++ }()

		it := asm.NewInstructionIterator(v)

		for it.Next() {
			op := it.Op()

			if op == vm.JUMP || op == vm.JUMPI {
				jumps++

			loopJump:
				for n := 1; n <= maxContext; n++ {
					jumpsStats, ok := jumpPatterns[n]
					if !ok {
						jumpsStats = make(map[string]int)
						jumpPatterns[n] = jumpsStats
					}

					var prev []asm.Command
					if beforeJump {
						prev = it.PreviousBeforeJump(n)
					} else if beforeJumpDest {
						prev = it.PreviousBefore(n, vm.JUMPDEST)
					} else {
						prev = it.Previous(n)
					}

					if len(prev) == 0 {
						jumpsStats[""]++
						continue
					}

					switch prev[len(prev)-1].Op {
					/*case vm.DUP, vm.DUP1, vm.DUP2, vm.DUP3, vm.DUP4, vm.DUP5, vm.DUP6, vm.DUP7, vm.DUP8, vm.DUP9, vm.DUP10, vm.DUP11, vm.DUP12, vm.DUP13, vm.DUP14, vm.DUP15, vm.DUP16:
					jumpsStats[vm.DUP.String()]++
					continue*/
					case vm.PUSH, vm.PUSH1, vm.PUSH2, vm.PUSH3, vm.PUSH4, vm.PUSH5, vm.PUSH6, vm.PUSH7, vm.PUSH8, vm.PUSH9, vm.PUSH10, vm.PUSH11, vm.PUSH12, vm.PUSH13, vm.PUSH14, vm.PUSH15, vm.PUSH16, vm.PUSH17, vm.PUSH18, vm.PUSH19, vm.PUSH20, vm.PUSH21, vm.PUSH22, vm.PUSH23, vm.PUSH24, vm.PUSH25, vm.PUSH26, vm.PUSH27, vm.PUSH28, vm.PUSH29, vm.PUSH30, vm.PUSH31, vm.PUSH32:
						jumpsStats[vm.PUSH.String()]++
						continue
					}

					// PUSH+SINGLE_OP
					if len(prev) >= 2 {
						/*
							POP 1 + PUSH 1
							ISZERO
							NOT
							BALANCE
						*/

						switch prev[len(prev)-1].Op {
						case vm.ISZERO, vm.NOT, vm.BALANCE:
							switch prev[len(prev)-2].Op {
							case vm.DUP, vm.DUP1, vm.DUP2, vm.DUP3, vm.DUP4, vm.DUP5, vm.DUP6, vm.DUP7, vm.DUP8, vm.DUP9, vm.DUP10, vm.DUP11, vm.DUP12, vm.DUP13, vm.DUP14, vm.DUP15, vm.DUP16:
								jumpsStats[vm.DUP.String()]++
								continue loopJump
							case vm.PUSH, vm.PUSH1, vm.PUSH2, vm.PUSH3, vm.PUSH4, vm.PUSH5, vm.PUSH6, vm.PUSH7, vm.PUSH8, vm.PUSH9, vm.PUSH10, vm.PUSH11, vm.PUSH12, vm.PUSH13, vm.PUSH14, vm.PUSH15, vm.PUSH16, vm.PUSH17, vm.PUSH18, vm.PUSH19, vm.PUSH20, vm.PUSH21, vm.PUSH22, vm.PUSH23, vm.PUSH24, vm.PUSH25, vm.PUSH26, vm.PUSH27, vm.PUSH28, vm.PUSH29, vm.PUSH30, vm.PUSH31, vm.PUSH32:
								jumpsStats[vm.PUSH.String()]++
								continue loopJump
							}
						}
					}

					/*
						POP 0 PUSH 1
						ADDRESS
						ORIGIN
						CALLER
						CALLVALUE
						CALLDATASIZE
						CODESIZE
						GASPRICE
						COINBASE
						TIMESTAMP
						NUMBER
						DIFFICULTY
						GASLIMIT
						PC
						GAS
					*/

					jumpsStats[asm.Commands(prev).String()]++
				}
			}
		}

		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}

	//spew.Dump(jumpPatterns, len(jumpPatterns), i)

	fmt.Println("different contracts", i)
	fmt.Println("different jumps", jumps)
	fmt.Println()

	type pair struct {
		k string
		v int
	}

	for n, jumpStats := range jumpPatterns {
		f, err := os.OpenFile("./patterns_"+strconv.Itoa(n), os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		var jumpPatternsArr []pair
		for k, v := range jumpStats {
			jumpPatternsArr = append(jumpPatternsArr, pair{k, v})
		}

		sort.Slice(jumpPatternsArr, func(i, j int) bool {
			return jumpPatternsArr[i].v > jumpPatternsArr[j].v
		})

		for _, kv := range jumpPatternsArr {
			fmt.Fprintf(f, "%s,%d,%.4f\n", kv.k, kv.v, float64(kv.v)/float64(jumps)*100)
		}

		f.Close()
	}
}

var stopping = []vm.OpCode{
	vm.RETURN,
	vm.REVERT,
	vm.SELFDESTRUCT,
	vm.STOP,
}

func (p *processor) isStop(op vm.OpCode, arg uint64, jumpDests map[uint64]struct{}, latest uint64) bool {
	if static.IsStop(op) {
		return true
	}
	if op == vm.JUMP || op == vm.JUMPI {
		if _, ok := jumpDests[arg]; !ok && arg > 2 && arg < latest {
			return true
		}
	}

	return false
}

func newExecutionPath(current []asm.Command) []asm.Command {
	n := make([]asm.Command, len(current))
	copy(n, current)
	return n
}

var solidityPrefix = []vm.OpCode{vm.PUSH1, vm.PUSH1, vm.MSTORE} // 0x60 0x60 0x52

func hasSolidityPrefix(it *asm.InstructionIterator) (string, bool) {
	prefix, ok := it.HasPrefix(solidityPrefix)
	return asm.OpsToString(prefix), ok
}

type toPut struct {
	bucket  []byte
	keys    [][]byte
	values  [][]byte
	idx     int
	isEmpty bool
	db      *ethdb.BoltDatabase
}

const BunchSize = 50000

func NewToPut(bucket []byte, size int, db *ethdb.BoltDatabase) *toPut {
	return &toPut{
		bucket:  bucket,
		keys:    make([][]byte, size),
		values:  make([][]byte, size),
		isEmpty: true,
		db:      db,
	}
}

func (p *toPut) Add(key, value []byte) {
	if p.idx >= len(p.keys)-1 {
		p.Flush()
	}

	if p.isEmpty {
		p.isEmpty = false
	}

	p.keys[p.idx] = key
	p.values[p.idx] = value
	p.idx++
}

func (p *toPut) Flush() error {
	if p.isEmpty {
		return nil
	}

	tuples := common.NewTuples(p.idx+1, 3, 1)
	for i := 0; i < p.idx; i++ {
		if err := tuples.Append(p.bucket, p.keys[i], p.values[i]); err != nil {
			return fmt.Errorf("tuples.Append failed: %w", err)
		}
	}
	sort.Sort(tuples)

	_, err := p.db.MultiPut(tuples.Values...)
	if err != nil {
		return fmt.Errorf("db.MultiPut failed: %w", err)
	}

	p.idx = 0
	p.isEmpty = true
	return nil
}
