package ir

import (
	"fmt"
	"os"
)

type accum int

const (
	preAccum = accum(iota)
	postAccum
)

func computeOffsets(vars []IrVar, typ IrVarType, baseOffset int, sign int, accum accum) (int, error) {
	for i := range vars {
		irvar := &vars[i]
		if irvar.VarType != typ {
			continue
		}

		size, err := SizeOfType(irvar.Type)
		if err != nil {
			return 0, err
		}

		if accum == postAccum {
			irvar.offset = baseOffset * sign
			baseOffset += size
		} else {
			baseOffset += size
			irvar.offset = baseOffset * sign
		}
	}

	return baseOffset, nil
}

// irFunction is a function in IR.
//
// The frameSize and the localsSize are only computed once the args,
// rets, and locals sections are fully defined.
type irFunction struct {
	id     string  // Name of function.
	vars   []IrVar // Variables in the order in which they were defined.
	frame  irFrame
	offset uint64 // Offset of function relative to program data.
}

func (f *irFunction) lookupVar(id string) (IrVar, error) {
	for _, irvar := range f.vars {
		if irvar.Id == id {
			return irvar, nil
		}
	}

	return IrVar{}, fmt.Errorf("Undefined variable %q", id)
}

func (f *irFunction) addVar(id string, irvar IrVar) error {
	if _, err := f.lookupVar(id); err == nil {
		return fmt.Errorf("Variable %q already defined in this context", id)
	}

	size, err := SizeOfType(irvar.Type)
	if err != nil {
		return err
	}

	irvar.offset = int(f.frame.localsSize)

	f.vars = append(f.vars, irvar)
	f.frame.frameSize += uint16(size)
	f.frame.localsSize += uint16(size)
	return nil
}

func (f *irFunction) decl() irDecl {
	var args []IrIntType
	var rets []IrIntType
	for _, irvar := range f.vars {
		if irvar.VarType == ArgVar {
			args = append(args, irvar.Type)
		} else if irvar.VarType == RetVar {
			rets = append(rets, irvar.Type)
		}
	}

	return irDecl{f.id, FunctionDecl, 0, IrFunctionType{args, rets}}
}

func (f *irFunction) args() []IrVar {
	var args []IrVar
	for _, irvar := range f.vars {
		if irvar.VarType == ArgVar {
			args = append(args, irvar)
		}
	}

	return args
}

func (f *irFunction) rets() []IrVar {
	var rets []IrVar
	for _, irvar := range f.vars {
		if irvar.VarType == RetVar {
			rets = append(rets, irvar)
		}
	}

	return rets
}

// Call stack:
//   rets (reverse order)
//   args (reverse order)
//                           <- fp
//   locals (reverse order)
//   pc
//   fp                      <- sp
//
// Call protocol:
//
// The caller initiates the call protocol. To make the call, the
// caller is responsible for allocating (or pushing) the rets and the
// args in reverse order onto the stack. The callee is responsible for
// allocating (or pushing) the locals, the pc, and the fp onto the
// stack.
//
// To make the return, the callee is responsible for deallocating (or
// popping) the fp, the pc, the locals, and the args from the stack,
// leaving only the rets. The rets are then managed by the caller.
//
// Locals size (enter size):
//
// The function needs to know the size in bytes for the locals to
// allocate them on the stack.
//
// This size is not an operand to the 'call' opcode because that would
// leak the implementation details of the function and therefore
// whenever the function's frame size changed, all call sites would
// have to be recompiled, which is not possible.
//
// There's also no 'enter' or 'leave' opcodes, which would mean
// additional operations that every function had to call. Instead, the
// first i16 word of the function's body is the size of the
// locals. The 'call' opcode jumps to the function's pc and reads the
// i16 word directly from the function's op data, which also advances
// the pc to the first proper instruction of the function's body.
//
// PC handling:
//
// The pc is handled by the 'call' opcode. The 'call' opcode pushes
// the pc from the register onto the stack. The 'return' opcode pops
// the pc from the stack back to the register.
//
// FP handling:
//
// The fp is handled by the 'call' opcode. The caller does not push
// the fp to the stack, but the callee knows that the fp is equal to
// the sp at the start of the called function (i.e., in the 'call'
// opcode). The 'call' opcode pushes the fp onto the stack to save it
// for later, and the 'return' opcode pops it back to the fp register.
//
// Offsets:
//
// The following call stack shows how offsets are calculated. Note
// that the local1 does not have offset 0 but rather size of its
// type.
//
//   retn    -- offset(retn-1)   - sizeof(retn)
//   ...
//   ret2    -- offset(ret1)     - sizeof(ret2)
//   ret1    -- offset(argn)     - sizeof(ret1)
//   ...
//   argn    -- offset(argn-1)   - sizeof(argn)
//   arg2    -- offset(arg1)     - sizeof(arg2)
//   arg1    -- 0                - sizeof(arg1)
//   ...
//   local1  -- 0
//   local2  -- offset(local1)   + sizeof(local2)
//   localn  -- offset(localn-1) + sizeof(localn)
//
// Steps:
// 1. Caller allocates (or pushes) rets.
// 2. Caller allocates (or pushes) args.
// 3. Caller invokes 'call' with the offset of the callee's function.
// 4. Callee allocates locals (enter size is stored in function's body)
// 5. Callee runs.
// 6. Callee invokes 'return' with the frame size.
func (f *irFunction) computeFrame() error {
	const pcSize = 8

	baseOffsets := []int{
		ArgVar:   0,
		RetVar:   0,
		LocalVar: 0,
	}

	// Compute offset for locals.
	{
		var err error
		baseOffsets[LocalVar], err = computeOffsets(f.vars, LocalVar, baseOffsets[LocalVar], 1 /* sign */, postAccum)
		if err != nil {
			return err
		}
	}

	{
		var err error
		baseOffsets[ArgVar], err = computeOffsets(f.vars, ArgVar, 0 /* baseOffset */, -1 /* sign */, preAccum)
		if err != nil {
			return err
		}
	}

	{
		var err error
		baseOffsets[RetVar], err = computeOffsets(f.vars, RetVar, baseOffsets[ArgVar], -1 /* sign */, preAccum)
		if err != nil {
			return err
		}
	}

	f.frame = irFrame{uint16(baseOffsets[ArgVar] + baseOffsets[LocalVar]), uint16(baseOffsets[LocalVar])}

	fmt.Fprintf(os.Stderr, "DEBUG frame size:%d enter:%d\n", f.frame.frameSize, f.frame.localsSize)

	for _, irvar := range f.vars {
		var varType string
		switch irvar.VarType {
		case ArgVar:
			varType = "arg"
		case RetVar:
			varType = "ret"
		case LocalVar:
			varType = "local"
		}

		size, err := SizeOfType(irvar.Type)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "  %s %s %d +%d\n", varType, irvar.Id, irvar.offset, size)
	}

	return nil
}
