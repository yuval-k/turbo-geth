package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core/asm"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

func main() {
	jumps()
	//cleanup()
	//cleanupBeggining()
	//cleanupErrors()
	// begginings()
}

type command struct {
	pc  uint64
	op  vm.OpCode
	arg []byte
}

func jumps() {
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

func cleanupBeggining() {
	db, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbCleaned, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes_cleaned_begginings")
	if err != nil {
		log.Fatal(err)
	}
	defer dbCleaned.Close()

	var incorrectBeggining int

	i := 0
	expected := [3]string{"PUSH1", "PUSH1", "MSTORE"}
	beggining := [3]string{}

	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		it := asm.NewInstructionIterator(v)

		n := 0

	loop:
		for it.Next() {
			if n <= 2 {
				it.Error()
				beggining[n] = it.Op().String()
			}
			if n == 2 {
				if expected != beggining {
					i++
					incorrectBeggining++
					return true, nil
				}
			}
			if n >= 3 {
				break loop
			}
			n++
		}

		dbCleaned.Put(dbutils.CodeBucket, k, v)
		i++
		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}

	fmt.Println("different contracts", i)
	fmt.Println("incorrect begginings", incorrectBeggining)
}

func cleanupErrors() {
	db, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes_cleaned_begginings")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbCleaned, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes_cleaned_begginings_no_errors")
	if err != nil {
		log.Fatal(err)
	}
	defer dbCleaned.Close()

	var incorrectContract int

	i := 0

	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		it := asm.NewInstructionIterator(v)

		for it.Next() {
			if it.Error() != nil {
				i++
				incorrectContract++
				return true, nil
			}
		}

		dbCleaned.Put(dbutils.CodeBucket, k, v)
		i++
		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}

	fmt.Println("different contracts", i)
	fmt.Println("incorrect contract", incorrectContract)
}

func cleanup() {
	db, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbCleaned, err := ethdb.NewBoltDatabase("/mnt/sdb/contract_codes_cleaned")
	if err != nil {
		log.Fatal(err)
	}

	defer dbCleaned.Close()

	incorrectCodes := make(map[vm.OpCode]int)

	i := 0

	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		i++

		it := asm.NewInstructionIterator(v)
		var previous *command

		for it.Next() {
			if previous == nil {
				// first call
			}

			op := it.Op()

			if op.IsUnknown() {
				incorrectCodes[op]++
				return true, nil
			}
		}

		dbCleaned.Put(dbutils.CodeBucket, k, v)
		return true, nil
	})
	if err != nil {
		log.Println("err", err)
	}

	fmt.Println("different contracts", i)

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
	err = db.Walk(dbutils.CodeBucket, make([]byte, 32), 0, func(k, v []byte) (bool, error) {
		it := asm.NewInstructionIterator(v)

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
