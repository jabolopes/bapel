package ir

import (
	"fmt"
)

// irFunction is a function in IR.
type irFunction struct {
	id   string  // Name of function.
	vars []IrVar // Variables in the order in which they were defined.
}

func (f *irFunction) lookupVar(id string) (IrDecl, error) {
	for _, irvar := range f.vars {
		if irvar.Id == id {
			return irvar.decl(), nil
		}
	}

	return IrDecl{}, fmt.Errorf("undefined variable %q", id)
}

func (f *irFunction) addVar(id string, irvar IrVar) error {
	if _, err := f.lookupVar(id); err == nil {
		return fmt.Errorf("Variable %q already defined in this context", id)
	}

	f.vars = append(f.vars, irvar)
	return nil
}

func (f *irFunction) decl() IrDecl {
	var args []IrType
	var rets []IrType
	for _, irvar := range f.vars {
		if irvar.VarType == ArgVar {
			args = append(args, irvar.Type)
		} else if irvar.VarType == RetVar {
			rets = append(rets, irvar.Type)
		}
	}

	return NewConstantDecl(f.id, NewFunctionType(IrFunctionType{args, rets}))
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
