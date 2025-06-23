package ir

import (
	"fmt"
	"strings"
)

type IrFunction struct {
	Export   bool
	ID       string
	TypeVars []VarKind
	Args     []IrDecl
	RetType  IrType
	Body     IrTerm

	// Position in source file.
	Pos Pos
}

func (f IrFunction) String() string {
	var b strings.Builder
	if f.Export {
		b.WriteString("export ")
	}
	b.WriteString(fmt.Sprintf("fn %s", f.ID))
	if len(f.TypeVars) > 0 {
		b.WriteString("[")
		b.WriteString("'")
		b.WriteString(f.TypeVars[0].Var)
		b.WriteString(" ")
		b.WriteString(f.TypeVars[0].Kind.String())
		for _, tvar := range f.TypeVars[1:] {
			b.WriteString(", ")
			b.WriteString("'")
			b.WriteString(tvar.Var)
			b.WriteString(" ")
			b.WriteString(tvar.Kind.String())
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
	b.WriteString(") -> ")
	b.WriteString(f.RetType.String())
	b.WriteString(" ")
	b.WriteString(f.Body.String())
	return b.String()
}

func (f IrFunction) Decl() IrDecl {
	argTypes := make([]IrType, len(f.Args))
	for i := range f.Args {
		argTypes[i] = f.Args[i].Term.Type
	}

	typ := ForallVars(f.TypeVars, NewFunctionType(NewTupleType(argTypes), f.RetType))
	decl := NewTermDecl(f.ID, typ, f.Export)
	decl.Pos = f.Pos
	return decl
}

func NewFunction(export bool, id string, typeVars []VarKind, args []IrDecl, retType IrType, body IrTerm) IrFunction {
	return IrFunction{export, id, typeVars, args, retType, body, Pos{}}
}
