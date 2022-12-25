package vm

import (
	"errors"
	"fmt"
	"io"
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

func opIfThen(machine *Machine) error {
	endOffset := machine.Tape().GetI64()
	if machine.Stack().PopI8() == 0 {
		machine.pc += endOffset
	}
	return nil
}

func opIfElse(machine *Machine) error {
	endOffset := machine.Tape().GetI64()
	if machine.Stack().PopI8() != 0 {
		machine.pc += endOffset
	}
	return nil
}

func opElse(machine *Machine) error {
	machine.pc += machine.Tape().GetI64()
	return nil
}

func opStackAlloc(machine *Machine) error {
	size := machine.Tape().GetI16()

	stackSize := len(machine.stack)
	machine.Stack().Extend(uint64(size))

	fmt.Printf("DEBUG: stack alloc %d: sp:%d -> sp:%d\n", size, stackSize, len(machine.stack))
	return nil
}

func opPopVarI8(machine *Machine) error {
	offset := machine.Tape().GetI16()
	machine.Frame().SetVarI8(uint64(offset), machine.Stack().PopI8())
	return nil
}

func opPopVarI16(machine *Machine) error {
	offset := machine.Tape().GetI16()
	machine.Frame().SetVarI16(uint64(offset), machine.Stack().PopI16())
	return nil
}

func opPopVarI32(machine *Machine) error {
	offset := machine.Tape().GetI16()
	machine.Frame().SetVarI32(uint64(offset), machine.Stack().PopI32())
	return nil
}

func opPopVarI64(machine *Machine) error {
	offset := machine.Tape().GetI16()
	machine.Frame().SetVarI64(uint64(offset), machine.Stack().PopI64())
	return nil
}

func opAddI8(machine *Machine) error {
	stack := machine.Stack()
	stack.PushI8(stack.PopI8() + stack.PopI8())
	return nil
}

func opAddI16(machine *Machine) error {
	stack := machine.Stack()
	stack.PushI16(stack.PopI16() + stack.PopI16())
	return nil
}

func opAddI32(machine *Machine) error {
	stack := machine.Stack()
	stack.PushI32(stack.PopI32() + stack.PopI32())
	return nil
}

func opAddI64(machine *Machine) error {
	stack := machine.Stack()
	stack.PushI64(stack.PopI64() + stack.PopI64())
	return nil
}

var opPrint = []func(*Machine) error{
	// Immediate mode.
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Tape().GetI8())
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Tape().GetI16())
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Tape().GetI32())
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Tape().GetI64())
		return nil
	},
	// Var mode.
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Frame().VarI8(uint64(machine.Tape().GetI16())))
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Frame().VarI16(uint64(machine.Tape().GetI16())))
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Frame().VarI32(uint64(machine.Tape().GetI16())))
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Frame().VarI64(uint64(machine.Tape().GetI16())))
		return nil
	},
	// Stack mode.
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Stack().PopI8())
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Stack().PopI16())
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Stack().PopI32())
		return nil
	},
	func(machine *Machine) error {
		fmt.Printf("%d\n", machine.Stack().PopI64())
		return nil
	},
}

var opPush = []func(*Machine) error{
	// Immediate mode.
	func(machine *Machine) error {
		machine.Stack().PushN(machine.Tape().GetN(1))
		return nil
	},
	func(machine *Machine) error {
		machine.Stack().PushN(machine.Tape().GetN(2))
		return nil
	},
	func(machine *Machine) error {
		machine.Stack().PushN(machine.Tape().GetN(4))
		return nil
	},
	func(machine *Machine) error {
		machine.Stack().PushN(machine.Tape().GetN(8))
		return nil
	},
	// Var mode.
	func(machine *Machine) error {
		value := machine.Frame().VarI8(uint64(machine.Tape().GetI16()))
		machine.Stack().PushI8(value)
		return nil
	},
	func(machine *Machine) error {
		value := machine.Frame().VarI16(uint64(machine.Tape().GetI16()))
		machine.Stack().PushI16(value)
		return nil
	},
	func(machine *Machine) error {
		value := machine.Frame().VarI32(uint64(machine.Tape().GetI16()))
		machine.Stack().PushI32(value)
		return nil
	},
	func(machine *Machine) error {
		value := machine.Frame().VarI64(uint64(machine.Tape().GetI16()))
		machine.Stack().PushI64(value)
		return nil
	},
	// Stack mode.
	func(machine *Machine) error { return errors.New("Unimplemented") },
	func(machine *Machine) error { return errors.New("Unimplemented") },
	func(machine *Machine) error { return errors.New("Unimplemented") },
	func(machine *Machine) error { return errors.New("Unimplemented") },
}
