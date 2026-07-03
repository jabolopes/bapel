package ir

import (
	"fmt"
)

type FunctionArg struct {
	ID   string
	Type IrType
}

func (t FunctionArg) String() string {
	return fmt.Sprintf("%s: %s", t.ID, t.Type)
}

type IrFunction struct {
	Export     bool
	ID         string
	TypeParams []TypeParam
	Args       []FunctionArg
	RetType    IrType
	Body       IrTerm

	// Position in source file.
	Pos Pos
}

func (t IrFunction) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	if t.Export {
		fmt.Fprint(f, "export ")
	}

	fmt.Fprintf(f, "fn %s", t.ID)

	if len(t.TypeParams) > 0 {
		fmt.Fprint(f, "[")
		Interleave(t.TypeParams, func() { fmt.Fprint(f, ", ") }, func(_ int, tp TypeParam) {
			fmt.Fprintf(f, "'%s %s", tp.Var, tp.Kind)
		})
		fmt.Fprint(f, "]")
	}

	fmt.Fprint(f, "(")
	Interleave(t.Args, func() { fmt.Fprint(f, ", ") }, func(_ int, arg FunctionArg) {
		fmt.Fprint(f, arg.String())
	})
	fmt.Fprintf(f, ") -> %s %s", t.RetType, t.Body)
}

func (t IrFunction) Decl() IrDecl {
	argTypes := make([]IrType, len(t.Args))
	for i := range t.Args {
		argTypes[i] = t.Args[i].Type
	}

	typ := ForallVars(t.TypeParams, NewFunctionType(NewTupleType(argTypes), t.RetType))
	decl := NewTermDecl(t.ID, typ, t.Export)
	decl.Pos = t.Pos
	return decl
}

func NewFunction(export bool, id string, typeParams []TypeParam, args []FunctionArg, retType IrType, body IrTerm) IrFunction {
	return IrFunction{export, id, typeParams, args, retType, body, Pos{}}
}
