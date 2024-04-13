package ir

type IrFunction struct {
	ID       string
	TypeVars []string
	Args     []IrDecl
	Rets     []IrDecl
	Body     IrTerm
}

func (f IrFunction) Decl() IrDecl {
	argTypes := make([]IrType, len(f.Args))
	for i := range f.Args {
		argTypes[i] = f.Args[i].Type()
	}

	retTypes := make([]IrType, len(f.Rets))
	for i := range f.Rets {
		retTypes[i] = f.Rets[i].Type()
	}

	typ := NewForallType(f.TypeVars, NewFunctionType(NewTupleType(argTypes), NewTupleType(retTypes)))
	return NewTermDecl(f.ID, typ)
}

func NewFunction(id string, typeVars []string, args, rets []IrDecl, body IrTerm) IrFunction {
	return IrFunction{id, typeVars, args, rets, body}
}
