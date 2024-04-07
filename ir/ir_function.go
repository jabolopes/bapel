package ir

type IrFunction struct {
	ID       string
	TypeVars []string
	Args     []IrDecl
	Rets     []IrDecl
}

func NewFunction(id string, typeVars []string, args, rets []IrDecl) IrFunction {
	return IrFunction{id, typeVars, args, rets}
}
