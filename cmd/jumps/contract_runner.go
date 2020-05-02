package static

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/asm"
	"github.com/ledgerwatch/turbo-geth/core/vm"
)

type ContractRunner struct {
	*Contract
	Jumpi []int // all jumpi positions, sorted
	debug bool

	errPath Path
}

// slice of sorted by PC Paths
// we don't jump as default
type Paths []Path

func (ps *Paths) Add(p Path, debug bool) {
	if len(p) == 0 {
		return
	}

	*ps = append(*ps, p)

	if debug {
		fmt.Println("add:", p)
		fmt.Println()
	}
}

func (ps Paths) Last() Path {
	return ps[len(ps)-1]
}

func (ps Paths) Next(debug bool) (Path, bool) {
	if len(ps) == 0 {
		return Path{}, false
	}

	last := ps.Last()
	if len(ps) == 1 {
		return last, true
	}

	// last.RemoveCycle()

	if debug {
		fmt.Println("last", last)
	}

	nextPath, wasChanged := last.Next()

	if debug {
		fmt.Println("next", nextPath)
	}

	return nextPath, wasChanged
}

type Path []Step

func (p Path) Next() (Path, bool) {
	nextPath := Path(make([]Step, len(p)))
	var pathStep Step
	var wasChanged bool
	var wasChangedIdx int
	for i := len(p) - 1; i >= 0; i-- {
		pathStep = p[i]
		if !pathStep.Jumped && !wasChanged {
			pathStep.Jumped = true
			wasChanged = true
			wasChangedIdx = i
		}
		nextPath[i] = pathStep
	}
	nextPath = nextPath[:wasChangedIdx+1]

	return nextPath, wasChanged
}

func (p *Path) RemoveCycle() bool {
	removed := p.RemoveTrailingCycles()

	// remove repeats. save "jumped" state
	removedDupl := p.Deduplicate()

	/*
	if removed || removedDupl {
		fmt.Println("!!!!!!!!!!!!!!!!!!", removed || removedDupl, removed, removedDupl)
	}
	*/

	return removed || removedDupl
}

func (p *Path) RemoveTrailingCycles() bool {
	if len(*p) < 3 {
		// 3 is the min code size with 2 cycles and an additional step
		return false
	}

	// remove loops starting the end
	var st Step
	possibleCycle := Path(make([]Step, 0, 10))

	tail := (*p)[len(*p)-1]
	possibleCycle = append(possibleCycle, tail)

	// fixme: it's not all the cases, but we can try dynamic programming later
	for i := len(*p) - 2; i >= len(*p)/2; i-- {
		st = (*p)[i]

		if tail.JumpPC == st.JumpPC {
			break
		}

		possibleCycle = append(possibleCycle, st)
	}

	// fmt.Println("!1", len(possibleCycle), len(*p)/2, possibleCycle)

	var hasCycle bool
	if len(possibleCycle) != 0 {
		// to form a cycle it should be Jumped=true
		possibleCycle[0].Jumped = true

		fromIdx := len(*p) - 1
		var toIdx int

		n := 0
		for i := len(*p) - 1 - len(possibleCycle); i >= 0; i-- {
			j := n % len(possibleCycle)

			if (*p)[i].JumpPC != possibleCycle[j].JumpPC {
				break
			}

			toIdx = i
			n++
		}

		//fmt.Println("N=", n, len(possibleCycle))
		if n < len(possibleCycle) {
			// fmt.Println("N=", n, len(possibleCycle))
			return false
		}

		hasCycle = fromIdx != toIdx
		if hasCycle {
			//fmt.Println("!!!", fromIdx, toIdx, len(*p), possibleCycle)

			*p = (*p)[:toIdx]

			toIdx = 0
			//fmt.Println("possible", possibleCycle)
			for i := 0; i < len(possibleCycle); i++ {
				if !possibleCycle[i].Jumped {
					//fmt.Println("ADD", possibleCycle[i], i)
					break
				} else {
					toIdx++
				}
			}
			possibleCycle = possibleCycle[toIdx:]

			for i := len(possibleCycle) - 1; i >= 0; i-- {
				*p = append(*p, possibleCycle[i])
			}

			(*p)[len(*p)-1].Jumped = true
		}
	}

	return hasCycle
}

func (p *Path) Deduplicate() bool {
	if len(*p) <= 1 {
		return false
	}

	toIdx := len(*p) - 1
	last := (*p)[toIdx]
	var hasCycle bool
	for i := len(*p) - 2; i >= 0; i-- {
		if last.JumpPC == (*p)[i].JumpPC {
			hasCycle = true
			if (*p)[i].Jumped {
				last.Jumped = true
			}
		} else {
			if hasCycle {
				*p = append((*p)[:i+1], (*p)[toIdx:]...)
			}
			last = (*p)[i]
			toIdx = i
			hasCycle = false
		}
	}
	if hasCycle {
		*p = (*p)[toIdx:]
	}

	return hasCycle
}

type Step struct {
	JumpPC int
	Jumped bool
}

var ErrUnknownOpcode = errors.New("unknown opcode")
var ErrMaxJumps = errors.New("max jumps or recursion")

const maxJumps = 1000 // works the same as 1_000_000

