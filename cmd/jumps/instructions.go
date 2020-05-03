package static

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/davecgh/go-spew/spew"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/math"
	"github.com/ledgerwatch/turbo-geth/core/vm"
)

var (
	ErrInvalidJump   = errors.New("evm: invalid jump destination")
	ErrNonStatic     = errors.New("non static jump")
	ErrNoValueStatic = errors.New("no value")
	ErrReturn        = errors.New("op.RETURN")
	ErrRevert        = errors.New("op.REVERT")
	ErrSelfDestruct  = errors.New("op.SELFDESTRUCT")
	ErrStop          = errors.New("op.STOP")

	ErrTimeout = errors.New("execution timeout")

	tt255   = math.BigPow(2, 255)
	bigZero = new(big.Int)
)

func NotStaticIfOneNotStatic(cmds ...*cell) bool {
	for _, cmd := range cmds {
		if !cmd.static {
			return false
		}
	}
	return true
}

func opAdd(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	y.static = NotStaticIfOneNotStatic(x, y)

	if y.static && y.IsValue() && x.IsValue() {
		math.U256(y.v.Add(x.v, y.v))
	} else {
		y.unset(interpreter)
	}

	y.AddHistory(vm.ADD, *pc, y.static)

	return nil, nil
}

func opSub(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	y.static = NotStaticIfOneNotStatic(x, y)

	if y.static && y.IsValue() && x.IsValue() {
		math.U256(y.v.Sub(x.v, y.v))
	} else {
		y.unset(interpreter)
	}

	y.AddHistory(vm.SUB, *pc, y.static)

	return nil, nil
}

func opMul(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	x.static = NotStaticIfOneNotStatic(x, y)

	if x.static && x.IsValue() && y.IsValue() {
		math.U256(x.v.Mul(x.v, y.v))
	} else {
		x.unset(interpreter)
	}
	stack.push(x)

	x.AddHistory(vm.MUL, *pc, x.static)

	return nil, nil
}

func opDiv(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if y.static && y.v.Sign() == 0 {
		y.set(0, interpreter)
		y.static = true
	} else {
		y.static = NotStaticIfOneNotStatic(x, y)
		if y.static && y.IsValue() && x.IsValue() {
			math.U256(y.v.Div(x.v, y.v))
		} else {
			y.unset(interpreter)
		}
	}

	y.AddHistory(vm.DIV, *pc, y.static)

	return nil, nil
}

func opSdiv(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	if x.static && x.IsValue() {
		x.v = math.U256(x.v)
	}
	if y.static && y.IsValue() {
		y.v = math.U256(y.v)
	}

	res := &cell{getInt(0, interpreter), NotStaticIfOneNotStatic(x, y), nil}

	ySign := y.Sign()
	xSign := x.Sign()
	if (ySign != nil && *ySign == 0) || (xSign != nil && *xSign == 0) {
		res.static = true
	} else {
		if xSign == nil || ySign == nil {
			res.static = false
			res.unset(interpreter)
		} else {
			if x.Sign() != y.Sign() {
				res.v.Div(x.v.Abs(x.v), y.v.Abs(y.v))
				res.v.Neg(res.v)
			} else {
				res.v.Div(x.v.Abs(x.v), y.v.Abs(y.v))
			}
		}

		if res.IsValue() {
			res.v = math.U256(res.v)
		}
	}

	res.AddHistory(vm.SDIV, *pc, res.static)
	stack.push(res)

	return nil, nil
}

func opMod(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	ySign := y.Sign()
	if ySign != nil && *ySign == 0 {
		x.set(0, interpreter)
		x.static = true
	} else {
		x.static = NotStaticIfOneNotStatic(x, y)
		if !x.static {
			x.unset(interpreter)
		}

		if x.static && x.IsValue() && y.IsValue() {
			math.U256(x.v.Mod(x.v, y.v))
		}
	}

	x.AddHistory(vm.MOD, *pc, x.static)
	stack.push(x)

	return nil, nil
}

func opSmod(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	res := &cell{getInt(0, interpreter), NotStaticIfOneNotStatic(x, y), nil}

	ySign := y.Sign()
	if ySign != nil && *ySign == 0 {
		res.static = true
	} else {
		if !res.static {
			res.unset(interpreter)
		} else {
			if x.IsValue() && y.IsValue() {
				xSign := x.Sign()
				if xSign != nil && *xSign < 0 {
					res.v.Mod(x.v.Abs(x.v), y.v.Abs(y.v))
					res.v.Neg(res.v)
				} else {
					res.v.Mod(x.v.Abs(x.v), y.v.Abs(y.v))
				}

				math.U256(res.v)
			}
		}
	}

	res.AddHistory(vm.SMOD, *pc, res.static)
	stack.push(res)

	return nil, nil
}

