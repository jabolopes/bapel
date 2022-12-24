package ir

import (
	"errors"
	"fmt"

	"github.com/jabolopes/bapel/vm"
	"github.com/zyedidia/generic/stack"
)

type blockType int

const (
	functionBlock = blockType(iota)
	argsBlock
	retsBlock
	localsBlock
	ifThenBlock
	ifElseBlock
	elseBlock
)

type IrAssembler struct {
	assemblers      *stack.Stack[*ByteAssembler]
	blocks          *stack.Stack[blockType]
	functions       []irFunction
	currentFunction *irFunction
}

func (a *IrAssembler) asm() *ByteAssembler {
	return a.assemblers.Peek()
}

func (a *IrAssembler) defineArg(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.args[id] = irVar{a.currentFunction.argsBytes(), size}
	return nil
}

func (a *IrAssembler) defineRet(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.rets[id] = irVar{a.currentFunction.retsBytes(), size}
	return nil
}

func (a *IrAssembler) defineLocal(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.locals[id] = irVar{a.currentFunction.localsBytes(), size}
	return nil
}

func (a *IrAssembler) endFunction() error {
	if a.blocks.Pop() != functionBlock {
		return errors.New("expected function block")
	}

	a.functions = append(a.functions, *a.currentFunction)
	a.currentFunction = nil
	return nil
}

func (a *IrAssembler) endArgs() error {
	if a.blocks.Pop() != argsBlock {
		return errors.New("expected args block")
	}
	return nil
}

func (a *IrAssembler) endRets() error {
	if a.blocks.Pop() != retsBlock {
		return errors.New("expected rets block")
	}

	return nil
}

func (a *IrAssembler) endLocals() error {
	// TODO: Validate there's a current ongoing function.

	if a.blocks.Pop() != localsBlock {
		return errors.New("expected locals block")
	}

	if err := a.StackAlloc(a.currentFunction.localsBytes()); err != nil {
		return err
	}

	return nil
}

func (a *IrAssembler) endIf() error {
	block := a.blocks.Pop()
	if block != ifThenBlock && block != ifElseBlock {
		return errors.New("expected if block")
	}

	nested := a.assemblers.Pop()

	if block == ifThenBlock {
		a.asm().PutOpCode(vm.IfThen)
	} else {
		a.asm().PutOpCode(vm.IfElse)
	}

	a.asm().
		PutI64(uint64(len(nested.Data()))).
		append(nested.Data())
	return nil
}

func (a *IrAssembler) endElse() error {
	if a.blocks.Pop() != elseBlock {
		return errors.New("expected else block")
	}

	elseAsm := a.assemblers.Pop()

	// Finish the 'if' section by adding the 'else' (aka jump) instruction. This
	// is important so that the length of the 'if' section is correct when jumping
	// to the else branch.
	a.asm().
		PutOpCode(vm.Else).
		PutI64(uint64(len(elseAsm.Data())))

	ifAsm := a.assemblers.Pop()

	a.asm().
		PutOpCode(vm.IfThen).
		PutI64(uint64(len(ifAsm.Data()))).
		append(ifAsm.Data()).
		append(elseAsm.Data())

	return nil
}

func (a *IrAssembler) StackAlloc(size uint16) error {
	a.asm().
		PutOpCode(vm.StackAlloc).
		PutI16(size)
	return nil
}

func (a *IrAssembler) Function(id string) error {
	// TODO: Validate there's no current ongoing function.

	a.blocks.Push(functionBlock)

	a.currentFunction = &irFunction{
		id,
		uint64(a.asm().Len()),
		map[string]irVar{}, /* args */
		map[string]irVar{}, /* rets */
		map[string]irVar{}, /* locals */
	}

	return nil
}

func (a *IrAssembler) Args() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(argsBlock)
	return nil
}

func (a *IrAssembler) Rets() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(retsBlock)
	return nil
}

func (a *IrAssembler) Locals() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(localsBlock)
	return nil
}

func (a *IrAssembler) DefineVar(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	switch block := a.blocks.Peek(); block {
	case argsBlock:
		return a.defineArg(id, size)
	case retsBlock:
		return a.defineRet(id, size)
	case localsBlock:
		return a.defineLocal(id, size)
	default:
		return fmt.Errorf("Cannot declare variable inside block type %d", block)
	}
}

