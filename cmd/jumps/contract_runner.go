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
	if len(p) == 0 {
		return
	}

	fmt.Println("path", p)

	/*
	sort.SliceStable(p, func(i, j int) bool {
		return p[i].jumpPC < p[j].jumpPC
	})
	*/

	/*
	var lastStep step
	var toCut int
	for i := len(p) - 1; i >= 0; i-- {
		if i == len(p) - 1 {
			lastStep = p[i]
		} else {
			if p[i].jumpPC == lastStep.jumpPC {
				p[i] = lastStep
				toCut++
			} else {
				break
			}
		}
	}
	p = p[:len(p)-toCut]
	*/

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

	nextPath, wasChanged := last.next()
	//fmt.Println("!!! 1", len(last), len(nextPath), nextPath.hasCycle())
	nextPath.removeCycle()
	nextPath, wasChanged = last.next()
	//fmt.Println("!!! 2", len(last), len(nextPath), nextPath.hasCycle())
	//fmt.Println()

	/*
	spew.Println(last)
	spew.Println(nextPath)
	 */

	return nextPath, wasChanged
}

type path []step

func (p path) hasCycle() bool {
	m := make(map[int]struct{}, len(p))
	var ok bool
	for _, st := range p {
		if _, ok = m[st.jumpPC]; ok {
			return true
		}
		m[st.jumpPC] = struct{}{}
	}

	return false
}

func (p path) next() (path, bool) {
	nextPath := path(make([]step, len(p)))
	var pathStep step
	var wasChanged bool
	var wasChangedIdx int
	for i := len(p) - 1; i >= 0; i-- {
		pathStep = p[i]
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

func (p *path) removeCycle() {
	// remove repeats. save "jumped" state

	// remove loops starting the end
	m := make(map[int]struct{}, len(*p))
	var ok bool
	var cycleHead int
	for _, st := range *p {
		if _, ok = m[st.jumpPC]; ok {
			cycleHead = st.jumpPC
			break
		}
		m[st.jumpPC] = struct{}{}
	}

	if ok {
		for i, st := range *p {
			if st.jumpPC == cycleHead {
				*p = (*p)[:i+1]
				break
			}
		}
	}

	return
}

type step struct {
	jumpPC int
	jumped bool
}

var errUnknownOpcode = errors.New("unknown opcode")
var errMaxJumps = errors.New("max jumps or recursion")

const maxJumps = 1000 // works the same as 1_000_000

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

	jumpPaths := &paths{}

	if len(c.jumpi) > 0 {
		jumpPaths.add(path{
			{c.jumpi[0], false},
		})
	}

	firstRun := true
	paths := 0

	pathsLoop:
	for {
		currentCodePath, ok := jumpPaths.next()
		if !ok && !firstRun {
			break
		}

		jumps := 0
		firstRun = false
		paths++

		//fmt.Println("111", paths, spew.Sdump(currentCodePath))

		gotPath := path(make([]step, 0, len(currentCodePath)))

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

			//fmt.Println(pc, op.String())

			// we don't count consumed gas so we need something to stop
			if op == vm.JUMP || op == vm.JUMPI {
				if jumps >= maxJumps {
					// too much to analyze
					jumpPaths.add(gotPath)
					spew.Println("path current - too much to analize", pc, op.String(), paths, len(c.jumpi), currentCodePath, gotPath)
					// return errMaxJumps
					continue pathsLoop
				}
				jumps++
			}

			if op == vm.JUMPI {
				jumpiIdx++

				found := sort.SearchInts(c.jumpi, int(pc))
				if found >= len(c.jumpi) {
					panic(fmt.Sprintln(found, c.jumpi[found], pc))
				}

				if jumpiIdx >= len(currentCodePath)  {
					onThePath = false
				}

				// default case
				operation.execute = opJumpiNotJUMP
				var willJump bool
				if onThePath {
					currentJumpi := currentCodePath[jumpiIdx]

					if currentJumpi.jumpPC == int(pc) {
						if currentJumpi.jumped {
							willJump = true
							operation.execute = opJumpiJUMP
						}
					}
				} else {
					onThePath = false
				}

				gotPath = append(gotPath, step{int(pc), willJump})
			}

			if operation.execute == nil {
				jumpPaths.add(gotPath)
				//spew.Println("path current - operation.execute == nil", paths, currentCodePath, gotPath)
				continue pathsLoop
				//return errUnknownOpcode
				/*
					return fmt.Errorf("%w: operation is nil PC %[2]d(%[2]x), OP %d\nCode hash:\n%v",
						errUnknownOpcode,
						pc,
						op,
						c.CodeHash.String())
				*/
			}

			prevPc := pc
			_, err = operation.execute(&pc, nil, c.Contract, nil, stack)
			if op != vm.JUMP && op != vm.JUMPI {
				pc++
			}

			if pc == prevPc {
				fmt.Println("!!!!!!!!!!!", pc, op.String())
				// it's a self-loop
				gotPath.removeCycle()

				if len(gotPath) > 0 {
					for i := len(gotPath) - 1; i >= 0; i-- {
						if !gotPath[i].jumped {
							gotPath[i].jumped = true
							break
						}
					}
				}
				jumpPaths.add(gotPath)
				continue pathsLoop
			}

			if err != nil {
				if len(gotPath) == 0 {
					// can't reach even first jumpi
					break pathsLoop
				}
				jumpPaths.add(gotPath)
				//spew.Println("path current - err", paths, currentCodePath, gotPath, err)
				continue pathsLoop
				//return err
			}
		}
	}

	fmt.Println("total paths", paths)

	return nil
}
