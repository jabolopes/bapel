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
	tape := machine.Tape()
	stack := machine.Stack()

	// Get opcode operands.
	pc := tape.GetI64()

	// Remember caller's registers
	callerPc := machine.pc
	callerFp := machine.fp
	callerSp := len(machine.stack)

	// Jump to new address.
	machine.pc = pc

	// Get the locals size from the function's pc.
	enterSize := tape.GetI16()

	// Allocate frame by reserving locals.
	stack.Extend(uint64(enterSize))

	// Set fp. The callee fp must point to the base of the locals.
	machine.fp = uint64(len(machine.stack))

	// Save caller's pc.
	fmt.Printf("DEBUG push %d = %d\n", len(machine.stack), callerPc)
	stack.PushI64(callerPc)

	// Save caller's fp.
	stack.PushI64(callerFp)

	{
		fmt.Printf("DEBUG call %d %d pc:%d fp:%d sp:%d", pc, enterSize, callerPc, callerFp, callerSp)
		fmt.Printf(" -> pc:%d fp:%d sp:%d\n", machine.pc, machine.fp, len(machine.stack))
	}
	return nil
}

func opReturn(machine *Machine) error {
	tape := machine.Tape()
	stack := machine.Stack()

	// Get opcode operands.
	leaveSize := uint64(tape.GetI16())

	calleePc := machine.pc
	calleeFp := machine.fp
	calleeSp := len(machine.stack)

	// Restore caller's fp.
	machine.fp = stack.PopI64()

	// Restore caller's pc.
	machine.pc = stack.PopI64()
	fmt.Printf("DEBUG pop %d = %d\n", len(machine.stack), machine.pc)

	// Deallocate frame by dropping locals and arguments.
	stack.Drop(leaveSize)

	{
		fmt.Printf("DEBUG return %d pc:%d fp:%d sp:%d", leaveSize, calleePc, calleeFp, calleeSp)
		fmt.Printf(" -> pc:%d fp:%d sp:%d\n", machine.pc, machine.fp, len(machine.stack))
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

func opPush(base OpCode) map[OpCode]func(*Machine) error {
	return map[OpCode]func(*Machine) error{
		// Immediate mode.
		unaryOpCode(base, ImmediateMode, I8): func(machine *Machine) error {
			machine.Stack().PushN(machine.Tape().GetN(1))
			return nil
		},
		unaryOpCode(base, ImmediateMode, I16): func(machine *Machine) error {
			machine.Stack().PushN(machine.Tape().GetN(2))
			return nil
		},
		unaryOpCode(base, ImmediateMode, I32): func(machine *Machine) error {
			machine.Stack().PushN(machine.Tape().GetN(4))
			return nil
		},
		unaryOpCode(base, ImmediateMode, I64): func(machine *Machine) error {
			machine.Stack().PushN(machine.Tape().GetN(8))
			return nil
		},
		// Var mode.
		unaryOpCode(base, VarMode, I8): func(machine *Machine) error {
			value := machine.Frame().VarI8(uint64(machine.Tape().GetI16()))
			machine.Stack().PushI8(value)
			return nil
		},
		unaryOpCode(base, VarMode, I16): func(machine *Machine) error {
			value := machine.Frame().VarI16(uint64(machine.Tape().GetI16()))
			machine.Stack().PushI16(value)
			return nil
		},
		unaryOpCode(base, VarMode, I32): func(machine *Machine) error {
			value := machine.Frame().VarI32(uint64(machine.Tape().GetI16()))
			machine.Stack().PushI32(value)
			return nil
		},
		unaryOpCode(base, VarMode, I64): func(machine *Machine) error {
			value := machine.Frame().VarI64(uint64(machine.Tape().GetI16()))
			machine.Stack().PushI64(value)
			return nil
		},
		// Stack mode.
		unaryOpCode(base, StackMode, I8):  func(machine *Machine) error { return errors.New("Unimplemented") },
		unaryOpCode(base, StackMode, I16): func(machine *Machine) error { return errors.New("Unimplemented") },
		unaryOpCode(base, StackMode, I32): func(machine *Machine) error { return errors.New("Unimplemented") },
		unaryOpCode(base, StackMode, I64): func(machine *Machine) error { return errors.New("Unimplemented") },
	}
}

// opPop pops from the stack.
//
// Immediate mode: unimplemented
//
// Var mode:
//   pop(i16 offset)
//
//   Pops a value from the stack and copies it to the given variable.
//
//   offset: variable to copy the popped value to, identified by its
//   offset relative to the fp.
//
// Stack mode:
//   pop()
//
//   Pops a value from the stack (and discards it).
//
//   No operands.
var opPop = []func(*Machine) error{
	// Immediate mode.
	func(machine *Machine) error { return errors.New("Unimplemented") },
	func(machine *Machine) error { return errors.New("Unimplemented") },
	func(machine *Machine) error { return errors.New("Unimplemented") },
	func(machine *Machine) error { return errors.New("Unimplemented") },
	// Var mode.
	func(machine *Machine) error {
		offset := machine.Tape().GetI16()
		machine.Frame().SetVarI8(uint64(offset), machine.Stack().PopI8())
		return nil
	},
	func(machine *Machine) error {
		offset := machine.Tape().GetI16()
		machine.Frame().SetVarI16(uint64(offset), machine.Stack().PopI16())
		return nil
	},
	func(machine *Machine) error {
		offset := machine.Tape().GetI16()
		machine.Frame().SetVarI32(uint64(offset), machine.Stack().PopI32())
		return nil
	},
	func(machine *Machine) error {
		offset := machine.Tape().GetI16()
		machine.Frame().SetVarI64(uint64(offset), machine.Stack().PopI64())
		return nil
	},
	// Stack mode.
	func(machine *Machine) error {
		_ = machine.Stack().PopI8()
		return nil
	},
	func(machine *Machine) error {
		_ = machine.Stack().PopI16()
		return nil
	},
	func(machine *Machine) error {
		_ = machine.Stack().PopI32()
		return nil
	},
	func(machine *Machine) error {
		_ = machine.Stack().PopI64()
		return nil
	},
}
