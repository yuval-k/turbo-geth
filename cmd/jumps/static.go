package static

import "github.com/ledgerwatch/turbo-geth/core/vm"

type Command struct {
	In           int
	Out          int
	IsStaticFunc func(...*cell) bool
	IsStatic     bool
	execute      executionFunc
}

func Static(_ ...*cell) bool {
	return true
}

func NotStatic(_ ...*cell) bool {
	return false
}

func NewCommand(in int, out int, fn executionFunc, isStaticFunc func(...*cell) bool) Command {
	return Command{in, out, isStaticFunc, false, fn}
}

var Commands = map[vm.OpCode]Command{
	vm.STOP:       NewCommand(0, 0, opStop, Static),
	vm.ADD:        NewCommand(2, 1, opAdd, NotStaticIfOneNotStatic),
	vm.MUL:        NewCommand(2, 1, opMul, NotStaticIfOneNotStatic),
	vm.SUB:        NewCommand(2, 1, opSub, NotStaticIfOneNotStatic),
	vm.DIV:        NewCommand(2, 1, opDiv, NotStaticIfOneNotStatic),
	vm.SDIV:       NewCommand(2, 1, opSdiv, NotStaticIfOneNotStatic),
	vm.MOD:        NewCommand(2, 1, opMod, NotStaticIfOneNotStatic),
	vm.SMOD:       NewCommand(2, 1, opSmod, NotStaticIfOneNotStatic),
	vm.ADDMOD:     NewCommand(3, 1, opAddmod, NotStaticIfOneNotStatic),
	vm.MULMOD:     NewCommand(3, 1, opMulmod, NotStaticIfOneNotStatic),
	vm.EXP:        NewCommand(2, 1, opExp, NotStaticIfOneNotStatic),
	vm.SIGNEXTEND: NewCommand(2, 1, opSignExtend, NotStaticIfOneNotStatic),

	vm.LT:     NewCommand(2, 1, opLt, NotStaticIfOneNotStatic),
	vm.GT:     NewCommand(2, 1, opGt, NotStaticIfOneNotStatic),
	vm.SLT:    NewCommand(2, 1, opSlt, NotStaticIfOneNotStatic),
	vm.SGT:    NewCommand(2, 1, opSgt, NotStaticIfOneNotStatic),
	vm.EQ:     NewCommand(2, 1, opEq, NotStaticIfOneNotStatic),
	vm.ISZERO: NewCommand(1, 1, opIszero, NotStaticIfOneNotStatic),
	vm.AND:    NewCommand(2, 1, opAnd, NotStaticIfOneNotStatic),
	vm.OR:     NewCommand(2, 1, opOr, NotStaticIfOneNotStatic),
	vm.XOR:    NewCommand(2, 1, opXor, NotStaticIfOneNotStatic),
	vm.NOT:    NewCommand(1, 1, opNot, NotStaticIfOneNotStatic),
	vm.BYTE:   NewCommand(2, 1, opByte, NotStaticIfOneNotStatic),
	vm.SHL:    NewCommand(2, 1, opSHL, NotStaticIfOneNotStatic),
	vm.SHR:    NewCommand(2, 1, opSHR, NotStaticIfOneNotStatic),
	vm.SAR:    NewCommand(2, 1, opSAR, NotStaticIfOneNotStatic),

	vm.SHA3: NewCommand(2, 1, opSha3, NotStaticIfOneNotStatic),

	vm.ADDRESS:        NewCommand(0, 1, opAddress, Static), // +++
	vm.BALANCE:        NewCommand(1, 1, opBalance, NotStatic),
	vm.ORIGIN:         NewCommand(0, 1, opOrigin, NotStatic),
	vm.CALLER:         NewCommand(0, 1, opCaller, NotStatic),
	vm.CALLVALUE:      NewCommand(0, 1, opCallValue, NotStatic),
	vm.CALLDATALOAD:   NewCommand(1, 1, opCallDataLoad, NotStatic),
	vm.CALLDATASIZE:   NewCommand(0, 1, opCallDataSize, NotStatic),
	vm.CALLDATACOPY:   NewCommand(3, 1, opCallDataCopy, NotStatic),
	vm.CODESIZE:       NewCommand(0, 1, opCodeSize, NotStatic),
	vm.CODECOPY:       NewCommand(3, 0, opCodeCopy, Static), // +++
	vm.GASPRICE:       NewCommand(0, 1, opGasprice, NotStatic),
	vm.EXTCODESIZE:    NewCommand(1, 1, opExtCodeSize, NotStatic),               // +++
	vm.EXTCODECOPY:    NewCommand(4, 0, opExtCodeCopy, Static),                  // +++
	vm.RETURNDATASIZE: NewCommand(0, 1, opReturnDataSize, NotStatic),            // fixme // +++
	vm.RETURNDATACOPY: NewCommand(3, 0, opReturnDataCopy, Static),               // fixme // +++
	vm.EXTCODEHASH:    NewCommand(1, 1, opExtCodeHash, NotStaticIfOneNotStatic), // +++

	vm.BLOCKHASH:   NewCommand(1, 1, opBlockhash, NotStatic),
	vm.COINBASE:    NewCommand(0, 1, opCoinbase, NotStatic),
	vm.TIMESTAMP:   NewCommand(0, 1, opTimestamp, NotStatic),
	vm.NUMBER:      NewCommand(0, 1, opNumber, NotStatic),
	vm.DIFFICULTY:  NewCommand(0, 1, opDifficulty, NotStatic),
	vm.GASLIMIT:    NewCommand(0, 1, opGasLimit, NotStatic),
	vm.CHAINID:     NewCommand(0, 1, opChainID, Static), // +++
	vm.SELFBALANCE: NewCommand(0, 1, opSelfBalance, NotStatic),

	vm.POP:      NewCommand(1, 0, opPop, Static),
	vm.MLOAD:    NewCommand(1, 1, opMload, NotStatic),
	vm.MSTORE:   NewCommand(2, 0, opMstore, Static),
	vm.MSTORE8:  NewCommand(2, 0, opMstore8, Static),
	vm.SLOAD:    NewCommand(1, 1, opSload, NotStatic),
	vm.SSTORE:   NewCommand(2, 1, opSstore, Static),
	vm.JUMP:     NewCommand(1, 0, opJump, Static),
	vm.JUMPI:    NewCommand(2, 0, opJumpi, Static),
	vm.PC:       NewCommand(0, 1, opPc, Static), // +++
	vm.MSIZE:    NewCommand(0, 1, opMsize, Static),
	vm.GAS:      NewCommand(0, 1, opGas, NotStatic),
	vm.JUMPDEST: NewCommand(0, 0, opJumpdest, Static),

	vm.PUSH1:  NewCommand(0, 1, opPush1, Static),          // +++
	vm.PUSH2:  NewCommand(0, 1, makePush(2, 2), Static),   // +++
	vm.PUSH3:  NewCommand(0, 1, makePush(3, 3), Static),   // +++
	vm.PUSH4:  NewCommand(0, 1, makePush(4, 4), Static),   // +++
	vm.PUSH5:  NewCommand(0, 1, makePush(5, 5), Static),   // +++
	vm.PUSH6:  NewCommand(0, 1, makePush(6, 6), Static),   // +++
	vm.PUSH7:  NewCommand(0, 1, makePush(7, 7), Static),   // +++
	vm.PUSH8:  NewCommand(0, 1, makePush(8, 8), Static),   // +++
	vm.PUSH9:  NewCommand(0, 1, makePush(9, 9), Static),   // +++
	vm.PUSH10: NewCommand(0, 1, makePush(10, 10), Static), // +++
	vm.PUSH11: NewCommand(0, 1, makePush(11, 11), Static), // +++
	vm.PUSH12: NewCommand(0, 1, makePush(12, 12), Static), // +++
	vm.PUSH13: NewCommand(0, 1, makePush(13, 13), Static), // +++
	vm.PUSH14: NewCommand(0, 1, makePush(14, 14), Static), // +++
	vm.PUSH15: NewCommand(0, 1, makePush(15, 15), Static), // +++
	vm.PUSH16: NewCommand(0, 1, makePush(16, 16), Static), // +++
	vm.PUSH17: NewCommand(0, 1, makePush(17, 17), Static), // +++
	vm.PUSH18: NewCommand(0, 1, makePush(18, 18), Static), // +++
	vm.PUSH19: NewCommand(0, 1, makePush(19, 19), Static), // +++
	vm.PUSH20: NewCommand(0, 1, makePush(20, 20), Static), // +++
	vm.PUSH21: NewCommand(0, 1, makePush(21, 21), Static), // +++
	vm.PUSH22: NewCommand(0, 1, makePush(22, 22), Static), // +++
	vm.PUSH23: NewCommand(0, 1, makePush(23, 23), Static), // +++
	vm.PUSH24: NewCommand(0, 1, makePush(24, 24), Static), // +++
	vm.PUSH25: NewCommand(0, 1, makePush(25, 25), Static), // +++
	vm.PUSH26: NewCommand(0, 1, makePush(26, 26), Static), // +++
	vm.PUSH27: NewCommand(0, 1, makePush(27, 27), Static), // +++
	vm.PUSH28: NewCommand(0, 1, makePush(28, 28), Static), // +++
	vm.PUSH29: NewCommand(0, 1, makePush(29, 29), Static), // +++
	vm.PUSH30: NewCommand(0, 1, makePush(30, 30), Static), // +++
	vm.PUSH31: NewCommand(0, 1, makePush(31, 31), Static), // +++
	vm.PUSH32: NewCommand(0, 1, makePush(32, 32), Static), // +++

	// fixme can be analysed better
	vm.DUP1:  NewCommand(1, 2, makeDup(1), NotStaticIfOneNotStatic),
	vm.DUP2:  NewCommand(2, 3, makeDup(2), NotStaticIfOneNotStatic),
	vm.DUP3:  NewCommand(3, 4, makeDup(3), NotStaticIfOneNotStatic),
	vm.DUP4:  NewCommand(4, 5, makeDup(4), NotStaticIfOneNotStatic),
	vm.DUP5:  NewCommand(5, 6, makeDup(5), NotStaticIfOneNotStatic),
	vm.DUP6:  NewCommand(6, 7, makeDup(6), NotStaticIfOneNotStatic),
	vm.DUP7:  NewCommand(7, 8, makeDup(7), NotStaticIfOneNotStatic),
	vm.DUP8:  NewCommand(8, 9, makeDup(8), NotStaticIfOneNotStatic),
	vm.DUP9:  NewCommand(9, 10, makeDup(9), NotStaticIfOneNotStatic),
	vm.DUP10: NewCommand(10, 11, makeDup(10), NotStaticIfOneNotStatic),
	vm.DUP11: NewCommand(11, 12, makeDup(11), NotStaticIfOneNotStatic),
	vm.DUP12: NewCommand(12, 13, makeDup(12), NotStaticIfOneNotStatic),
	vm.DUP13: NewCommand(13, 14, makeDup(13), NotStaticIfOneNotStatic),
	vm.DUP14: NewCommand(14, 15, makeDup(14), NotStaticIfOneNotStatic),
	vm.DUP15: NewCommand(15, 16, makeDup(15), NotStaticIfOneNotStatic),
	vm.DUP16: NewCommand(16, 17, makeDup(16), NotStaticIfOneNotStatic),

	// special cases!
	vm.SWAP1:  NewCommand(1, 1, makeSwap(1), NotStaticIfOneNotStatic),
	vm.SWAP2:  NewCommand(2, 2, makeSwap(2), NotStaticIfOneNotStatic),
	vm.SWAP3:  NewCommand(3, 3, makeSwap(3), NotStaticIfOneNotStatic),
	vm.SWAP4:  NewCommand(4, 4, makeSwap(4), NotStaticIfOneNotStatic),
	vm.SWAP5:  NewCommand(5, 5, makeSwap(5), NotStaticIfOneNotStatic),
	vm.SWAP6:  NewCommand(6, 6, makeSwap(6), NotStaticIfOneNotStatic),
	vm.SWAP7:  NewCommand(7, 7, makeSwap(7), NotStaticIfOneNotStatic),
	vm.SWAP8:  NewCommand(8, 8, makeSwap(8), NotStaticIfOneNotStatic),
	vm.SWAP9:  NewCommand(9, 9, makeSwap(9), NotStaticIfOneNotStatic),
	vm.SWAP10: NewCommand(10, 10, makeSwap(10), NotStaticIfOneNotStatic),
	vm.SWAP11: NewCommand(11, 11, makeSwap(11), NotStaticIfOneNotStatic),
	vm.SWAP12: NewCommand(12, 12, makeSwap(12), NotStaticIfOneNotStatic),
	vm.SWAP13: NewCommand(13, 13, makeSwap(13), NotStaticIfOneNotStatic),
	vm.SWAP14: NewCommand(14, 14, makeSwap(14), NotStaticIfOneNotStatic),
	vm.SWAP15: NewCommand(15, 15, makeSwap(15), NotStaticIfOneNotStatic),
	vm.SWAP16: NewCommand(16, 16, makeSwap(16), NotStaticIfOneNotStatic),

	vm.LOG0: NewCommand(2, 0, makeLog(0), Static),
	vm.LOG1: NewCommand(3, 0, makeLog(1), Static),
	vm.LOG2: NewCommand(4, 0, makeLog(2), Static),
	vm.LOG3: NewCommand(5, 0, makeLog(3), Static),
	vm.LOG4: NewCommand(6, 0, makeLog(4), Static),

	vm.CREATE:       NewCommand(3, 1, opCreate, NotStatic),   // fixme
	vm.CALL:         NewCommand(7, 1, opCall, NotStatic),     // fixme
	vm.CALLCODE:     NewCommand(7, 1, opCallCode, NotStatic), // fixme
	vm.RETURN:       NewCommand(2, 0, opReturn, Static),
	vm.DELEGATECALL: NewCommand(6, 1, opDelegateCall, NotStatic), // fixme
	vm.CREATE2:      NewCommand(4, 1, opCreate2, NotStatic),
	vm.STATICCALL:   NewCommand(6, 1, opStaticCall, NotStatic), // fixme

	vm.REVERT:       NewCommand(2, 0, opRevert, Static),
	vm.SELFDESTRUCT: NewCommand(1, 0, opSuicide, Static),
}

func IsStop(op vm.OpCode) bool {
	switch op {
	case vm.RETURN, vm.REVERT, vm.SELFDESTRUCT, vm.STOP:
		return true
	}

	return false
}