func NewContractRunner(code []byte, codeAddress common.Address, debug bool) *ContractRunner {
	contract := NewContract(code)
	contract.CodeAddr = &codeAddress

	jumpis := codeOpCodesPositions(code, vm.JUMPI)

	/*
		if debug {
			fmt.Println("JUMPIs")
			for i, jump := range jumpis {
				fmt.Printf("%d: %[2]d %[2]x\n", i, jump)
			}
			spew.Dump(asm.PrintDisassembledBytes(code))
		}
	*/

	return &ContractRunner{
		contract,
		jumpis,
		debug,
		Path{},
	}
}

func (c *ContractRunner) Run(ctx context.Context) (int, error) {
	if len(c.Code) == 0 {
		return 0, nil
	}

	jumpPaths := &Paths{}

	if len(c.Jumpi) > 0 {
		jumpPaths.Add(Path{
			{c.Jumpi[0], false},
		}, c.debug)
	}

	firstRun := true
	paths := 0
	var err error

pathsLoop:
	for {
		select {
		case <-ctx.Done():
			err = ErrTimeout
			break
		default:
			// nothing to do
		}

		var fullPath []string
		if c.debug {
			fmt.Println("======================================================================\n\n\n")
		}

		currentCodePath, ok := jumpPaths.Next(c.debug)
		if !ok && !firstRun {
			break
		}

		jumps := 0
		firstRun = false
		paths++

		gotPath := Path(make([]Step, 0, len(currentCodePath)))

		var (
			operation Command
			op        vm.OpCode    // current opcode
			stack     = newstack() // local stack
			pc        = uint64(0)  // program counter
		)

		jumpiIdx := -1
		onThePath := true

		innerCtx, innerCancel := context.WithTimeout(context.Background(), 3*time.Second)

		var lastDestPC uint64

		for {
			select {
			case <-ctx.Done():
				err = ErrTimeout

				break pathsLoop
			case <-innerCtx.Done():
				innerCancel()
				continue pathsLoop
			default:
				// nothing to do
			}

			op = c.GetOp(pc)
			operation = Commands[op]

			if c.debug {
				fullPath = append(fullPath, asm.PrintCommand(op, pc, nil))
			}

			// we don't count consumed gas so we need something to stop
			if op == vm.JUMP || op == vm.JUMPI {
				if jumps >= maxJumps {
					// too much to analyze
					jumpPaths.Add(gotPath, c.debug)

					// return errMaxJumps
					innerCancel()
					continue pathsLoop
				}
				jumps++
			}

			if op == vm.JUMPI {
				jumpiIdx++

				found := sort.SearchInts(c.Jumpi, int(pc))
				if found >= len(c.Jumpi) {
					innerCancel()
					panic(fmt.Sprintln(found, c.Jumpi[found], pc))
				}

				if jumpiIdx >= len(currentCodePath) {
					onThePath = false
				}

				// default case
				operation.execute = opJumpiNotJUMP
				var willJump bool

				if onThePath {
					currentJumpi := currentCodePath[jumpiIdx]

					if currentJumpi.JumpPC == int(pc) {
						if currentJumpi.Jumped {
							willJump = true
							operation.execute = opJumpiJUMP
						}
					}
				} else {
					onThePath = false
				}

				gotPath = append(gotPath, Step{int(pc), willJump})
			}

			if operation.execute == nil {
				jumpPaths.Add(gotPath, c.debug)

				innerCancel()
				continue pathsLoop
			}

			prevPc := pc
			_, err = operation.execute(&pc, nil, c.Contract, nil, stack)
			if op != vm.JUMP && op != vm.JUMPI {
				pc++
			}

			nextOp := c.GetOp(pc)
			if nextOp == vm.JUMPDEST {
				if lastDestPC == pc {
					// it's a cycle
					gotPath[len(gotPath)-1].Jumped = true
					jumpPaths.Add(gotPath, c.debug)
					innerCancel()
					continue pathsLoop
				} else {
					lastDestPC = pc
				}
			}

			if err == nil && pc == prevPc {
				// it's a self-loop
				gotPath.RemoveCycle()

				if len(gotPath) > 0 {
					for i := len(gotPath) - 1; i >= 0; i-- {
						if !gotPath[i].Jumped {
							gotPath[i].Jumped = true
							break
						}
					}
				}
				jumpPaths.Add(gotPath, c.debug)

				innerCancel()
				continue pathsLoop
			} else if err == nil && prevPc > pc && nextOp == vm.JUMPDEST {
				// cycle
				hasCycle := gotPath.RemoveCycle()

				if hasCycle {
					jumpPaths.Add(gotPath, c.debug)
					innerCancel()
					continue pathsLoop
				}
			}

			if err != nil {
				if len(gotPath) == 0 {
					// can't reach even first jumpi

					innerCancel()

					break pathsLoop
				}

				jumpPaths.Add(gotPath, c.debug)

				// if we found a non-static jump we can stop
				if errors.Is(err, ErrNonStatic) {
					innerCancel()

					c.errPath = gotPath
					if c.debug {
						fmt.Printf("Full path\n%v\n\n", fullPath)
					}

					break pathsLoop
				}
				if errors.Is(err, ErrNoValueStatic) {
					innerCancel()

					c.errPath = gotPath

					break pathsLoop
				}

				innerCancel()
				continue pathsLoop
			}
		}
	}

	return paths, err
}

func (c ContractRunner) ErrPath() Path {
	return c.errPath
}