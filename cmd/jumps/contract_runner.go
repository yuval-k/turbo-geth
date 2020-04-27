package main

import (
	"errors"
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

// slice of sorted by PC paths
// we don't jump as default
type paths []path

func (ps *paths) add(p path) {
	sort.SliceStable(p, func(i, j int) bool {
		return p[i].jumpPC < p[j].jumpPC
	})

	*ps = append(*ps, p)
}

func (ps paths) last() path {
	return ps[len(ps)-1]
}

func (ps paths) next() (path, bool) {
	if len(ps) == 0 {
		return path{}, false
	}
	if len(ps) == 1 {
		return ps.last(), true
	}

	last := ps[len(ps)-1]
	nextPath := make([]step, len(last))
	var pathStep step
	var wasChanged bool
	var wasChangedIdx int
	for i := len(last) - 1; i >= 0; i-- {
		pathStep = last[i]
		if !pathStep.jumped && !wasChanged {
			pathStep.jumped = true
			wasChanged = true
			wasChangedIdx = i
		}
		nextPath[i] = pathStep
	}

	nextPath = nextPath[:wasChangedIdx+1]

	return nextPath, wasChanged
}

type path []step

type step struct {
	jumpPC int
	jumped bool
}

var errUnknownOpcode = errors.New("unknown opcode")
var errMaxJumps = errors.New("max jumps or recursion")

const maxJumps = 100 // works the same as 1_000_000

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

	jumps := 0
	jumpPaths := &paths{}

	if len(c.jumpi) > 0 {
		jumpPaths.add(path{
			{c.jumpi[0], false},
		})
	}

	firstRun := true
	paths := 0

	pathsLoop:
	for currentCodePath, ok := jumpPaths.next(); ok || firstRun; {
		firstRun = false
		paths++

		var (
			operation Command
			op        vm.OpCode    // current opcode
			stack     = newstack() // local stack
			pc        = uint64(0)  // program counter
			//res       []byte       // result of the opcode execution function
			err error
		)

		jumpiIdx := -1
		onThePath := true
		for {
			op = c.GetOp(pc)
			operation = Commands[op]

			// we don't count consumed gas so we need something to stop
			if op == vm.JUMP || op == vm.JUMPI {
				if jumps >= maxJumps {
					// too much to analize
					return errMaxJumps
				}
				jumps++
			}

			if op == vm.JUMPI {
				//fmt.Println(paths, jumpiIdx, ok, firstRun)

				jumpiIdx++

				found := sort.SearchInts(c.jumpi, int(pc))
				if found >= len(c.jumpi) {
					panic(fmt.Sprintln(found, c.jumpi[found], pc))
				}

				if len(currentCodePath) >= jumpiIdx+1 {
					onThePath = false
				}

				// default case
				operation.execute = opJumpiNotJUMP

				if onThePath {
					currentJumpi := currentCodePath[jumpiIdx]
					if currentJumpi.jumpPC == int(pc) {
						if currentJumpi.jumped {
							operation.execute = opJumpiJUMP
						}
					}
				} else {
					onThePath = false
				}
			}

			if operation.execute == nil {
				continue pathsLoop
				return errUnknownOpcode
				/*
					return fmt.Errorf("%w: operation is nil PC %[2]d(%[2]x), OP %d\nCode hash:\n%v",
						errUnknownOpcode,
						pc,
						op,
						c.CodeHash.String())
				*/
			}

			_, err = operation.execute(&pc, nil, c.Contract, nil, stack)
			if err != nil {
				continue pathsLoop
				return err
			}

			if op != vm.JUMP && op != vm.JUMPI {
				pc++
			}
		}
	}

	fmt.Println("total paths", paths)

	return nil
}