var one = big.NewInt(1)

func opExp(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	base, err := stack.pop()
	if err != nil {
		return nil, err
	}

	exponent, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(exponent.v)

	if exponent.static && exponent.IsValue() {
		cmpToOne := exponent.v.Cmp(one)

		if cmpToOne <= 0 {
			if cmpToOne < 0 { // Exponent is zero
				// x ^ 0 == 1
				base.set(1, interpreter)
				base.static = true
			}
			/*
				// nothing to do if cmpToOne == 0
				if cmpToOne == 0 { // Exponent is one
					// x ^ 1 == x
				}
			*/

			base.AddHistory(vm.EXP, *pc, base.static)
			stack.push(base)
			return nil, nil
		}
	}

	baseSign := base.Sign()
	if baseSign != nil && *baseSign == 0 {
		// 0 ^ y, if y != 0, == 0
		base.set(0, interpreter)
		base.static = true

		base.AddHistory(vm.EXP, *pc, base.static)
		stack.push(base)

		return nil, nil
	}

	res := NewNonStaticCell()

	exponentSign := exponent.Sign()
	if exponentSign != nil && baseSign != nil {
		res.v = math.Exp(base.v, exponent.v)
		res.static = true
	}

	res.AddHistory(vm.EXP, *pc, res.static)
	stack.push(res)

	interpreter.IntPool.Put(base.v)

	return nil, nil
}

var n31 = big.NewInt(31)

func opSignExtend(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	back, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(back.v)

	// fixme it might be that we should push a NonStatic value in else branch
	if back.static && back.IsValue() && back.v.Cmp(n31) < 0 {
		num, err := stack.pop()
		if err != nil {
			return nil, err
		}

		if !num.static || !num.IsValue() {
			num.unset(interpreter)
			num.static = false
		} else {
			bit := uint(back.v.Uint64()*8 + 7)

			mask := back.v.Lsh(common.Big1, bit)
			mask.Sub(mask, common.Big1)
			if num.v.Bit(int(bit)) > 0 {
				num.v.Or(num.v, mask.Not(mask))
			} else {
				num.v.And(num.v, mask)
			}

			math.U256(num.v)
			num.static = true
		}

		num.AddHistory(vm.EXP, *pc, num.static)
		stack.push(num)
	}

	return nil, nil
}

func opNot(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if x.static && x.IsValue() {
		math.U256(x.v.Not(x.v))
	} else {
		x.unset(interpreter)
		x.static = false
	}

	x.AddHistory(vm.EXP, *pc, x.static)

	return nil, nil
}

func opLt(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if x.static && x.IsValue() && y.static && y.IsValue() {
		if x.v.Cmp(y.v) < 0 {
			y.set(1, interpreter)
		} else {
			y.set(0, interpreter)
		}

		y.static = true
	} else {
		y.unset(interpreter)
		y.static = false
	}

	y.AddHistory(vm.LT, *pc, y.static)

	return nil, nil
}

func opGt(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if x.static && x.IsValue() && y.static && y.IsValue() {
		if x.v.Cmp(y.v) < 0 {
			y.set(1, interpreter)
		} else {
			y.set(0, interpreter)
		}

		y.static = true
	} else {
		y.unset(interpreter)
		y.static = false
	}

	y.AddHistory(vm.GT, *pc, y.static)

	return nil, nil
}

func opSlt(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if x.static && x.IsValue() && y.static && y.IsValue() {
		xSign := x.v.Cmp(tt255)
		ySign := y.v.Cmp(tt255)

		switch {
		case xSign >= 0 && ySign < 0:
			y.set(1, interpreter)

		case xSign < 0 && ySign >= 0:
			y.set(0, interpreter)

		default:
			if x.v.Cmp(y.v) < 0 {
				y.set(1, interpreter)
			} else {
				y.set(0, interpreter)
			}
		}

		y.static = true
	} else {
		y.unset(interpreter)
		y.static = false
	}

	y.AddHistory(vm.SLT, *pc, y.static)

	return nil, nil
}

