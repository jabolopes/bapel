package bplparser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
)

func newFunctionType(arg, ret ir.IrType) ir.IrType {
	return ir.NewFunctionType(ir.NewTupleType([]ir.IrType{arg}), ir.NewTupleType([]ir.IrType{ret}))
}

func TestParseType(t *testing.T) {
	a := ir.NewVarType("a")
	i8 := ir.NewNameType("i8")
	i16 := ir.NewNameType("i16")

	structType1 := ir.NewStructType([]ir.StructField{
		ir.StructField{"a", i8},
	})

	structType2 := ir.NewStructType([]ir.StructField{
		ir.StructField{"a", i8},
		ir.StructField{"b", i16},
	})

	tupleType0 := ir.NewTupleType(nil)
	tupleType2 := ir.NewTupleType([]ir.IrType{i8, i16})
	tupleTypeAa := ir.NewTupleType([]ir.IrType{a, a})

	tests := []struct {
		input string
		want  ir.IrType
	}{
		// Typename.
		{"i8", i8},
		{"i16", i16},
		// Struct.
		{"{a i8}", structType1},
		{"{a i8, b i16}", structType2},
		// Tuple.
		{"()", tupleType0},
		{"(i8, i16)", tupleType2},
		// Array.
		{"[i8 10]", ir.NewArrayType(i8, 10)},
		// Function.
		{"i8 -> i16", newFunctionType(i8, i16)},
		{"i8 -> (i8, i16)", newFunctionType(i8, tupleType2)},
		{"(i8, i16) -> i16", newFunctionType(tupleType2, i16)},
		{"(i8, i16) -> (i8, i16)", newFunctionType(tupleType2, tupleType2)},
		// Forall.
		{"forall ['a] 'a -> 'a", ir.Forall("a", newFunctionType(a, a))},
		{"forall ['a] ('a, 'a) -> 'a", ir.Forall("a", newFunctionType(tupleTypeAa, a))},
		{"forall ['a] 'a -> ('a, 'a)", ir.Forall("a", newFunctionType(a, tupleTypeAa))},
		{"forall ['a] ('a, 'a) -> ('a, 'a)", ir.Forall("a", newFunctionType(tupleTypeAa, tupleTypeAa))},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))
		parser.Scan()
		if got, err := parser.parseType(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("parseType(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
		}
	}
}
