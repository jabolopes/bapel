package asm

import (
	"errors"
	"fmt"

	"github.com/jabolopes/go-vm/vm"
	"github.com/zyedidia/generic/stack"
)

type OpVar struct {
	offset uint16
	size   uint16
}

type blockType int

const (
	functionBlock = blockType(iota)
	argsBlock
	retsBlock
	localsBlock
	ifTrueBlock
	ifFalseBlock
)

type OpFunction struct {
	id     string
	offset uint64
	args   map[string]OpVar
	rets   map[string]OpVar
	locals map[string]OpVar
}

func (f OpFunction) ArgsBytes() uint16 {
	var size uint16
	for _, arg := range f.args {
		size += arg.size
	}
	return size
}

func (f OpFunction) RetsBytes() uint16 {
	var size uint16
	for _, arg := range f.rets {
		size += arg.size
	}
	return size
}

func (f OpFunction) LocalsBytes() uint16 {
	var size uint16
	for _, local := range f.locals {
		size += local.size
	}
	return size
}

type OpAssembler struct {
	assemblers      *stack.Stack[*Assembler]
	blocks          *stack.Stack[blockType]
	functions       []OpFunction
	currentFunction *OpFunction
}

func (a *OpAssembler) asm() *Assembler {
	return a.assemblers.Peek()
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

	a.currentFunction = &OpFunction{
		id,
		uint64(a.asm().Len()),
		map[string]OpVar{}, /* args */
		map[string]OpVar{}, /* rets */
		map[string]OpVar{}, /* locals */
	}

	return nil
}

func (a *OpAssembler) EndFunction() error {
	if a.blocks.Pop() != functionBlock {
		return errors.New("expected function block")
	}

	a.functions = append(a.functions, *a.currentFunction)
	a.currentFunction = nil
	return nil
}

func (a *OpAssembler) Args() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(argsBlock)
	return nil
}

func (a *OpAssembler) EndArgs() error {
	if a.blocks.Pop() != argsBlock {
		return errors.New("expected args block")
	}
	return nil
}

func (a *OpAssembler) Rets() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(retsBlock)
	return nil
}

func (a *OpAssembler) EndRets() error {
	if a.blocks.Pop() != retsBlock {
		return errors.New("expected rets block")
	}

	return nil
}

func (a *OpAssembler) Locals() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(localsBlock)
	return nil
}

func (a *OpAssembler) EndLocals() error {
	// TODO: Validate there's a current ongoing function.

	if a.blocks.Pop() != localsBlock {
		return errors.New("expected locals block")
	}

	if err := a.StackAlloc(a.currentFunction.LocalsBytes()); err != nil {
		return err
	}

	return nil
}

func (a *OpAssembler) DefineArg(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.args[id] = OpVar{a.currentFunction.ArgsBytes(), size}
	return nil
}

func (a *OpAssembler) DefineRet(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.rets[id] = OpVar{a.currentFunction.RetsBytes(), size}
	return nil
}

func (a *OpAssembler) DefineLocal(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	a.currentFunction.locals[id] = OpVar{a.currentFunction.LocalsBytes(), size}
	return nil
}

func (a *OpAssembler) DefineVar(id string, size uint16) error {
	// TODO: Validate there's a current ongoing function.

	switch block := a.blocks.Peek(); block {
	case argsBlock:
		return a.DefineArg(id, size)
	case retsBlock:
		return a.DefineRet(id, size)
	case localsBlock:
		return a.DefineLocal(id, size)
	default:
		return fmt.Errorf("Cannot declare id inside block type %d", block)
	}
}

func (a *OpAssembler) If(which bool) error {
	// TODO: Validate there's a current ongoing function.

	a.assemblers.Push(NewAssembler())
	if which {
		a.blocks.Push(ifTrueBlock)
	} else {
		a.blocks.Push(ifFalseBlock)
	}
	return nil
}

func (a *OpAssembler) EndIf() error {
	block := a.blocks.Pop()
	if block != ifTrueBlock && block != ifFalseBlock {
		return errors.New("expected if block")
	}

	nested := a.assemblers.Pop()

	if block == ifTrueBlock {
		a.asm().PutOpCode(vm.IfTrue)
	} else {
		a.asm().PutOpCode(vm.IfFalse)
	}

	a.asm().
		PutI64(uint64(len(nested.Data()))).
		append(nested.Data())
	return nil
}

func (a *OpAssembler) End() error {
	switch block := a.blocks.Peek(); block {
	case functionBlock:
		return a.EndFunction()
	case argsBlock:
		return a.EndArgs()
	case retsBlock:
		return a.EndRets()
	case localsBlock:
		return a.EndLocals()
	case ifTrueBlock, ifFalseBlock:
		return a.EndIf()
	default:
		return fmt.Errorf("Unknown block type %d", block)
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

func (a *OpAssembler) PrintI8() error {
	a.asm().PutOpCode(vm.PrintI8)
	return nil
}

func (a *OpAssembler) PrintI16() error {
	a.asm().PutOpCode(vm.PrintI16)
	return nil
}

func (a *OpAssembler) PrintI32() error {
	a.asm().PutOpCode(vm.PrintI32)
	return nil
}

func (a *OpAssembler) PrintI64() error {
	a.asm().PutOpCode(vm.PrintI64)
	return nil
}

func (a *OpAssembler) AddI8() error {
	a.asm().PutOpCode(vm.AddI8)
	return nil
}

func (a *OpAssembler) AddI16() error {
	a.asm().PutOpCode(vm.AddI16)
	return nil
}

func (a *OpAssembler) AddI32() error {
	a.asm().PutOpCode(vm.AddI32)
	return nil
}

func (a *OpAssembler) AddI64() error {
	a.asm().PutOpCode(vm.AddI64)
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
		stack.New[*Assembler](), /* assemblers */
		stack.New[blockType](),  /* blocks */
		[]OpFunction{},
		nil, /* currentFunction */
	}
	assembler.assemblers.Push(NewAssembler())
	return assembler
}
