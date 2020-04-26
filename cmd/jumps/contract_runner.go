package main

import (
	"fmt"
	"sort"

	"github.com/davecgh/go-spew/spew"

	"github.com/ledgerwatch/turbo-geth/core/asm"
	"github.com/ledgerwatch/turbo-geth/core/vm"
)

type ContractRunner struct {
	*Contract
	jumpi []int // all jumpi positions, sorted
}

func NewContractRunner(code []byte, debug bool) *ContractRunner {
	contract := NewContract(code)

	jumpis := codeOpCodesPositions(code, vm.JUMPI)

	if debug {
		fmt.Println("JUMPIs")
		for i, jump := range jumpis {
			fmt.Printf("%d: %[2]d %[2]x\n", i, jump)
		}
		spew.Dump(asm.PrintDisassembledBytes(code))
	}

	return &ContractRunner{
		contract,
		jumpis,
	}
}

func (c *ContractRunner) Run() error {
	if len(c.Code) == 0 {
		return nil
	}

	var (
		operation Command
		op        vm.OpCode    // current opcode
		stack     = newstack() // local stack
		pc        = uint64(0)  // program counter
		res       []byte       // result of the opcode execution function
		err       error
	)

	for {
		op = c.GetOp(pc)
		operation = Commands[op]

		if op == vm.JUMPI {
			found := sort.SearchInts(c.jumpi, int(pc))
			if found != len(c.jumpi) {
				panic(fmt.Sprintln(found, c.jumpi[found], pc))
			}
		}

		res, err = operation.execute(&pc, nil, c.Contract, nil, stack)
		_ = res

		if err != nil {
			return err
		}
		if op != vm.JUMP && op != vm.JUMPI {
			pc++
		}
	}

	return nil
}