func opSgt(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if x.static && x.IsValue() && y.static && y.IsValue() {
		xSign := x.v.Cmp(tt255)
		ySign := y.v.Cmp(tt255)

		switch {
		case xSign >= 0 && ySign < 0:
			y.set(0, interpreter)

		case xSign < 0 && ySign >= 0:
			y.set(1, interpreter)

		default:
			if x.v.Cmp(y.v) > 0 {
				y.set(1, interpreter)
			} else {
				y.set(0, interpreter)
			}
		}

		y.static = true
	} else {
		y.unset(interpreter)
		y.static = false
	}

	y.AddHistory(vm.SGT, *pc, y.static)

	return nil, nil
}

func opEq(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if x.static && x.IsValue() && y.static && y.IsValue() {
		if x.v.Cmp(y.v) == 0 {
			y.set(1, interpreter)
		} else {
			y.set(0, interpreter)
		}

		y.static = true
	} else {
		y.unset(interpreter)
		y.static = false
	}

	y.AddHistory(vm.EQ, *pc, y.static)

	return nil, nil
}

func opIszero(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if x.static && x.IsValue() {
		if x.v.Sign() > 0 {
			x.set(0, interpreter)
		} else {
			x.set(1, interpreter)
		}

		x.static = true
	} else {
		x.unset(interpreter)
		x.static = false
	}

	x.AddHistory(vm.ISZERO, *pc, x.static)

	return nil, nil
}

func opAnd(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	x.static = NotStaticIfOneNotStatic(x, y)
	if x.static && x.IsValue() && y.static && y.IsValue() {
		x.v = x.v.And(x.v, y.v)
	} else {
		x.unset(interpreter)
	}

	x.AddHistory(vm.AND, *pc, x.static)
	stack.push(x)

	return nil, nil
}

func opOr(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	y.static = NotStaticIfOneNotStatic(x, y)

	if y.static && y.IsValue() && x.static && x.IsValue() {
		y.v.Or(x.v, y.v)
	} else {
		y.unset(interpreter)
	}

	y.AddHistory(vm.OR, *pc, y.static)

	return nil, nil
}

func opXor(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.peek()
	if err != nil {
		return nil, err
	}

	y.static = NotStaticIfOneNotStatic(x, y)

	if y.static && y.IsValue() && x.static && x.IsValue() {
		y.v.Xor(x.v, y.v)
	} else {
		y.unset(interpreter)
	}

	y.AddHistory(vm.XOR, *pc, y.static)

	return nil, nil
}

func opByte(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	th, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(th.v)

	val, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if th.static && th.IsValue() {
		if th.v.Cmp(common.Big32) < 0 {
			if val.static && val.IsValue() {
				b := math.Byte(val.v, 32, int(th.v.Int64()))
				val.set(0, interpreter).SetUint64(uint64(b))
				val.static = true
			} else {
				val.unset(interpreter)
				val.static = false
			}
		} else {
			val.set(0, interpreter)
			val.static = true
		}
	} else {
		val.unset(interpreter)
		val.static = false
	}

	val.AddHistory(vm.BYTE, *pc, val.static)

	return nil, nil
}

func opAddmod(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	z, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(z.v)

	if z.static && z.IsValue() {
		if z.v.Cmp(bigZero) > 0 {
			if x.static && x.IsValue() && y.static && y.IsValue() {
				x.v.Add(x.v, y.v)
				x.v.Mod(x.v, z.v)
				x.static = true
			} else {
				x.unset(interpreter)
				x.static = false
			}
		} else {
			x.set(0, interpreter)
			x.static = true
		}
	} else {
		x.unset(interpreter)
		x.static = false
	}

	x.AddHistory(vm.ADDMOD, *pc, x.static)
	stack.push(x)

	return nil, nil
}

func opMulmod(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	z, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(z.v)

	if z.static && z.IsValue() {
		if z.v.Cmp(bigZero) > 0 {
			if x.static && x.IsValue() && y.static && y.IsValue() {
				x.v.Mul(x.v, y.v)
				x.v.Mod(x.v, z.v)
				x.static = true
			} else {
				x.unset(interpreter)
				x.static = false
			}
		} else {
			x.set(0, interpreter)
			x.static = true
		}
	} else {
		x.unset(interpreter)
		x.static = false
	}

	x.AddHistory(vm.MULMOD, *pc, x.static)
	stack.push(x)

	return nil, nil
}

