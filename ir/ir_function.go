package ir

import (
	"fmt"
)

func computeOffsets(vars []IrVar, typ IrVarType, baseOffset int) (int, error) {
	for i := range vars {
		irvar := &vars[i]
		if irvar.VarType != typ {
			continue
		}

		irvar.offset = uint16(baseOffset)

		size, err := SizeOfType(irvar.Type)
		if err != nil {
			return 0, err
		}

		baseOffset += size
	}

	return baseOffset, nil
}

// IrFunction is a function in IR.
//
// The frameSize and the localsSize are only computed once the args,
// rets, and locals sections are fully defined.
type IrFunction struct {
	id     string  // Name of function.
	vars   []IrVar // Variables in the order in which they were defined.
	frame  irFrame
	offset uint64 // Offset of function relative to program data.
}

func (f *IrFunction) lookupVar(id string) (IrVar, error) {
	for _, irvar := range f.vars {
		if irvar.Id == id {
			return irvar, nil
		}
	}

	return IrVar{}, fmt.Errorf("Undefined variable %q", id)
}

func (f *IrFunction) addVar(id string, irvar IrVar) error {
	if _, err := f.lookupVar(id); err == nil {
		return fmt.Errorf("Variable %q already defined in this context", id)
	}

	f.vars = append(f.vars, irvar)
	return nil
}

// Call stack (grows upwards)
//   rets (reverse order)
//   pc
//   args (reverse order)
//   locals (reverse order)
//
// FP handling:
//   The caller does not push the fp to the stack, but the callee
//   knows that the sp is the fp at the start of the called
//   function. The callee can push the fp to the stack to save it for
//   later, and pop back to the fp just before returning to the
//   caller.
//
//   push fp
//   fp <- sp - 8 (subtract the effect of pushing the fp)
//   (execute function body)
//   pop fp
//
// Example:
//   ...
//   ret2
//   ret1
//   pc
//   arg2
//   arg1
//   ...
//   local2
//   local1
//   ...
//
// Offsets:
//   Note: At the start of the function, the fp is the sp.
//
//   offset(Local, n) = sizeIndexes(Local, [1:n-1])
//   offset(Arg, n) = offset(Local, n) + sizeIndexes(Arg, [1:n-1])
//   offset(Ret, n) = offset(Arg, n) + size(pc) + sizeIndexes(Ret, [1:n-1])
func (f *IrFunction) computeFrame() error {
	const pcSize = 8

	baseOffsets := []int{
		ArgVar:   0,
		RetVar:   0,
		LocalVar: 0,
	}

	// Compute offset for locals.
	{
		var err error
		baseOffsets[LocalVar], err = computeOffsets(f.vars, LocalVar, baseOffsets[LocalVar])
		if err != nil {
			return err
		}
	}

	{
		var err error
		baseOffsets[ArgVar], err = computeOffsets(f.vars, ArgVar, baseOffsets[LocalVar])
		if err != nil {
			return err
		}
	}

	{
		var err error
		baseOffsets[RetVar], err = computeOffsets(f.vars, RetVar, baseOffsets[ArgVar])
		if err != nil {
			return err
		}
	}

	f.frame = irFrame{uint16(baseOffsets[ArgVar]), uint16(baseOffsets[LocalVar])}

	for _, irvar := range f.vars {
		if irvar.VarType == LocalVar {
			size, err := SizeOfType(irvar.Type)
			if err != nil {
				return err
			}

			fmt.Printf("DEBUG: %s %d %d\n", irvar.Id, irvar.offset, size)
		}
	}

	for _, irvar := range f.vars {
		if irvar.VarType == ArgVar {
			size, err := SizeOfType(irvar.Type)
			if err != nil {
				return err
			}

			fmt.Printf("DEBUG: %s %d %d\n", irvar.Id, irvar.offset, size)
		}
	}

	for _, irvar := range f.vars {
		if irvar.VarType == RetVar {
			size, err := SizeOfType(irvar.Type)
			if err != nil {
				return err
			}

			fmt.Printf("DEBUG: %s %d %d\n", irvar.Id, irvar.offset, size)
		}
	}

	return nil
}

func (f *IrFunction) Vars() []IrVar { return f.vars }
