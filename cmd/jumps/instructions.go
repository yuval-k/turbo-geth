package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/davecgh/go-spew/spew"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/vm"
)

var (
	errInvalidJump           = errors.New("evm: invalid jump destination")
	errNonStatic             = errors.New("non static jump")
	errNoValueStatic             = errors.New("no value")
	errReturn = errors.New("op.RETURN")
	errRevert = errors.New("op.REVERT")
	errSelfDestruct = errors.New("op.SELFDESTRUCT")
	errStop = errors.New("op.STOP")
)

func NotStaticIfOneNotStatic(cmds ...*cell) bool {
	for _, cmd := range cmds {
		if !cmd.static {
			return false
		}
	}
	return true
}

func opAdd(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y)

	return nil, nil
}

func opSub(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y)

	return nil, nil
}

func opMul(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	y.static = NotStaticIfOneNotStatic(x, y)

	return nil, nil
}

// fixme: UNSAFE!!! Possibly it's wrong to use NotStaticIfOneNotStatic. Check original opDiv
func opDiv(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y)

	return nil, nil
}

func opSdiv(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	res := NewCell(NotStaticIfOneNotStatic(x, y))
	stack.push(res)

	return nil, nil
}

func opMod(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	x.static = NotStaticIfOneNotStatic(x, y)
	stack.push(x)

	return nil, nil
}

func opSmod(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	res := NewCell(NotStaticIfOneNotStatic(x, y))
	stack.push(res)

	return nil, nil
}

func opExp(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	base, exponent := stack.pop(), stack.pop()
	stack.push(NewCell(NotStaticIfOneNotStatic(base, exponent)))

	return nil, nil
}

// fixme: isStatic depends on the code and data
func opSignExtend(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	/*
		back := stack.pop()
		if back.Cmp(big.NewInt(31)) < 0 {
			bit := uint(back.Uint64()*8 + 7)
			num := stack.pop()
			mask := back.Lsh(common.Big1, bit)
			mask.Sub(mask, common.Big1)
			if num.Bit(int(bit)) > 0 {
				num.Or(num, mask.Not(mask))
			} else {
				num.And(num, mask)
			}

			stack.push(math.U256(num))
		}
	*/

	stack.pop()
	return nil, nil
}

func opNot(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, _ *Stack) ([]byte, error) {
	return nil, nil
}

func opLt(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y) // it could be TRUE, but let's decide a bit more strict

	return nil, nil
}

func opGt(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y) // it could be TRUE, but let's decide a bit more strict

	return nil, nil
}

func opSlt(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()

	x.static = NotStaticIfOneNotStatic(x, y) // it could be TRUE, but let's decide a bit more strict
	y.static = NotStaticIfOneNotStatic(x, y) // it could be TRUE, but let's decide a bit more strict

	return nil, nil
}

func opSgt(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()

	x.static = NotStaticIfOneNotStatic(x, y) // it could be TRUE, but let's decide a bit more strict
	y.static = NotStaticIfOneNotStatic(x, y) // it could be TRUE, but let's decide a bit more strict

	return nil, nil
}

func opEq(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y) // it could be TRUE, but let's decide a bit more strict

	return nil, nil
}

func opIszero(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, _ *Stack) ([]byte, error) {
	return nil, nil
}

func opAnd(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	x.static = NotStaticIfOneNotStatic(x, y)
	stack.push(x)

	return nil, nil
}

func opOr(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y)
	stack.push(y)

	return nil, nil
}

func opXor(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.static = NotStaticIfOneNotStatic(x, y)
	stack.push(y)

	return nil, nil
}

func opByte(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	th, val := stack.pop(), stack.peek()
	val.static = NotStaticIfOneNotStatic(th, val) // it could be TRUE, but let's decide a bit more strict

	return nil, nil
}

func opAddmod(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.pop()

	x.static = NotStaticIfOneNotStatic(x, y, z) // it could be (x, z), but let's decide a bit more strict

	return nil, nil
}

func opMulmod(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	x.static = NotStaticIfOneNotStatic(x, y, z) // it could be (x, z), but let's decide a bit more strict

	return nil, nil
}

// opSHL implements Shift Left
// The SHL instruction (shift left) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the left by arg1 number of bits.
func opSHL(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := stack.pop(), stack.peek()

	value.static = NotStaticIfOneNotStatic(shift, value) // it could be SAME, but let's decide a bit more strict

	return nil, nil
}

