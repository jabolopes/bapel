package vm

import "fmt"

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

func opPushU8(machine *Machine) error {
	Stack{&machine.stack}.PushU8(machine.pc[0])
	machine.pc = machine.pc[1:]
	return nil
}

func opPushU16(machine *Machine) error {
	const size = 2
	Stack{&machine.stack}.PushN(machine.pc[:size])
	machine.pc = machine.pc[size:]
	return nil
}

func opPushU32(machine *Machine) error {
	const size = 4
	Stack{&machine.stack}.PushN(machine.pc[:size])
	machine.pc = machine.pc[size:]
	return nil
}

func opPushU64(machine *Machine) error {
	const size = 8
	Stack{&machine.stack}.PushN(machine.pc[:size])
	machine.pc = machine.pc[size:]
	return nil
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
