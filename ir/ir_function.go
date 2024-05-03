package ir

import (
	"fmt"
	"strings"
)

type IrFunction struct {
	ID       string
	TypeVars []string
	Args     []IrDecl
	Rets     []IrDecl
	Body     IrTerm
}

func (f IrFunction) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("func %s", f.ID))
	if len(f.TypeVars) > 0 {
		b.WriteString("[")
		b.WriteString("'")
		b.WriteString(f.TypeVars[0])
		for _, tvar := range f.TypeVars[1:] {
			b.WriteString(", ")
			b.WriteString("'")
			b.WriteString(tvar)
		}
		b.WriteString("]")
	}
	b.WriteString("(")
	if len(f.Args) > 0 {
		b.WriteString(f.Args[0].String())
		for _, arg := range f.Args[1:] {
			b.WriteString(", ")
			b.WriteString(arg.String())
		}
	}
	b.WriteString(") -> (")
	if len(f.Rets) > 0 {
		b.WriteString(f.Rets[0].String())
		for _, ret := range f.Rets[1:] {
			b.WriteString(", ")
			b.WriteString(ret.String())
		}
	}
	b.WriteString(") ")
	b.WriteString(f.Body.String())
	return b.String()
}

func (f IrFunction) Decl() IrDecl {
	argTypes := make([]IrType, len(f.Args))
	for i := range f.Args {
		argTypes[i] = f.Args[i].Term.Type
	}

	retTypes := make([]IrType, len(f.Rets))
	for i := range f.Rets {
		retTypes[i] = f.Rets[i].Term.Type
	}

	typ := NewForallType(f.TypeVars, NewFunctionType(NewTupleType(argTypes), NewTupleType(retTypes)))
	return NewTermDecl(f.ID, typ)
}

func NewFunction(id string, typeVars []string, args, rets []IrDecl, body IrTerm) IrFunction {
	return IrFunction{id, typeVars, args, rets, body}
}
