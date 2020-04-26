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
type InstructionIterator struct {
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

func (c Command) String() string {
	var res string
	if c.Arg != nil && 0 < len(c.Arg) {
		res += fmt.Sprintf("%05x: %v 0x%x\n", c.PC, c.Op, c.Arg)
	} else {
		res += fmt.Sprintf("%05x: %v\n", c.PC, c.Op)
	}
	return res
}

// Create a new instruction iterator.
func NewInstructionIterator(code []byte) *InstructionIterator {
	it := new(InstructionIterator)
	it.code = code
	return it
}

// Returns true if there is a next instruction and moves on.
func (it *InstructionIterator) Next() bool {
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

func OpsToString(ops []vm.OpCode) string {
	str := ""
	for _, op := range ops {
		str = fmt.Sprintf("%s_%s", str, op.String())
	}
	return str
}

func (it *InstructionIterator) HasPrefix(prefix []vm.OpCode) ([]vm.OpCode, bool) {
	var actualPrefix []vm.OpCode
	for _, op := range prefix {
		if it.started {
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
			return actualPrefix, false
		}

		it.op = vm.OpCode(it.code[it.pc])
		actualPrefix = append(actualPrefix, it.op)

		if it.op.IsPush() {
			a := uint64(it.op) - uint64(vm.PUSH1) + 1
			u := it.pc + 1 + a
			if uint64(len(it.code)) <= it.pc || uint64(len(it.code)) < u {
				it.error = fmt.Errorf("incomplete push instruction at %v", it.pc)
				return actualPrefix, false
			}
			it.arg = it.code[it.pc+1 : u]
		} else {
			it.arg = nil
		}

		if it.op != op {
			return actualPrefix, false
		}
	}

	it.op = 0
	it.arg = nil
	it.pc = 0
	it.previous = nil
	it.error = nil
	it.started = false

	return nil, true
}

func (it *InstructionIterator) Previous(n int) []Command {
	if !it.started {
		return nil
	}
	start := len(it.previous) - n
	if start < 0 {
		start = 0
	}

	return it.previous[start:]
}

func (it *InstructionIterator) PreviousOpCode() *Command {
	if !it.started {
		return nil
	}
	return &it.previous[len(it.previous)-1]
}

// Up to previous OpCode
func (it *InstructionIterator) PreviousOp(n int, code vm.OpCode) {
	if !it.started {
		return
	}

	count := 0
	var prev Command
	for i := len(it.previous) - 1; i >= 0; i-- {
		prev = it.previous[len(it.previous)-1]
		it.previous = it.previous[:len(it.previous)-1]

		it.pc = prev.PC
		it.arg = prev.Arg
		it.op = prev.Op

		if it.previous[i].Op == code {
			count++
			if count == n {
				return
			}
		}
	}

	return
}

/*
func (it *InstructionIterator) PreviousUntilStackValue(stackV int) ([]Command, bool) {
	if !it.started {
		return nil, false
	}

	var (
		count int
		found bool
		prev  Command
		cmd   vm.Command
	)
	history := make([]Command, 0, 5)

loop:
	for i := len(it.previous) - 1; i >= 0; i-- {
		prev = it.previous[i]
		cmd = vm.Commands[prev.Op]

		history = append(history, prev)

		if prev.Op == vm.JUMPDEST || IsStop(prev.Op) {
			found = false
			break loop
		}

		count = count - cmd.In + cmd.Out
		if stackV > 0 {
			if count >= stackV {
				found = true
				break loop
			}
		}
		if stackV == 0 {
			if count == stackV {
				found = true
				break loop
			}
		}
		if stackV < 0 {
			if count < stackV {
				found = true
				break loop
			}
		}
	}

	for i := len(history)/2 - 1; i >= 0; i-- {
		opp := len(history) - 1 - i
		history[i], history[opp] = history[opp], history[i]
	}

	return history, found
}

 */

func IsDup(op vm.OpCode) bool {
	switch op {
	case vm.DUP1, vm.DUP2, vm.DUP3, vm.DUP4, vm.DUP5, vm.DUP6, vm.DUP7, vm.DUP8, vm.DUP9, vm.DUP10, vm.DUP11, vm.DUP12, vm.DUP13, vm.DUP14, vm.DUP15, vm.DUP16:
		return true
	}
	return false
}

func (it *InstructionIterator) Current() *Command {
	if !it.started {
		return nil
	}

	return &Command{
		it.pc,
		it.arg,
		it.op,
	}
}

// Up to previous opCode
func (it *InstructionIterator) PreviousBefore(n int, code vm.OpCode) []Command {
	if !it.started {
		return nil
	}
	start := len(it.previous) - n
	if start < 0 {
		start = 0
	}
	for i := start; i < len(it.previous); i++ {
		if it.previous[i].Op == code {
			start = i + 1
		}
	}

	return it.previous[start:]
}

func (it *InstructionIterator) PreviousBeforeJump(n int) []Command {
	if !it.started {
		return nil
	}
	start := len(it.previous) - n
	if start < 0 {
		start = 0
	}

	for i := start; i < len(it.previous); i++ {
		if it.previous[i].Op == vm.JUMP || it.previous[i].Op == vm.JUMPI {
			start = i + 1
		}
	}

	return it.previous[start:]
}

func (it *InstructionIterator) Last() *Command {
	idx := len(it.previous) - 1
	if idx < 0 {
		return nil
	}
	return &it.previous[idx]
}

func (it *InstructionIterator) Prev(n int) *Command {
	idx := len(it.previous) - n
	if idx < 0 {
		return nil
	}
	return &it.previous[len(it.previous)-1]
}

// Returns any error that may have been encountered.
func (it *InstructionIterator) Error() error {
	return it.error
}

// Returns the PC of the current instruction.
func (it *InstructionIterator) PC() uint64 {
	return it.pc
}

// Returns the opcode of the current instruction.
func (it *InstructionIterator) Op() vm.OpCode {
	return it.op
}

// Returns the argument of the current instruction.
func (it *InstructionIterator) Arg() []byte {
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
func Disassembled(code string) (string, error) {
	script, err := hex.DecodeString(code)
	if err != nil {
		return "", err
	}
	return DisassembledBytes(script)
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
func DisassembledBytes(script []byte) (string, error) {
	var res string
	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			res += fmt.Sprintf("%05x: %v 0x%x\n", it.PC(), it.Op(), it.Arg())
		} else {
			res += fmt.Sprintf("%05x: %v\n", it.PC(), it.Op())
		}
	}
	return res, it.Error()
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
