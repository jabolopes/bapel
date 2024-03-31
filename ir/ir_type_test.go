package ir

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func vars(xs ...string) []string {
	return append([]string{}, xs...)
}

func forallType(vars []string, typ IrType) IrType {
	if vars == nil {
		panic("forall type cannot have empty type variables")
	}

	return IrType{
		Case: ForallType,
		Forall: &struct {
			Vars []string
			Type IrType
		}{vars, typ},
	}
}

func TestNewForallType(t *testing.T) {
	tests := []struct {
		vars []string
		typ  IrType
		want IrType
	}{
		{nil, NewNameType("i8"), NewNameType("i8")},
		{
			vars("a"), NewVarType("a"),
			forallType(vars("a"), NewVarType("a")),
		},
		{
			vars("a", "b"), NewFunctionType(NewVarType("a"), NewVarType("b")),
			forallType(vars("a", "b"), NewFunctionType(NewVarType("a"), NewVarType("b"))),
		},
		{
			vars("a", "a"), NewFunctionType(NewVarType("a"), NewVarType("a")),
			forallType(vars("a"), NewFunctionType(NewVarType("a"), NewVarType("a"))),
		},
		{
			vars("a"), NewForallType(vars("b"), NewFunctionType(NewVarType("a"), NewVarType("b"))),
			forallType(vars("a", "b"), NewFunctionType(NewVarType("a"), NewVarType("b"))),
		},
	}

	for _, test := range tests {
		if got := NewForallType(test.vars, test.typ); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("NewForallType(%v, %v) = %v; want %v", test.vars, test.typ, got, test.want)
		}
	}
}
