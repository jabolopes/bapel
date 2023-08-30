package ir

// irFunction is a function in IR.
type irFunction struct {
	id     string
	args   []IrDecl
	rets   []IrDecl
	locals []IrDecl
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
