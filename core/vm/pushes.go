package vm

import (
	"github.com/holiman/uint256"

	"github.com/ledgerwatch/turbo-geth/common"
)

// opPush1 is a specialized version of pushN
func opPush1(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	var (
		codeLen = uint64(len(callContext.contract.Code))
		integer = new(uint256.Int)
	)
	*pc++
	if *pc < codeLen {
		callContext.stack.Push(integer.SetUint64(uint64(callContext.contract.Code[*pc])))
	} else {
		callContext.stack.Push(integer)
	}
	return nil, nil
}

func opPush2(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 2
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush3(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 3
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush4(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 4
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush5(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 5
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush6(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 6
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush7(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 7
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush8(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 8
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush9(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 9
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush10(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 10
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush11(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 11
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush12(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 12
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush13(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 13
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush14(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 14
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush15(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 15
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush16(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 16
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush17(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 17
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush18(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 18
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush19(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 19
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush20(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 20
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush21(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 21
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush22(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 22
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush23(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 23
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush24(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 24
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush25(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 25
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush26(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 26
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush27(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 27
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush28(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 28
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush29(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 29
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush30(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 30
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush31(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 31
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}

func opPush32(pc *uint64, _ *EVMInterpreter, callContext *callCtx) ([]byte, error) {
	pushByteSize := 32
	codeLen := len(callContext.contract.Code)

	startMin := int(*pc + 1)
	if startMin >= codeLen {
		startMin = codeLen
	}
	endMin := startMin + pushByteSize
	if startMin+pushByteSize >= codeLen {
		endMin = codeLen
	}

	integer := new(uint256.Int)
	callContext.stack.Push(common.SetBytesRightPadded(integer, callContext.contract.Code[startMin:endMin], pushByteSize))

	*pc += uint64(pushByteSize)
	return nil, nil
}