func (a *IrAssembler) IfThen() error {
	// TODO: Validate there's a current ongoing function.

	a.assemblers.Push(NewByteAssembler())
	a.blocks.Push(ifThenBlock)
	return nil
}

func (a *IrAssembler) IfElse() error {
	// TODO: Validate there's a current ongoing function.

	a.assemblers.Push(NewByteAssembler())
	a.blocks.Push(ifElseBlock)
	return nil
}

func (a *IrAssembler) Else() error {
	if a.blocks.Pop() != ifThenBlock {
		return errors.New("expected if block")
	}

	a.assemblers.Push(NewByteAssembler())
	a.blocks.Push(elseBlock)
	return nil
}

func (a *IrAssembler) End() error {
	switch block := a.blocks.Peek(); block {
	case functionBlock:
		return a.endFunction()
	case argsBlock:
		return a.endArgs()
	case retsBlock:
		return a.endRets()
	case localsBlock:
		return a.endLocals()
	case ifThenBlock, ifElseBlock:
		return a.endIf()
	case elseBlock:
		return a.endElse()
	default:
		return fmt.Errorf("Unexpected block type %d", block)
	}
}

func (a *IrAssembler) PushImmediate(typ IrType, value uint64) error {
	// TODO: Validate whether typecast truncates the value and return an
	// error in that case.

	switch typ {
	case I8:
		a.asm().
			PutOpCode(vm.PushI8).
			PutI8(byte(value))
	case I16:
		a.asm().
			PutOpCode(vm.PushI16).
			PutI16(uint16(value))
	case I32:
		a.asm().
			PutOpCode(vm.PushI32).
			PutI32(uint32(value))
	case I64:
		a.asm().
			PutOpCode(vm.PushI64).
			PutI64(value)
	default:
		return fmt.Errorf("Unhandled optype %d", typ)
	}
	return nil
}

func (a *IrAssembler) PushLocal(id string) error {
	// TODO: Validate there's a current ongoing function.

	local, ok := a.currentFunction.locals[id]
	if !ok {
		return fmt.Errorf("Undeclared local %q", id)
	}

	switch local.size {
	case 1:
		a.asm().PutOpCode(vm.PushLocalI8)
	case 2:
		a.asm().PutOpCode(vm.PushLocalI16)
	case 4:
		a.asm().PutOpCode(vm.PushLocalI32)
	case 8:
		a.asm().PutOpCode(vm.PushLocalI64)
	}

	a.asm().PutI16(local.offset)
	return nil
}

func (a *IrAssembler) PopLocal(id string) error {
	// TODO: Validate there's a current ongoing function.

	local, ok := a.currentFunction.locals[id]
	if !ok {
		return fmt.Errorf("Undeclared local %q", id)
	}

	switch local.size {
	case 1:
		a.asm().PutOpCode(vm.PopLocalI8)
	case 2:
		a.asm().PutOpCode(vm.PopLocalI16)
	case 4:
		a.asm().PutOpCode(vm.PopLocalI32)
	case 8:
		a.asm().PutOpCode(vm.PopLocalI64)
	}

	a.asm().PutI16(local.offset)
	return nil
}

func (a *IrAssembler) Add(typ IrType) error {
	opcodes := []uint64{
		I8:  vm.AddI8,
		I16: vm.AddI16,
		I32: vm.AddI32,
		I64: vm.AddI64,
	}
	a.asm().PutOpCode(opcodes[typ])
	return nil
}

func (a *IrAssembler) Print(typ IrType) error {
	opcodes := []uint64{
		I8:  vm.PrintI8,
		I16: vm.PrintI16,
		I32: vm.PrintI32,
		I64: vm.PrintI64,
	}
	a.asm().PutOpCode(opcodes[typ])
	return nil
}

func (a *IrAssembler) Program() vm.OpProgram {
	return vm.OpProgram{
		a.asm().Data(),
		[]vm.OpFunction{},
	}
}

func New() *IrAssembler {
	assembler := &IrAssembler{
		stack.New[*ByteAssembler](), /* assemblers */
		stack.New[blockType](),      /* blocks */
		[]irFunction{},
		nil, /* currentFunction */
	}
	assembler.assemblers.Push(NewByteAssembler())
	return assembler
}
