package asm

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

type OpAssembler struct {
	assemblers      *stack.Stack[*ByteAssembler]
	blocks          *stack.Stack[blockType]
	functions       []opFunction
	currentFunction *opFunction
}

func (a *OpAssembler) asm() *ByteAssembler {
	return a.assemblers.Peek()
}

func (a *OpAssembler) defineArg(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.args[id] = opVar{a.currentFunction.argsBytes(), size}
	return nil
}

func (a *OpAssembler) defineRet(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.rets[id] = opVar{a.currentFunction.retsBytes(), size}
	return nil
}

func (a *OpAssembler) defineLocal(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.locals[id] = opVar{a.currentFunction.localsBytes(), size}
	return nil
}

func (a *OpAssembler) endFunction() error {
	if a.blocks.Pop() != functionBlock {
		return errors.New("expected function block")
	}

	a.functions = append(a.functions, *a.currentFunction)
	a.currentFunction = nil
	return nil
}

func (a *OpAssembler) endArgs() error {
	if a.blocks.Pop() != argsBlock {
		return errors.New("expected args block")
	}
	return nil
}

func (a *OpAssembler) endRets() error {
	if a.blocks.Pop() != retsBlock {
		return errors.New("expected rets block")
	}

	return nil
}

func (a *OpAssembler) endLocals() error {
	// TODO: Validate there's a current ongoing function.

	if a.blocks.Pop() != localsBlock {
		return errors.New("expected locals block")
	}

	if err := a.StackAlloc(a.currentFunction.localsBytes()); err != nil {
		return err
	}

	return nil
}

func (a *OpAssembler) endIf() error {
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

func (a *OpAssembler) endElse() error {
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

func (a *OpAssembler) StackAlloc(size uint16) error {
	a.asm().
		PutOpCode(vm.StackAlloc).
		PutI16(size)
	return nil
}

func (a *OpAssembler) Function(id string) error {
	// TODO: Validate there's no current ongoing function.

	a.blocks.Push(functionBlock)

	a.currentFunction = &opFunction{
		id,
		uint64(a.asm().Len()),
		map[string]opVar{}, /* args */
		map[string]opVar{}, /* rets */
		map[string]opVar{}, /* locals */
	}

	return nil
}

func (a *OpAssembler) Args() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(argsBlock)
	return nil
}

func (a *OpAssembler) Rets() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(retsBlock)
	return nil
}

func (a *OpAssembler) Locals() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(localsBlock)
	return nil
}

func (a *OpAssembler) DefineVar(id string, size uint16) error {
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

func (a *OpAssembler) IfThen() error {
	// TODO: Validate there's a current ongoing function.

	a.assemblers.Push(NewByteAssembler())
	a.blocks.Push(ifThenBlock)
	return nil
}

func (a *OpAssembler) IfElse() error {
	// TODO: Validate there's a current ongoing function.

	a.assemblers.Push(NewByteAssembler())
	a.blocks.Push(ifElseBlock)
	return nil
}

func (a *OpAssembler) Else() error {
	if a.blocks.Pop() != ifThenBlock {
		return errors.New("expected if block")
	}

	a.assemblers.Push(NewByteAssembler())
	a.blocks.Push(elseBlock)
	return nil
}

func (a *OpAssembler) End() error {
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

func (a *OpAssembler) PushI8(value byte) error {
	a.asm().
		PutOpCode(vm.PushI8).
		PutI8(value)
	return nil
}

func (a *OpAssembler) PushI16(value uint16) error {
	a.asm().
		PutOpCode(vm.PushI16).
		PutI16(value)
	return nil
}

func (a *OpAssembler) PushI32(value uint32) error {
	a.asm().
		PutOpCode(vm.PushI32).
		PutI32(value)
	return nil
}

func (a *OpAssembler) PushI64(value uint64) error {
	a.asm().
		PutOpCode(vm.PushI64).
		PutI64(value)
	return nil
}

func (a *OpAssembler) PushLocal(id string) error {
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

func (a *OpAssembler) PopLocal(id string) error {
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

func (a *OpAssembler) Add(typ OpType) error {
	opcodes := []uint64{
		I8:  vm.AddI8,
		I16: vm.AddI16,
		I32: vm.AddI32,
		I64: vm.AddI64,
	}
	a.asm().PutOpCode(opcodes[typ])
	return nil
}

func (a *OpAssembler) Print(typ OpType) error {
	opcodes := []uint64{
		I8:  vm.PrintI8,
		I16: vm.PrintI16,
		I32: vm.PrintI32,
		I64: vm.PrintI64,
	}
	a.asm().PutOpCode(opcodes[typ])
	return nil
}

func (a *OpAssembler) Program() vm.OpProgram {
	return vm.OpProgram{
		a.asm().Data(),
		[]vm.OpFunction{},
	}
}

func New() *OpAssembler {
	assembler := &OpAssembler{
		stack.New[*ByteAssembler](), /* assemblers */
		stack.New[blockType](),      /* blocks */
		[]opFunction{},
		nil, /* currentFunction */
	}
	assembler.assemblers.Push(NewByteAssembler())
	return assembler
}
