package vm

import (
	"fmt"
	"io"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Instruction set
//
// call <index:u16>
//
// pushu8 <value:u8>
// ...
//
// pushl8 <localOffset:u16>
// ...
//
// print8
// ...
type OpCode = uint64

const (
	Halt = OpCode(iota)

	Call
	Return

	StackAlloc

	PushI8
	PushI16
	PushI32
	PushI64

	PushL8
	PushL16
	PushL32
	PushL64

	PrintI8
	PrintI16
	PrintI32
	PrintI64
)

func opHalt(*Machine) error {
	return io.EOF
}

func opCall(machine *Machine) error {
	functionIndex := machine.Tape().GetU16()
	function := machine.program.Functions[functionIndex]

	stack := machine.Stack().
		PushU64(machine.pc).
		PushU64(machine.locals)
	machine.pc = function.Offset
	machine.locals = uint64(len(machine.stack))
	stack.Extend(function.Locals)
	return nil
}

func opReturn(machine *Machine) error {
	stack := machine.Stack()
	machine.stack = machine.stack[:machine.locals]
	machine.locals = stack.PopU64()
	machine.pc = stack.PopU64()
	return nil
}

func opStackAlloc(machine *Machine) error {
	size := machine.Tape().GetU16()
	machine.Stack().Extend(uint64(size))
	return nil
}

func opPushImmediate[T constraints.Integer]() func(*Machine) error {
	var value T
	size := uint64(unsafe.Sizeof(value))
	return func(machine *Machine) error {
		value := machine.Tape().GetN(size)
		machine.Stack().PushN(value)
		return nil
	}
}

func opPushLocalI8(machine *Machine) error {
	localOffset := Tape{machine.program.Data, &machine.pc}.GetU16()
	actualOffset := machine.locals + uint64(localOffset)
	machine.Stack().PushU8(machine.stack[actualOffset])
	return nil
}

func opPrintI8(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU8())
	return nil
}

func opPrintI16(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU16())
	return nil
}

func opPrintI32(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU32())
	return nil
}

func opPrintI64(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU64())
	return nil
}