// opSHL implements Shift Left
// The SHL instruction (shift left) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the left by arg1 number of bits.
func opSHL(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(shift.v)

	value, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if shift.static && shift.IsValue() {
		math.U256(shift.v)
	}
	if value.static && value.IsValue() {
		math.U256(value.v)
	}

	if shift.static && shift.IsValue() {
		if shift.v.Cmp(common.Big256) >= 0 {
			value.set(0, interpreter)
			value.static = true
		} else {
			if value.static && value.IsValue() {
				n := uint(shift.v.Uint64())
				math.U256(value.v.Lsh(value.v, n))
				value.static = true
			} else {
				value.unset(interpreter)
				value.static = false
			}
		}
	} else {
		value.unset(interpreter)
		value.static = false
	}

	value.AddHistory(vm.SHL, *pc, value.static)

	return nil, nil
}

// opSHR implements Logical Shift Right
// The SHR instruction (logical shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with zero fill.
func opSHR(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(shift.v)

	value, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if shift.static && shift.IsValue() {
		math.U256(shift.v)
	}
	if value.static && value.IsValue() {
		math.U256(value.v)
	}

	if shift.static && shift.IsValue() {
		if shift.v.Cmp(common.Big256) >= 0 {
			value.set(0, interpreter)
			value.static = true
		} else {
			if value.static && value.IsValue() {
				n := uint(shift.v.Uint64())
				math.U256(value.v.Rsh(value.v, n))
				value.static = true
			} else {
				value.unset(interpreter)
				value.static = false
			}
		}
	} else {
		value.unset(interpreter)
		value.static = false
	}

	value.AddHistory(vm.SHR, *pc, value.static)

	return nil, nil
}