// opSHR implements Logical Shift Right
// The SHR instruction (logical shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with zero fill.
func opSHR(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := stack.pop(), stack.peek()
	value.static = NotStaticIfOneNotStatic(shift, value) // it could be SAME, but let's decide a bit more strict
	return nil, nil
}

// opSAR implements Arithmetic Shift Right
// The SAR instruction (arithmetic shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with sign extension.
func opSAR(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Note, S256 returns (potentially) a new bigint, so we're popping, not peeking this one
	shift, value := stack.pop(), stack.pop()
	value.static = NotStaticIfOneNotStatic(shift, value) // it could be SAME, but let's decide a bit more strict
	stack.push(value)
	return nil, nil
}

func opSha3(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()

	stack.push(NewNonStaticCell()) // fixme it's stricter than it could be

	return nil, nil
}

func opAddress(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewStaticCell())
	return nil, nil
}

func opBalance(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	slot := stack.peek()
	slot.static = false
	return nil, nil
}

func opOrigin(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opCaller(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opCallValue(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opCallDataLoad(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.push(NewNonStaticCell())

	return nil, nil
}

func opCallDataSize(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())

	return nil, nil
}

func opCallDataCopy(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	stack.pop()

	return nil, nil
}

func opReturnDataSize(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell()) // fixme: stricter than it could be

	return nil, nil
}

func opReturnDataCopy(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	stack.pop()

	return nil, nil
}

func opExtCodeSize(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	slot := stack.peek()
	slot.static = false // fixme: stricter than it could be

	return nil, nil
}

func opCodeSize(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewStaticCell())

	return nil, nil
}

func opCodeCopy(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	stack.pop()

	return nil, nil
}

func opExtCodeCopy(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()

	return nil, nil
}

// opExtCodeHash returns the code hash of a specified account.
// There are several cases when the function is called, while we can relay everything
// to `state.GetCodeHash` function to ensure the correctness.
//   (1) Caller tries to get the code hash of a normal contract account, state
// should return the relative code hash and set it as the result.
//
//   (2) Caller tries to get the code hash of a non-existent account, state should
// return common.Hash{} and zero will be set as the result.
//
//   (3) Caller tries to get the code hash for an account without contract code,
// state should return emptyCodeHash(0xc5d246...) as the result.
//
//   (4) Caller tries to get the code hash of a precompiled account, the result
// should be zero or emptyCodeHash.
//
// It is worth noting that in order to avoid unnecessary create and clean,
// all precompile accounts on mainnet have been transferred 1 wei, so the return
// here should be emptyCodeHash.
// If the precompile account is not transferred any amount on a private or
// customized chain, the return value will be zero.
//
//   (5) Caller tries to get the code hash for an account which is marked as suicided
// in the current transaction, the code hash of this account should be returned.
//
//   (6) Caller tries to get the code hash for an account which is marked as deleted,
// this account should be regarded as a non-existent account and zero should be returned.
func opExtCodeHash(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	slot := stack.peek()
	slot.static = true
	return nil, nil
}

