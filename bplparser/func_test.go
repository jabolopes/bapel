package bplparser

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func newFunction(args, rets []ir.IrDecl) Source {
	return NewFunctionSource("f", args, rets)
}

func TestParseFunc(t *testing.T) {
	tests := []struct {
		input string
		want  Source
	}{
		{"func f() -> () {", newFunction(nil, nil)},
		{"func f(a i32) -> () {",
			newFunction(
				[]ir.IrDecl{ir.NewTermDecl("a", ir.NewIntType(ir.I32))},
				nil),
		},
		{"func f() -> (r i64) {",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("r", ir.NewIntType(ir.I64)),
				}),
		},
		{"func f(a [i32], b i64) -> () {",
			newFunction(
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(ir.NewIntType(ir.I32), math.MaxInt)),
					ir.NewTermDecl("b", ir.NewIntType(ir.I64)),
				},
				nil),
		},
		{"func f(a [i32], b i64) -> (r1 i32, r2 [i64]) {",
			newFunction(
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(ir.NewIntType(ir.I32), math.MaxInt)),
					ir.NewTermDecl("b", ir.NewIntType(ir.I64)),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r1", ir.NewIntType(ir.I32)),
					ir.NewTermDecl("r2", ir.NewArrayType(ir.NewIntType(ir.I64), math.MaxInt)),
				}),
		},
	}

	parser := NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		parser.SetLine(test.input)
		if got, err := parser.parseFunc(); !reflect.DeepEqual(got, test.want) || err != nil {
			t.Errorf("parseFunc() = %v, %v; want %v, %v", got, err, test.want, nil)
		}
	}
}
