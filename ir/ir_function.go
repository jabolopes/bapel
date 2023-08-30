package ir

import (
	"fmt"
)

// irFunction is a function in IR.
type irFunction struct {
	id     string
	args   []IrDecl
	rets   []IrDecl
	locals []IrDecl
}

func (f *irFunction) lookupVar(id string) (IrDecl, error) {
	for _, decl := range f.locals {
		if decl.ID == id {
			return decl, nil
		}
	}

	for _, decl := range f.rets {
		if decl.ID == id {
			return decl, nil
		}
	}

	for _, decl := range f.args {
		if decl.ID == id {
			return decl, nil
		}
	}

	return IrDecl{}, fmt.Errorf("undefined variable %q", id)
}

func (f *irFunction) addLocal(decl IrDecl) error {
	if _, err := f.lookupVar(decl.ID); err == nil {
		return fmt.Errorf("variable %q already defined in this function", decl.ID)
	}

	f.locals = append(f.locals, decl)
	return nil
}

func (f *irFunction) decl() IrDecl {
	args := make([]IrType, len(f.args))
	for i := range f.args {
		args[i] = f.args[i].Type
	}

	rets := make([]IrType, len(f.rets))
	for i := range f.rets {
		rets[i] = f.rets[i].Type
	}

	return NewTermDecl(f.id, NewFunctionType(IrFunctionType{args, rets}))
}

func NewFunction(id string, args, rets []IrDecl) irFunction {
	return irFunction{id, args, rets, nil}
}
