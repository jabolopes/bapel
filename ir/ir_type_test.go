package ir

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func vars(xs ...string) []string {
	return append([]string{}, xs...)
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
			Forall("a", NewVarType("a")),
		},
		{
			vars("a", "b"), NewFunctionType(NewVarType("a"), NewVarType("b")),
			Forall("a", Forall("b", NewFunctionType(NewVarType("a"), NewVarType("b")))),
		},
		{
			vars("a", "a"), NewFunctionType(NewVarType("a"), NewVarType("a")),
			Forall("a", Forall("a", NewFunctionType(NewVarType("a"), NewVarType("a")))),
		},
		{
			vars("a"), NewForallType("b", NewFunctionType(NewVarType("a"), NewVarType("b"))),
			Forall("a", Forall("b", NewFunctionType(NewVarType("a"), NewVarType("b")))),
		},
	}

	for _, test := range tests {
		if got := NewForallVarsType(test.vars, test.typ); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("NewForallType(%v, %v) = %v; want %v", test.vars, test.typ, got, test.want)
		}
	}
}
