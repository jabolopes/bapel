package vm

import (
	"fmt"
	"io"

	"github.com/jabolopes/bapel/ir"
	"golang.org/x/exp/constraints"
)

// opHalt halts the program.
//
// No operands.
func opHalt(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error { return io.EOF },
	}
}

// opCall calls a function at a given offset.
//
// call(i64 immediate offset)
//
// offset: function to call identified by its absolute offset. The
// offset is an index in the program data.
func opCall(base ir.OpCode) opFamilyMap {
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
				fmt.Printf("DEBUG call %d %d pc:%d fp:%d sp:%d", pc, enterSize, callerPc, callerFp, callerSp)
				fmt.Printf(" -> pc:%d fp:%d sp:%d\n", machine.pc, machine.fp, len(machine.stack))
			}
			return nil
		},
	}
}

// opReturn returns from a function.
//
// return(i16 immediate leaveSize)
//
// leaveSize: size in bytes to deallocate from the stack. This size
// includes the size of locals and args.
func opReturn(base ir.OpCode) opFamilyMap {
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
				fmt.Printf("DEBUG return %d pc:%d fp:%d sp:%d", leaveSize, calleePc, calleeFp, calleeSp)
				fmt.Printf(" -> pc:%d fp:%d sp:%d\n", machine.pc, machine.fp, len(machine.stack))
			}
			return nil
		},
	}
}

// opIfThen tests whether the value at the top of the stack is zero
// and if that is the case the pc is incremented by the value given by
// the operand. The value at the top of the stack is i8.
//
// ifThen(i64 immediate offset)
//
// offset: if the value
func opIfThen(base ir.OpCode) opFamilyMap {
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

func opIfElse(base ir.OpCode) opFamilyMap {
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

func opElse(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			machine.pc += machine.Tape().GetI64()
			return nil
		},
	}
}

func opPrintImpl[T constraints.Integer](value T) {
	fmt.Printf("%d\n", value)
}
