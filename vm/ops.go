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
	{
		fmt.Printf("DEBUG call pc:%d", machine.pc)
	}

	pc := machine.Tape().GetI64()
	machine.Stack().PushI64(machine.pc)
	machine.pc = pc

	{
		fmt.Printf(" -> pc:%d\n", machine.pc)
	}
	return nil
}

func opReturn(machine *Machine) error {
	{
		fmt.Printf("DEBUG return pc:%d", machine.pc)
	}

	machine.pc = machine.Stack().PopI64()

	{
		fmt.Printf(" -> pc:%d\n", machine.pc)
	}
	return nil
}

func opEnter(machine *Machine) error {
	// Allocate space in stack for locals.
	stack := machine.Stack()

	enterSize := uint64(machine.Tape().GetI16())
	{
		fmt.Printf("DEBUG enter %d sp:%d fp:%d", enterSize, len(machine.stack), machine.fp)
	}

	stack.Extend(enterSize)

	// Set fp (saving caller's fp also).
	callerFp := machine.fp
	machine.fp = uint64(len(machine.stack))
	stack.PushI64(callerFp)

	{
		fmt.Printf(" -> sp:%d fp:%d\n", len(machine.stack), machine.fp)
	}

	return nil
}

func opLeave(machine *Machine) error {
	{
		fmt.Printf("DEBUG leave sp:%d fp:%d", len(machine.stack), machine.fp)
	}

	// Restore caller's fp.
	stack := machine.Stack()
	machine.fp = stack.PopI64()

	// Deallocate stack space for locals and also arguments.
	stack.Drop(uint64(machine.Tape().GetI16()))

	{
		fmt.Printf(" -> sp:%d fp:%d\n", len(machine.stack), machine.fp)
	}

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
