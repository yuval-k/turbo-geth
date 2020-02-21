package vm

type Command struct {
	in           int
	out          int
	isStaticFunc func(...Command) bool
	isStatic     bool
}

func Static(_ ...Command) bool {
	return true
}

func NotStatic(_ ...Command) bool {
	return false
}

func NotStaticIfOneNotStatic(cmds ...Command) bool {
	for _, cmd := range cmds {
		if !cmd.isStatic {
			return false
		}
	}
	return true
}

func NewCommand(in int, out int, isStaticFunc func(...Command) bool) Command {
	return Command{in, out, isStaticFunc, false}
}

var Commands = map[OpCode]Command{
	STOP:       NewCommand(0, 0, Static),
	ADD:        NewCommand(2, 1, NotStaticIfOneNotStatic),
	MUL:        NewCommand(2, 1, NotStaticIfOneNotStatic),
	SUB:        NewCommand(2, 1, NotStaticIfOneNotStatic),
	DIV:        NewCommand(2, 1, NotStaticIfOneNotStatic),
	SDIV:       NewCommand(2, 1, NotStaticIfOneNotStatic),
	MOD:        NewCommand(2, 1, NotStaticIfOneNotStatic),
	SMOD:       NewCommand(2, 1, NotStaticIfOneNotStatic),
	ADDMOD:     NewCommand(3, 1, NotStaticIfOneNotStatic),
	MULMOD:     NewCommand(3, 1, NotStaticIfOneNotStatic),
	EXP:        NewCommand(2, 1, NotStaticIfOneNotStatic),
	SIGNEXTEND: NewCommand(2, 1, NotStaticIfOneNotStatic),

	LT:     NewCommand(2, 1, NotStaticIfOneNotStatic),
	GT:     NewCommand(2, 1, NotStaticIfOneNotStatic),
	SLT:    NewCommand(2, 1, NotStaticIfOneNotStatic),
	SGT:    NewCommand(2, 1, NotStaticIfOneNotStatic),
	EQ:     NewCommand(2, 1, NotStaticIfOneNotStatic),
	ISZERO: NewCommand(1, 1, NotStaticIfOneNotStatic),
	AND:    NewCommand(2, 1, NotStaticIfOneNotStatic),
	OR:     NewCommand(2, 1, NotStaticIfOneNotStatic),
	XOR:    NewCommand(2, 1, NotStaticIfOneNotStatic),
	NOT:    NewCommand(1, 1, NotStaticIfOneNotStatic),
	BYTE:   NewCommand(2, 1, NotStaticIfOneNotStatic),
	SHL:    NewCommand(2, 1, NotStaticIfOneNotStatic),
	SHR:    NewCommand(2, 1, NotStaticIfOneNotStatic),
	SAR:    NewCommand(2, 1, NotStaticIfOneNotStatic),

	SHA3: NewCommand(2, 1, NotStaticIfOneNotStatic),

	ADDRESS:        NewCommand(0, 1, Static),
	BALANCE:        NewCommand(1, 1, NotStaticIfOneNotStatic),
	ORIGIN:         NewCommand(0, 1, NotStatic),
	CALLER:         NewCommand(0, 1, NotStatic),
	CALLVALUE:      NewCommand(0, 1, NotStatic),
	CALLDATALOAD:   NewCommand(1, 1, NotStatic),
	CALLDATASIZE:   NewCommand(0, 1, NotStatic),
	CALLDATACOPY:   NewCommand(3, 1, NotStatic),
	CODESIZE:       NewCommand(0, 1, Static),
	CODECOPY:       NewCommand(0, 3, Static),
	GASPRICE:       NewCommand(0, 1, NotStatic),
	EXTCODESIZE:    NewCommand(1, 1, Static),
	EXTCODECOPY:    NewCommand(4, 0, Static),
	RETURNDATASIZE: NewCommand(0, 1, NotStatic), // fixme
	RETURNDATACOPY: NewCommand(3, 0, Static),    // fixme
	EXTCODEHASH:    NewCommand(1, 1, NotStaticIfOneNotStatic),

	BLOCKHASH:   {},
	COINBASE:    {},
	TIMESTAMP:   {},
	NUMBER:      {},
	DIFFICULTY:  {},
	GASLIMIT:    {},
	CHAINID:     {},
	SELFBALANCE: {},

	POP:      {},
	MLOAD:    {},
	MSTORE:   {},
	MSTORE8:  {},
	SLOAD:    {},
	SSTORE:   {},
	JUMP:     {},
	JUMPI:    {},
	PC:       {},
	MSIZE:    {},
	GAS:      {},
	JUMPDEST: {},

	PUSH1:  NewCommand(0, 1, Static),
	PUSH2:  NewCommand(0, 2, Static),
	PUSH3:  NewCommand(0, 3, Static),
	PUSH4:  NewCommand(0, 4, Static),
	PUSH5:  NewCommand(0, 5, Static),
	PUSH6:  NewCommand(0, 6, Static),
	PUSH7:  NewCommand(0, 7, Static),
	PUSH8:  NewCommand(0, 8, Static),
	PUSH9:  NewCommand(0, 9, Static),
	PUSH10: NewCommand(0, 10, Static),
	PUSH11: NewCommand(0, 11, Static),
	PUSH12: NewCommand(0, 12, Static),
	PUSH13: NewCommand(0, 13, Static),
	PUSH14: NewCommand(0, 14, Static),
	PUSH15: NewCommand(0, 15, Static),
	PUSH16: NewCommand(0, 16, Static),
	PUSH17: NewCommand(0, 17, Static),
	PUSH18: NewCommand(0, 18, Static),
	PUSH19: NewCommand(0, 19, Static),
	PUSH20: NewCommand(0, 20, Static),
	PUSH21: NewCommand(0, 21, Static),
	PUSH22: NewCommand(0, 22, Static),
	PUSH23: NewCommand(0, 23, Static),
	PUSH24: NewCommand(0, 24, Static),
	PUSH25: NewCommand(0, 25, Static),
	PUSH26: NewCommand(0, 26, Static),
	PUSH27: NewCommand(0, 27, Static),
	PUSH28: NewCommand(0, 28, Static),
	PUSH29: NewCommand(0, 29, Static),
	PUSH30: NewCommand(0, 30, Static),
	PUSH31: NewCommand(0, 31, Static),
	PUSH32: NewCommand(0, 32, Static),

	DUP1:   NewCommand(1, 2, NotStaticIfOneNotStatic),
	DUP2:   NewCommand(2, 3, NotStaticIfOneNotStatic),
	DUP3:   NewCommand(3, 4, NotStaticIfOneNotStatic),
	DUP4:   NewCommand(4, 5, NotStaticIfOneNotStatic),
	DUP5:   NewCommand(5, 6, NotStaticIfOneNotStatic),
	DUP6:   NewCommand(6, 7, NotStaticIfOneNotStatic),
	DUP7:   NewCommand(7, 8, NotStaticIfOneNotStatic),
	DUP8:   NewCommand(8, 9, NotStaticIfOneNotStatic),
	DUP9:   NewCommand(9, 10, NotStaticIfOneNotStatic),
	DUP10:  NewCommand(10, 11, NotStaticIfOneNotStatic),
	DUP11:  NewCommand(11, 12, NotStaticIfOneNotStatic),
	DUP12:  NewCommand(12, 13, NotStaticIfOneNotStatic),
	DUP13:  NewCommand(13, 14, NotStaticIfOneNotStatic),
	DUP14:  NewCommand(14, 15, NotStaticIfOneNotStatic),
	DUP15:  NewCommand(15, 16, NotStaticIfOneNotStatic),
	DUP16:  NewCommand(16, 17, NotStaticIfOneNotStatic),

	// special cases!
	SWAP1:  NewCommand(1, 1, NotStaticIfOneNotStatic),
	SWAP2:  NewCommand(2, 2, NotStaticIfOneNotStatic),
	SWAP3:  NewCommand(3, 3, NotStaticIfOneNotStatic),
	SWAP4:  NewCommand(4, 4, NotStaticIfOneNotStatic),
	SWAP5:  NewCommand(5, 5, NotStaticIfOneNotStatic),
	SWAP6:  NewCommand(6, 6, NotStaticIfOneNotStatic),
	SWAP7:  NewCommand(7, 7, NotStaticIfOneNotStatic),
	SWAP8:  NewCommand(8, 8, NotStaticIfOneNotStatic),
	SWAP9:  NewCommand(9, 9, NotStaticIfOneNotStatic),
	SWAP10: NewCommand(10, 10, NotStaticIfOneNotStatic),
	SWAP11: NewCommand(11, 11, NotStaticIfOneNotStatic),
	SWAP12: NewCommand(12, 12, NotStaticIfOneNotStatic),
	SWAP13: NewCommand(13, 13, NotStaticIfOneNotStatic),
	SWAP14: NewCommand(14, 14, NotStaticIfOneNotStatic),
	SWAP15: NewCommand(15, 15, NotStaticIfOneNotStatic),
	SWAP16: NewCommand(16, 16, NotStaticIfOneNotStatic),

	LOG0: NewCommand(2, 0, Static),
	LOG1: NewCommand(3, 0, Static),
	LOG2: NewCommand(4, 0, Static),
	LOG3: NewCommand(5, 0, Static),
	LOG4: NewCommand(6, 0, Static),

	PUSH: NewCommand(0, 1, Static),
	DUP:  NewCommand(1, 2, NotStaticIfOneNotStatic),
	SWAP: NewCommand(1, 1, NotStaticIfOneNotStatic),

	CREATE:       {},
	CALL:         {},
	CALLCODE:     {},
	RETURN:       {},
	DELEGATECALL: {},
	CREATE2:      NewCommand(4, 1, NotStatic),
	STATICCALL:   NewCommand(6, 1, NotStatic), // fixme

	REVERT:       NewCommand(2, 0, Static),
	SELFDESTRUCT: NewCommand(1, 0, Static),
}
