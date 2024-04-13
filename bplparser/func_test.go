package bplparser

import (
	"math"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func newFunction(tvars []string, args, rets []ir.IrDecl) Source {
	return NewFunctionSource(ir.NewFunction("f", tvars, args, rets))
}

func TestParseFunc(t *testing.T) {
	tests := []struct {
		input string
		want  Source
	}{
		{"func f() -> () {", newFunction(nil, nil, nil)},
		{"func f(a i32) -> () {",
			newFunction(
				nil,
				[]ir.IrDecl{ir.NewTermDecl("a", ir.NewNameType("i32"))},
				nil),
		},
		{"func f() -> (r i64) {",
			newFunction(
				nil,
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("r", ir.NewNameType("i64")),
				}),
		},
		{"func f(a [i32], b i64) -> () {",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(ir.NewNameType("i32"), math.MaxInt)),
					ir.NewTermDecl("b", ir.NewNameType("i64")),
				},
				nil),
		},
		{"func f(a [i32], b i64) -> (r1 i32, r2 [i64]) {",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(ir.NewNameType("i32"), math.MaxInt)),
					ir.NewTermDecl("b", ir.NewNameType("i64")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r1", ir.NewNameType("i32")),
					ir.NewTermDecl("r2", ir.NewArrayType(ir.NewNameType("i64"), math.MaxInt)),
				}),
		},
		{"func f['a](x 'a) -> (r 'a) {",
			newFunction(
				[]string{"a"},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.NewVarType("a")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r", ir.NewVarType("a")),
				}),
		},
		{"func f['a, 'b](x 'a, y 'b) -> (r1 'a, r2 'b) {",
			newFunction(
				[]string{"a", "b"},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.NewVarType("a")),
					ir.NewTermDecl("y", ir.NewVarType("b")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r1", ir.NewVarType("a")),
					ir.NewTermDecl("r2", ir.NewVarType("b")),
				}),
		},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.SetLine(test.input)
		if got, err := parser.parseFunc(); !reflect.DeepEqual(got, test.want) || err != nil {
			t.Errorf("parseFunc() = %v, %v; want %v, %v", got, err, test.want, nil)
		}
	}
}
