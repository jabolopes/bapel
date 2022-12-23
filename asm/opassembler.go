package asm

import (
	"fmt"

	"github.com/jabolopes/go-vm/vm"
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
	assembler       *Assembler
	blocks          []blockType
	functions       []OpFunction
	currentFunction *OpFunction
}

func (a *OpAssembler) StackAlloc(size uint16) error {
	a.assembler.PutOpCode(vm.StackAlloc)
	a.assembler.PutI16(size)
	return nil
}

func (a *OpAssembler) Function(id string) error {
	// TODO: Validate there's no current ongoing function.

	a.blocks = append(a.blocks, functionBlock)

	a.currentFunction = &OpFunction{
		id,
		uint64(a.assembler.Len()),
		map[string]OpVar{}, /* args */
		map[string]OpVar{}, /* rets */
		map[string]OpVar{}, /* locals */
	}

	return nil
}

func (a *OpAssembler) EndFunction() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks = a.blocks[:len(a.blocks)-1]
	a.functions = append(a.functions, *a.currentFunction)
	a.currentFunction = nil
	return nil
}

func (a *OpAssembler) Args() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks = append(a.blocks, argsBlock)
	return nil
}

func (a *OpAssembler) EndArgs() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks = a.blocks[:len(a.blocks)-1]
	return nil
}

func (a *OpAssembler) Rets() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks = append(a.blocks, retsBlock)
	return nil
}

func (a *OpAssembler) EndRets() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks = a.blocks[:len(a.blocks)-1]
	return nil
}

func (a *OpAssembler) Locals() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks = append(a.blocks, localsBlock)
	return nil
}

func (a *OpAssembler) EndLocals() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks = a.blocks[:len(a.blocks)-1]

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

	switch block := a.blocks[len(a.blocks)-1]; block {
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

func (a *OpAssembler) End() error {
	switch block := a.blocks[len(a.blocks)-1]; block {
	case functionBlock:
		return a.EndFunction()
	case argsBlock:
		return a.EndArgs()
	case retsBlock:
		return a.EndRets()
	case localsBlock:
		return a.EndLocals()
	default:
		return fmt.Errorf("Unknown block type %d", block)
	}
}

func (a *OpAssembler) PushI8(value byte) error {
	a.assembler.PutOpCode(vm.PushI8)
	a.assembler.PutI8(value)
	return nil
}

func (a *OpAssembler) PushI16(value uint16) error {
	a.assembler.PutOpCode(vm.PushI16)
	a.assembler.PutI16(value)
	return nil
}

func (a *OpAssembler) PushI32(value uint32) error {
	a.assembler.PutOpCode(vm.PushI32)
	a.assembler.PutI32(value)
	return nil
}

func (a *OpAssembler) PushI64(value uint64) error {
	a.assembler.PutOpCode(vm.PushI64)
	a.assembler.PutI64(value)
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
		a.assembler.PutOpCode(vm.PushLocalI8)
	case 2:
		a.assembler.PutOpCode(vm.PushLocalI16)
	case 4:
		a.assembler.PutOpCode(vm.PushLocalI32)
	case 8:
		a.assembler.PutOpCode(vm.PushLocalI64)
	}

	a.assembler.PutI16(local.offset)
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
		a.assembler.PutOpCode(vm.PopLocalI8)
	case 2:
		a.assembler.PutOpCode(vm.PopLocalI16)
	case 4:
		a.assembler.PutOpCode(vm.PopLocalI32)
	case 8:
		a.assembler.PutOpCode(vm.PopLocalI64)
	}

	a.assembler.PutI16(local.offset)
	return nil
}

func (a *OpAssembler) PrintI8() error {
	a.assembler.PutOpCode(vm.PrintI8)
	return nil
}

func (a *OpAssembler) PrintI16() error {
	a.assembler.PutOpCode(vm.PrintI16)
	return nil
}

func (a *OpAssembler) PrintI32() error {
	a.assembler.PutOpCode(vm.PrintI32)
	return nil
}

func (a *OpAssembler) PrintI64() error {
	a.assembler.PutOpCode(vm.PrintI64)
	return nil
}

func (a *OpAssembler) AddI8() error {
	a.assembler.PutOpCode(vm.AddI8)
	return nil
}

func (a *OpAssembler) AddI16() error {
	a.assembler.PutOpCode(vm.AddI16)
	return nil
}

func (a *OpAssembler) AddI32() error {
	a.assembler.PutOpCode(vm.AddI32)
	return nil
}

func (a *OpAssembler) AddI64() error {
	a.assembler.PutOpCode(vm.AddI64)
	return nil
}

func (a *OpAssembler) Program() vm.OpProgram {
	return vm.OpProgram{
		a.assembler.Data(),
		[]vm.OpFunction{},
	}
}

func New() *OpAssembler {
	return &OpAssembler{
		NewAssembler(),
		nil, /* blocks */
		[]OpFunction{},
		nil, /* currentFunction */
	}
}
