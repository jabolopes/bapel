package vm

import (
	"errors"
	"fmt"
	"os"
)

var errHalt = errors.New("HALT")

func opHalt(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error { return errHalt },
	}
}

func opCall(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			tape := machine.Tape()
			stack := machine.Stack()

			// Get opcode operands.
			pc := tape.GetI64()

			// Remember caller's registers
			callerPc := machine.pc
			callerFp := machine.fp
			callerSp := len(machine.stack)

			// Set fp. The callee fp must point to the base of the args.
			machine.fp = uint64(len(machine.stack))

			// Jump to new address.
			machine.pc = pc

			// Get the locals size from the function's pc.
			enterSize := tape.GetI16()

			// Allocate frame by reserving locals.
			stack.Extend(uint64(enterSize))

			// Save caller's pc.
			stack.PushI64(callerPc)

			// Save caller's fp.
			stack.PushI64(callerFp)

			{
				fmt.Fprintf(os.Stderr, "DEBUG call %d %d pc:%d fp:%d sp:%d", pc, enterSize, callerPc, callerFp, callerSp)
				fmt.Fprintf(os.Stderr, " -> pc:%d fp:%d sp:%d\n", machine.pc, machine.fp, len(machine.stack))
			}
			return nil
		},
	}
}

func opReturn(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
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

			// Deallocate frame by dropping locals and arguments.
			stack.Drop(leaveSize)

			{
				fmt.Fprintf(os.Stderr, "DEBUG return %d pc:%d fp:%d sp:%d", leaveSize, calleePc, calleeFp, calleeSp)
				fmt.Fprintf(os.Stderr, " -> pc:%d fp:%d sp:%d\n", machine.pc, machine.fp, len(machine.stack))
			}
			return nil
		},
	}
}

func opIfThen(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			endOffset := machine.Tape().GetI64()
			if machine.Stack().PopI8() == 0 {
				machine.pc += endOffset
			}
			return nil
		},
	}
}

func opIfElse(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			endOffset := machine.Tape().GetI64()
			if machine.Stack().PopI8() != 0 {
				machine.pc += endOffset
			}
			return nil
		},
	}
}

func opElse(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			machine.pc += machine.Tape().GetI64()
			return nil
		},
	}
}
