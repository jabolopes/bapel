package vm

import (
	"fmt"
	"io"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Instruction set
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
	Call = OpCode(iota)
	Halt

	PushU8
	PushU16
	PushU32
	PushU64

	PushL8
	PushL16
	PushL32
	PushL64

	PrintU8
	PrintU16
	PrintU32
	PrintU64
)

// func opCall(machine *Machine) error {
// 	functionIndex := Tape{&machine.pc}.GetU16()
// 	function := machine.program.functions[functionIndex]

// 	Stack{&machine.stack}.PushU16(machine.pc)
// 	return nil
// }

func opHalt(*Machine) error {
	return io.EOF
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

func opPushLocalU8(machine *Machine) error {
	localOffset := Tape{machine.program.Data, &machine.pc}.GetU16()
	actualOffset := machine.locals + uint64(localOffset)
	machine.Stack().PushU8(machine.stack[actualOffset])
	return nil
}

func opPrintU8(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU8())
	return nil
}

func opPrintU16(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU16())
	return nil
}

func opPrintU32(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU32())
	return nil
}

func opPrintU64(machine *Machine) error {
	fmt.Printf("%d\n", machine.Stack().PopU64())
	return nil
}
