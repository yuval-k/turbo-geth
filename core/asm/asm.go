// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Provides support for dealing with EVM assembly instructions (e.g., disassembling them).
package asm

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ledgerwatch/turbo-geth/core/vm"
)

// Iterator for disassembled EVM instructions
type instructionIterator struct {
	code     []byte
	pc       uint64
	arg      []byte
	op       vm.OpCode
	error    error
	started  bool
	previous []Command
}

type Command struct {
	PC  uint64
	Arg []byte
	Op  vm.OpCode
}

type Commands []Command

func (cs Commands) String() string {
	res := make([]string, len(cs))
	for i, c := range cs {
		res[i] = c.Op.String()
	}
	return strings.Join(res, "_")
}

// Create a new instruction iterator.
func NewInstructionIterator(code []byte) *instructionIterator {
	it := new(instructionIterator)
	it.code = code
	return it
}

// Returns true if there is a next instruction and moves on.
func (it *instructionIterator) Next() bool {
	if it.error != nil || uint64(len(it.code)) <= it.pc {
		// We previously reached an error or the end.
		return false
	}

	if it.started {
		prevArgs := append([]byte{}, it.arg...)
		it.previous = append(it.previous, Command{
			it.pc,
			prevArgs,
			it.op,
		})

		// Since the iteration has been already started we move to the next instruction.
		if it.arg != nil {
			it.pc += uint64(len(it.arg))
		}
		it.pc++
	} else {
		// We start the iteration from the first instruction.
		it.started = true
	}

	if uint64(len(it.code)) <= it.pc {
		// We reached the end.
		return false
	}

	it.op = vm.OpCode(it.code[it.pc])
	if it.op.IsPush() {
		a := uint64(it.op) - uint64(vm.PUSH1) + 1
		u := it.pc + 1 + a
		if uint64(len(it.code)) <= it.pc || uint64(len(it.code)) < u {
			it.error = fmt.Errorf("incomplete push instruction at %v", it.pc)
			return false
		}
		it.arg = it.code[it.pc+1 : u]
	} else {
		it.arg = nil
	}
	return true
}

// Up to previous JUMPDEST
func (it *instructionIterator) Previous(n int) []Command {
	if !it.started {
		return nil
	}
	start := len(it.previous) - n
	if start < 0 {
		start = 0
	}

	return it.previous[start:]
}

// Up to previous opCode
func (it *instructionIterator) PreviousBefore(n int, code vm.OpCode) []Command {
	if !it.started {
		return nil
	}
	start := len(it.previous) - n
	if start < 0 {
		start = 0
	}
	for i:=start;i<len(it.previous);i++ {
		if it.previous[i].Op == code {
			start = i+1
		}
	}

	return it.previous[start:]
}

func (it *instructionIterator) PreviousBeforeJump(n int) []Command {
	if !it.started {
		return nil
	}
	start := len(it.previous) - n
	if start < 0 {
		start = 0
	}

	for i:=start;i<len(it.previous);i++ {
		if it.previous[i].Op == vm.JUMP || it.previous[i].Op == vm.JUMPI {
			start = i+1
		}
	}

	return it.previous[start:]
}

func (it *instructionIterator) Last() Command {
	return it.previous[len(it.previous)-1]
}

// Returns any error that may have been encountered.
func (it *instructionIterator) Error() error {
	return it.error
}

// Returns the PC of the current instruction.
func (it *instructionIterator) PC() uint64 {
	return it.pc
}

// Returns the opcode of the current instruction.
func (it *instructionIterator) Op() vm.OpCode {
	return it.op
}

// Returns the argument of the current instruction.
func (it *instructionIterator) Arg() []byte {
	return it.arg
}

// Pretty-print all disassembled EVM instructions to stdout.
func PrintDisassembled(code string) error {
	script, err := hex.DecodeString(code)
	if err != nil {
		return err
	}
	return PrintDisassembledBytes(script)
}

// Pretty-print all disassembled EVM instructions to stdout.
func PrintDisassembledBytes(script []byte) error {
	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			fmt.Printf("%05x: %v 0x%x\n", it.PC(), it.Op(), it.Arg())
		} else {
			fmt.Printf("%05x: %v\n", it.PC(), it.Op())
		}
	}
	return it.Error()
}

// Pretty-print all disassembled EVM instructions to stdout.
func PrintDisassembledBytesUpTo(script []byte, toPC uint64) error {
	it := NewInstructionIterator(script)
	for it.Next() {
		if it.pc >= toPC {
			return nil
		}
		if it.Arg() != nil && 0 < len(it.Arg()) {
			fmt.Printf("%05x: %v 0x%x\n", it.PC(), it.Op(), it.Arg())
		} else {
			fmt.Printf("%05x: %v\n", it.PC(), it.Op())
		}
	}
	return it.Error()
}

// Return all disassembled EVM instructions in human-readable format.
func Disassemble(script []byte) ([]string, error) {
	instrs := make([]string, 0)

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			instrs = append(instrs, fmt.Sprintf("%05x: %v 0x%x\n", it.PC(), it.Op(), it.Arg()))
		} else {
			instrs = append(instrs, fmt.Sprintf("%05x: %v\n", it.PC(), it.Op()))
		}
	}
	if err := it.Error(); err != nil {
		return nil, err
	}
	return instrs, nil
}
