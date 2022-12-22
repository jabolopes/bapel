package vm

import (
	"fmt"
	"io"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Instruction set
//
// call <index:i16>
//
// pushi8 <value:i8>
// ...
//
// pushli8 <localOffset:i16>
// ...
//
// printi8
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

	PushLocalI8
	PushLocalI16
	PushLocalI32
	PushLocalI64

	PrintI8
	PrintI16
	PrintI32
	PrintI64
)

func opHalt(*Machine) error {
	return io.EOF
}

func opCall(machine *Machine) error {
	functionIndex := machine.Tape().GetI16()
	function := machine.program.Functions[functionIndex]

	stack := machine.Stack().
		PushI64(machine.pc).
		PushI64(machine.fp)
	machine.pc = function.Offset
	machine.fp = uint64(len(machine.stack))
	stack.Extend(function.Locals)
	return nil
}

func opReturn(machine *Machine) error {
	stack := machine.Stack()
	machine.stack = machine.stack[:machine.fp]
	machine.fp = stack.PopI64()
	machine.pc = stack.PopI64()
	return nil
}

func opStackAlloc(machine *Machine) error {
	size := machine.Tape().GetI16()

	stackSize := len(machine.stack)
	machine.Stack().Extend(uint64(size))

	fmt.Printf("DEBUG: stack alloc %d: sp:%d -> sp:%d\n", size, stackSize, len(machine.stack))
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
	offset := machine.Tape().GetI16()
	value := machine.Frame().LocalI8(uint64(offset))
	machine.Stack().PushI8(value)
	return nil
}

func opPushLocalI16(machine *Machine) error {
	offset := machine.Tape().GetI16()
	value := machine.Frame().LocalI16(uint64(offset))
	machine.Stack().PushI16(value)
	return nil
}

func opPushLocalI32(machine *Machine) error {
	offset := machine.Tape().GetI32()
	value := machine.Frame().LocalI32(uint64(offset))
	machine.Stack().PushI32(value)
	return nil
}

func opPushLocalI64(machine *Machine) error {
	offset := machine.Tape().GetI64()
	value := machine.Frame().LocalI64(uint64(offset))
	machine.Stack().PushI64(value)
	return nil
}

func opPrintI8(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopI8())
	return nil
}

func opPrintI16(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopI16())
	return nil
}

func opPrintI32(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopI32())
	return nil
}

func opPrintI64(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopI64())
	return nil
}
