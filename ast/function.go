package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type Function struct {
	Export   bool
	ID       string
	TypeVars []ir.VarKind
	Args     []ir.FunctionArg
	RetType  ir.IrType
	Body     Expr

	// Position in source file.
	Pos ir.Pos
}

func (t Function) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	if t.Export {
		fmt.Fprint(f, "pub ")
	}

	fmt.Fprintf(f, "fn %s", t.ID)

	if len(t.TypeVars) > 0 {
		fmt.Fprint(f, "[")
		ir.Interleave(t.TypeVars, func() { fmt.Fprint(f, ", ") }, func(_ int, varkind ir.VarKind) {
			fmt.Fprintf(f, "'%s %s", varkind.Var, varkind.Kind)
		})
		fmt.Fprint(f, "]")
	}

	fmt.Fprint(f, "(")
	ir.Interleave(t.Args, func() { fmt.Fprint(f, ", ") }, func(_ int, arg ir.FunctionArg) {
		fmt.Fprint(f, arg.String())
	})
	fmt.Fprintf(f, ") -> %s %s", t.RetType, t.Body)
}

func (t Function) Decl() ir.IrDecl {
	argTypes := make([]ir.IrType, len(t.Args))
	for i := range t.Args {
		argTypes[i] = t.Args[i].Type
	}

	typ := ir.ForallVars(t.TypeVars, ir.NewFunctionType(ir.NewTupleType(argTypes), t.RetType))
	decl := ir.NewTermDecl(t.ID, typ, t.Export)
	decl.Pos = t.Pos
	return decl
}

func NewFunction(pos ir.Pos, export bool, id string, typeVars []ir.VarKind, args []ir.FunctionArg, retType ir.IrType, body Expr) Function {
	return Function{export, id, typeVars, args, retType, body, pos}
}