func opGasprice(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opBlockhash(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()

	stack.push(NewNonStaticCell())

	return nil, nil
}

func opCoinbase(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opTimestamp(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opNumber(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opDifficulty(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opGasLimit(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opPop(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	return nil, nil
}

func opMload(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	v := stack.peek()
	v.static = false // fixme: not true if we introduce momory type
	return nil, nil
}

func opMstore(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()

	return nil, nil
}

func opMstore8(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	return nil, nil
}

func opSload(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	loc := stack.peek()
	loc.static = false // fixme: not true if we introduce momory type
	return nil, nil
}

func opSstore(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()

	return nil, nil
}

func opJump(pc *uint64, _ *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos := stack.pop()
	if !pos.static || !pos.IsValue() {
		return nil, fmt.Errorf("%w on %v", errNonStatic, spew.Sdump(pc))
	}
	if !contract.validJumpdest(pos.v) {
		return nil, fmt.Errorf("%w on %v", errInvalidJump, spew.Sdump(pc))
	}
	*pc = pos.v.Uint64()

	return nil, nil
}

func opJumpi(pc *uint64, _ *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos, cond := stack.pop(), stack.pop()

	if !pos.static {
		return nil, fmt.Errorf("jumpi: %w on %v", errNonStatic, spew.Sdump(pc))
	}
	if pos.static && !pos.IsValue() {
		return nil, fmt.Errorf("jumpi: %w on %v", errNoValueStatic, spew.Sdump(pc))
	}

	if cond.v.Sign() != 0 {
		if !contract.validJumpdest(pos.v) {
			return nil, errInvalidJump
		}
		*pc = pos.v.Uint64()
	} else {
		*pc++
	}

	return nil, nil
}

func opJumpiJUMP(pc *uint64, _ *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos, _ := stack.pop(), stack.pop()

	if !pos.static {
		return nil, fmt.Errorf("jumpi: %w on %v", errNonStatic, spew.Sdump(pc))
	}
	if pos.static && !pos.IsValue() {
		return nil, fmt.Errorf("jumpi: %w on %v", errNoValueStatic, spew.Sdump(pc))
	}

	if !contract.validJumpdest(pos.v) {
		return nil, fmt.Errorf("%w on %v", errInvalidJump, spew.Sdump(pc))
	}
	*pc = pos.v.Uint64()

	return nil, nil
}

func opJumpiNotJUMP(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos, _ := stack.pop(), stack.pop()

	if !pos.static {
		return nil, fmt.Errorf("jumpi: %w on %v", errNonStatic, spew.Sdump(pc))
	}
	if pos.static && !pos.IsValue() {
		return nil, fmt.Errorf("jumpi: %w on %v", errNoValueStatic, spew.Sdump(pc))
	}

	*pc++

	return nil, nil
}

func opJumpdest(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	return nil, nil
}

func opPc(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewStaticCell())
	return nil, nil
}

func opMsize(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opGas(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

func opCreate(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	stack.pop()

	stack.push(NewNonStaticCell())

	return nil, nil
}

func opCreate2(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()

	stack.push(NewNonStaticCell())
	return nil, nil
}

func opCall(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas in interpreter.evm.callGasTemp.
	stack.pop()

	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()

	stack.push(NewNonStaticCell())

	return nil, nil
}

func opCallCode(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	stack.pop()

	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()

	stack.push(NewNonStaticCell())

	return nil, nil
}

func opDelegateCall(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	stack.pop()

	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()

	stack.push(NewNonStaticCell())

	return nil, nil
}

func opStaticCall(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	stack.pop()

	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()
	stack.pop()

	stack.push(NewNonStaticCell())

	return nil, nil
}

func opReturn(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()
	
	return nil, errReturn
}

func opRevert(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	stack.pop()

	return nil, errRevert
}

func opStop(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	return nil, errStop
}

func opSuicide(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.pop()
	return nil, errSelfDestruct
}

// following functions are used by the instruction jump  table
type executionFunc func(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, memory *vm.Memory, stack *Stack) ([]byte, error)

// make log instruction function
func makeLog(size int) executionFunc {
	return func(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		stack.pop()
		stack.pop()
		for i := 0; i < size; i++ {
			stack.pop()
		}

		return nil, nil
	}
}

var PushDest = make(map[uint64]struct{})

// opPush1 is a specialized version of pushN
func opPush1(pc *uint64, _ *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	var (
		codeLen = uint64(len(contract.Code))
		integer = big.NewInt(0)
	)
	*pc += 1
	c := NewStaticCell()
	if *pc < codeLen {
		c.SetValue(integer.SetUint64(uint64(contract.Code[*pc])))

		PushDest[uint64(contract.Code[*pc])] = struct{}{}
	} else {
		c.SetValue(integer.SetUint64(0))
	}
	stack.push(c)

	return nil, nil
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, _ *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		codeLen := len(contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := big.NewInt(0)
		integer.SetBytes(common.RightPadBytes(contract.Code[startMin:endMin], pushByteSize))

		PushDest[integer.Uint64()] = struct{}{}

		c := NewStaticCell()
		c.SetValue(integer)
		stack.push(c)

		*pc += size
		return nil, nil
	}
}

// make dup instruction function
func makeDup(size int64) executionFunc {
	return func(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		stack.dup(int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size++
	return func(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		stack.swap(int(size))
		return nil, nil
	}
}

func opSelfBalance(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewNonStaticCell())
	return nil, nil
}

// opChainID implements CHAINID opcode
func opChainID(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	stack.push(NewStaticCell())
	return nil, nil
}