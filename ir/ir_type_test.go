package ir

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func tvars(xs ...VarKind) []VarKind {
	return append([]VarKind{}, xs...)
}

func TestNewForallType(t *testing.T) {
	tests := []struct {
		tvars []VarKind
		typ   IrType
		want  IrType
	}{
		{nil, NewNameType("i8"), NewNameType("i8")},
		{
			tvars(VarKind{"a", NewTypeKind()}), NewVarType("a"),
			Forall("a", NewTypeKind(), NewVarType("a")),
		},
		{
			tvars(VarKind{"a", NewTypeKind()}, VarKind{"b", NewTypeKind()}), NewFunctionType(NewVarType("a"), NewVarType("b")),
			Forall("a", NewTypeKind(), Forall("b", NewTypeKind(), NewFunctionType(NewVarType("a"), NewVarType("b")))),
		},
		{
			tvars(VarKind{"a", NewTypeKind()}, VarKind{"a", NewTypeKind()}), NewFunctionType(NewVarType("a"), NewVarType("a")),
			Forall("a", NewTypeKind(), Forall("a", NewTypeKind(), NewFunctionType(NewVarType("a"), NewVarType("a")))),
		},
		{
			tvars(VarKind{"a", NewTypeKind()}), NewForallType("b", NewTypeKind(), NewFunctionType(NewVarType("a"), NewVarType("b"))),
			Forall("a", NewTypeKind(), Forall("b", NewTypeKind(), NewFunctionType(NewVarType("a"), NewVarType("b")))),
		},
	}

	for _, test := range tests {
		if got := ForallVars(test.tvars, test.typ); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("ForallVars(%v, %v) = %v; want %v", test.tvars, test.typ, got, test.want)
		}
	}
}

func TestForallVars(t *testing.T) {
	tests := []struct {
		got  IrType
		want IrType
	}{
		{
			ForallVars([]VarKind{{"a", NewTypeKind()}, {"b", NewTypeKind()}}, NewFunctionType(Tvar("a"), Tvar("b"))),
			Forall("a", NewTypeKind(), Forall("b", NewTypeKind(), NewFunctionType(Tvar("a"), Tvar("b")))),
		},
	}

	for _, test := range tests {
		if !cmp.Equal(test.got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("got = %v; want %v", test.got, test.want)
		}
	}
}
