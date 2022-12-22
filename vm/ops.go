package vm

import (
	"fmt"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type OpCode = uint64

const (
	PushU8 = OpCode(iota)
	PushU16
	PushU32
	PushU64
	PrintU8
	PrintU16
	PrintU32
	PrintU64
)

func opPushGeneric[T constraints.Integer]() func(*Machine) error {
	var value T
	size := unsafe.Sizeof(value)
	return func(machine *Machine) error {
		Stack{&machine.stack}.PushN(machine.pc[:size])
		machine.pc = machine.pc[size:]
		return nil
	}
}

func opPrintU8(machine *Machine) error {
	stack := Stack{&machine.stack}
	fmt.Printf("%d\n", stack.PopU8())
	return nil
}

func opPrintU16(machine *Machine) error {
	stack := Stack{&machine.stack}
	fmt.Printf("%d\n", stack.PopU16())
	return nil
}

func opPrintU32(machine *Machine) error {
	stack := Stack{&machine.stack}
	fmt.Printf("%d\n", stack.PopU32())
	return nil
}

func opPrintU64(machine *Machine) error {
	stack := Stack{&machine.stack}
	fmt.Printf("%d\n", stack.PopU64())
	return nil
}
