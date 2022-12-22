package asm

import (
	"fmt"

	"github.com/jabolopes/go-vm/vm"
)

type OpLocal struct {
	offset uint16
	size   uint16
}

type OpFunction struct {
	argBytes   uint16
	localBytes uint16
	locals     map[string]OpLocal
}

type OpAssembler struct {
	assembler *Assembler

	functions map[string]OpFunction

	currentFunctionId     string
	currentFunctionOffset uint64
}

func (a *OpAssembler) StackAlloc(size uint16) error {
	a.assembler.PutOpCode(vm.StackAlloc)
	a.assembler.PutI16(size)
	return nil
}

func (a *OpAssembler) Function(id string, argBytes, localBytes uint16) error {
	// TODO: Validate there's no current ongoing function.

	a.currentFunctionId = id
	a.currentFunctionOffset = uint64(a.assembler.Len())

	if err := a.StackAlloc(localBytes); err != nil {
		return err
	}

	return nil
}

func (a *OpAssembler) LocalGet(id string) error {
	local, ok := a.functions[a.currentFunctionId].locals[id]
	if !ok {
		return fmt.Errorf("Undeclared local %q", id)
	}

	switch local.size {
	case 1:
		a.assembler.PutOpCode(vm.PushL8)
	case 2:
		a.assembler.PutOpCode(vm.PushL16)
	case 4:
		a.assembler.PutOpCode(vm.PushL32)
	case 8:
		a.assembler.PutOpCode(vm.PushL64)
	}

	a.assembler.PutI16(local.offset)

	return nil
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

func (a *OpAssembler) Program() vm.OpProgram {
	return vm.OpProgram{
		a.assembler.Data(),
		[]vm.OpFunction{},
	}
}

func New() *OpAssembler {
	return &OpAssembler{
		NewAssembler(),
		map[string]OpFunction{},
		"", /* currentFunctionId */
		0,  /* currentFunctionOffset  */
	}
}