// opSAR implements Arithmetic Shift Right
// The SAR instruction (arithmetic shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with sign extension.
func opSAR(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	shift, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(shift.v)

	value, err := stack.pop()
	if err != nil {
		return nil, err
	}

	if shift.static && shift.IsValue() {
		math.U256(shift.v)
	}
	if value.static && value.IsValue() {
		math.U256(value.v)
	}

	if shift.static && shift.IsValue() && value.static && value.IsValue() {
		if shift.v.Cmp(common.Big256) >= 0 {
			if value.v.Sign() >= 0 {
				value.set(0, interpreter)
			} else {
				value.v.SetInt64(-1)
			}
		} else {
			n := uint(shift.v.Uint64())
			value.v.Rsh(value.v, n)
		}

		math.U256(value.v)
		value.static = true
	} else {
		value.unset(interpreter)
		value.static = false
	}

	value.AddHistory(vm.SAR, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opSha3(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	y, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(y.v)

	value := NewNonStaticCell() // fixme it's stricter than it could be
	value.AddHistory(vm.SHA3, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opAddress(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewStaticCell()
	value.set(0, interpreter).SetBytes(contract.CodeAddr.Bytes())
	value.AddHistory(vm.ADDRESS, *pc, value.static)

	stack.push(value)
	return nil, nil
}

func opBalance(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	slot, err := stack.peek()
	if err != nil {
		return nil, err
	}
	slot.unset(interpreter)
	slot.static = false
	slot.AddHistory(vm.BALANCE, *pc, slot.static)

	return nil, nil
}

func opOrigin(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.ORIGIN, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opCaller(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.CALLER, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opCallValue(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.CALLVALUE, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opCallDataLoad(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)

	value := NewNonStaticCell()
	value.AddHistory(vm.CALLDATALOAD, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opCallDataSize(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.CALLDATASIZE, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opCallDataCopy(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(3, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func opReturnDataSize(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.RETURNDATASIZE, *pc, value.static)
	stack.push(value) // fixme: stricter than it could be

	return nil, nil
}

func opReturnDataCopy(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(3, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func opExtCodeSize(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	slot, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if slot.IsStatic() {
		slot.static = true
		// fixme we can set value here if we get extContract by address
	} else {
		slot.unset(interpreter)
		slot.static = false
	}
	slot.AddHistory(vm.EXTCODECOPY, *pc, slot.static)

	return nil, nil
}

func opCodeSize(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewStaticCell()
	// fixme we can set value here if we get the contract untrimmed code size
	value.AddHistory(vm.CODESIZE, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opCodeCopy(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(3, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func opExtCodeCopy(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(4, interpreter)
	if err != nil {
		return nil, err
	}

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
func opExtCodeHash(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	slot, err := stack.peek()
	if err != nil {
		return nil, err
	}

	if slot.static && slot.IsValue() {
		slot.unset(interpreter)
		slot.static = true
		/*
			address := common.BigToAddress(slot.v)
			if interpreter.evm.IntraBlockState.Empty(address) {
				slot = getInt(0, interpreter)
			} else {
				slot.SetBytes(interpreter.evm.IntraBlockState.GetCodeHash(address).Bytes())
			}
		*/
	} else {
		slot.unset(interpreter)
		slot.static = false
	}

	slot.AddHistory(vm.EXTCODEHASH, *pc, slot.static)

	return nil, nil
}

func opGasprice(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.GASPRICE, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opBlockhash(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	slot, err := stack.pop()
	if err != nil {
		return nil, err
	}

	if slot.IsStatic() && slot.IsValue() {
		slot.unset(interpreter)
		slot.static = true
	} else {
		slot.static = false
	}
	slot.AddHistory(vm.BLOCKHASH, *pc, slot.static)
	stack.push(slot)

	return nil, nil
}

func opCoinbase(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.COINBASE, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opTimestamp(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.TIMESTAMP, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opNumber(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.NUMBER, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opDifficulty(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.DIFFICULTY, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opGasLimit(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.GASLIMIT, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opPop(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	x, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(x.v)
	return nil, nil
}

func opMload(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	v, err := stack.peek()
	if err != nil {
		return nil, err
	}

	v.static = false // fixme: not true if we introduce a memory fake
	v.AddHistory(vm.MLOAD, *pc, v.static)

	return nil, nil
}

func opMstore(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(2, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func opMstore8(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(2, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func opSload(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	loc, err := stack.peek()
	if err != nil {
		return nil, err
	}

	loc = NewNonStaticCell() // fixme: not true if we introduce memory type
	loc.AddHistory(vm.SLOAD, *pc, loc.static)

	return nil, nil
}

func opSstore(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(2, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func opJump(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(pos.v)

	/*
		if pos.static {
			fmt.Printf("jumpiT: on %x to %v\nValue history %v\n", *pc, pos.v, pos.History())
		}
	*/

	if !pos.static {
		return nil, fmt.Errorf("opJumpi: %w on %v\nValue history %v\n", ErrNonStatic, spew.Sdump(pc), pos.History())
	}
	if pos.static && !pos.IsValue() {
		return nil, fmt.Errorf("jumpi: %w on %v\nValue history %v\n", ErrNoValueStatic, spew.Sdump(pc), pos.History())
	}
	if !contract.validJumpdest(pos.v) {
		return nil, fmt.Errorf("%w on %v\nValue history %v\n", ErrInvalidJump, spew.Sdump(pc), pos.History())
	}

	*pc = pos.v.Uint64()

	return nil, nil
}

func opJumpi(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(pos.v)

	cond, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(cond.v)

	/*
		if pos.static {
			fmt.Printf("jumpiT: on %x to %v\nValue history %v\n", *pc, pos.v, pos.History())
		}
	*/

	if cond.v.Sign() != 0 {
		if !pos.static {
			return nil, fmt.Errorf("opJumpi: %w on %v\nValue history %v\n", ErrNonStatic, spew.Sdump(pc), pos.History())
		}
		if pos.static && !pos.IsValue() {
			return nil, fmt.Errorf("jumpi: %w on %v\nValue history %v\n", ErrNoValueStatic, spew.Sdump(pc), pos.History())
		}

		if !contract.validJumpdest(pos.v) {
			return nil, fmt.Errorf("jumpi: %w on %v\nValue history %v\n", ErrInvalidJump, spew.Sdump(pc), pos.History())
		}
		*pc = pos.v.Uint64()
	} else {
		*pc++
	}

	return nil, nil
}

func opJumpiJUMP(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(pos.v)

	cond, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(cond.v)

	/*
		if pos.static {
			fmt.Printf("jumpiT: on %x to %v\nValue history %v\n", *pc, pos.v, pos.History())
		}
	*/

	if !pos.static {
		return nil, fmt.Errorf("opJumpiJUMP: %w on %v\nValue history %v\n", ErrNonStatic, spew.Sdump(pc), pos.History())
	}
	if pos.static && !pos.IsValue() {
		return nil, fmt.Errorf("jumpi: %w on %v. jump to %v\nValue history %v\n", ErrNoValueStatic, spew.Sdump(pc), pos.v, pos.History())
	}

	if !contract.validJumpdest(pos.v) {
		return nil, fmt.Errorf("jumpi: %w on %v\nValue history %v\n", ErrInvalidJump, spew.Sdump(pc), pos.History())
	}

	*pc = pos.v.Uint64()

	return nil, nil
}

func opJumpiNotJUMP(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	pos, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(pos.v)

	cond, err := stack.pop()
	if err != nil {
		return nil, err
	}
	defer interpreter.IntPool.Put(cond.v)

	*pc = *pc + 1

	return nil, nil
}

func opJumpdest(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	return nil, nil
}

func opPc(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewStaticCell()
	value.set(0, interpreter).SetUint64(*pc)

	value.AddHistory(vm.PC, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opMsize(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.MSIZE, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opGas(pc *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.GAS, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func opCreate(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(3, interpreter)
	if err != nil {
		return nil, err
	}

	value := NewNonStaticCell()
	value.AddHistory(vm.CREATE, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opCreate2(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(4, interpreter)
	if err != nil {
		return nil, err
	}

	value := NewNonStaticCell()
	value.AddHistory(vm.CREATE2, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opCall(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(7, interpreter)
	if err != nil {
		return nil, err
	}

	value := NewNonStaticCell()
	value.AddHistory(vm.CALL, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opCallCode(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	err := stack.remove(7, interpreter)
	if err != nil {
		return nil, err
	}

	value := NewNonStaticCell()
	value.AddHistory(vm.CALLCODE, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opDelegateCall(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	err := stack.remove(6, interpreter)
	if err != nil {
		return nil, err
	}

	value := NewNonStaticCell()
	value.AddHistory(vm.DELEGATECALL, *pc, value.static)
	stack.push(value)

	return nil, nil
}

func opStaticCall(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	err := stack.remove(6, interpreter)
	if err != nil {
		return nil, err
	}

	value := NewNonStaticCell()
	value.AddHistory(vm.STATICCALL, *pc, value.static) // fixme stricter than it could be
	stack.push(value)

	return nil, nil
}

func opReturn(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(2, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, ErrReturn
}

func opRevert(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(2, interpreter)
	if err != nil {
		return nil, err
	}

	return nil, ErrRevert
}

func opStop(_ *uint64, _ *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	return nil, ErrStop
}

func opSuicide(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	err := stack.remove(1, interpreter)
	if err != nil {
		return nil, err
	}
	return nil, ErrSelfDestruct
}

// following functions are used by the instruction jump  table
type executionFunc func(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, memory *vm.Memory, stack *Stack) ([]byte, error)

// make log instruction function
func makeLog(size int) executionFunc {
	return func(_ *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		err := stack.remove(2, interpreter)
		if err != nil {
			return nil, err
		}

		for i := 0; i < size; i++ {
			_, err = stack.pop()
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}
}

// opPush1 is a specialized version of pushN
func opPush1(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	var (
		codeLen = uint64(len(contract.Code))
		integer = getInt(0, interpreter)
	)
	*pc += 1

	c := NewStaticCell()

	if *pc < codeLen {
		integer.SetUint64(uint64(contract.Code[*pc]))
	}

	c.SetValue(integer)
	stack.push(c)

	c.AddHistory(vm.PUSH1, *pc, c.static)

	return nil, nil
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, interpreter *vm.EVMInterpreter, contract *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		codeLen := len(contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := getInt(0, interpreter)
		integer.SetBytes(common.RightPadBytes(contract.Code[startMin:endMin], pushByteSize))

		c := NewStaticCell()
		c.SetValue(integer)
		c.AddHistory(vm.PUSH, *pc, c.static)
		stack.push(c)

		*pc += size
		return nil, nil
	}
}

// make dup instruction function
func makeDup(size int64) executionFunc {
	return func(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		err := stack.dup(int(size), vm.DUP, *pc)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size++
	return func(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
		err := stack.swap(int(size), vm.SWAP, *pc)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func opSelfBalance(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewNonStaticCell()
	value.AddHistory(vm.SELFBALANCE, *pc, value.static)
	stack.push(value)
	return nil, nil
}

// opChainID implements CHAINID opcode
func opChainID(pc *uint64, interpreter *vm.EVMInterpreter, _ *Contract, _ *vm.Memory, stack *Stack) ([]byte, error) {
	value := NewStaticCell()
	value.AddHistory(vm.CHAINID, *pc, value.static)
	stack.push(value)
	return nil, nil
}

func getInt(n int64, interpreter *vm.EVMInterpreter) *big.Int {
	if interpreter == nil {
		panic(1)
	}
	if interpreter.IntPool == nil {
		panic(2)
	}
	if v := interpreter.IntPool.Get(); v == nil {
		panic(3)
	} else {
		return interpreter.IntPool.Get().SetInt64(n)
	}
}
