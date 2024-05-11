package bplparser

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func newFunction(tvars []ir.VarKind, args, rets []ir.IrDecl, body ir.IrTerm) Source {
	return NewFunctionSource(ir.NewFunction(false /* export */, "f", tvars, args, rets, body))
}

func TestParseFunc(t *testing.T) {
	body := ir.NewBlockTerm(nil)
	i32 := ir.NewNameType("i32")
	i64 := ir.NewNameType("i64")

	tests := []struct {
		input string
		want  Source
	}{
		{"func f() -> () {\n}", newFunction(nil, nil, nil, body)},
		{"func f(a i32) -> () {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{ir.NewTermDecl("a", i32)},
				nil,
				body),
		},
		{"func f() -> (r i64) {\n}",
			newFunction(
				nil,
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("r", i64),
				},
				body),
		},
		{"func f(a [i32], b i64) -> () {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(i32, math.MaxInt)),
					ir.NewTermDecl("b", i64),
				},
				nil,
				body),
		},
		{"func f(a [i32], b i64) -> (r1 i32, r2 [i64]) {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(i32, math.MaxInt)),
					ir.NewTermDecl("b", i64),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r1", i32),
					ir.NewTermDecl("r2", ir.NewArrayType(i64, math.MaxInt)),
				},
				body),
		},
		{"func f['a](x 'a) -> (r 'a) {\n}",
			newFunction(
				[]ir.VarKind{{"a", ir.NewTypeKind()}},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.Tvar("a")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r", ir.Tvar("a")),
				},
				body),
		},
		{"func f['a, 'b](x 'a, y 'b) -> (r1 'a, r2 'b) {\n}",
			newFunction(
				[]ir.VarKind{{"a", ir.NewTypeKind()}, {"b", ir.NewTypeKind()}},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.Tvar("a")),
					ir.NewTermDecl("y", ir.Tvar("b")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r1", ir.Tvar("a")),
					ir.NewTermDecl("r2", ir.Tvar("b")),
				},
				body),
		},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))
		parser.Scan()
		if got, err := parser.parseFunc(); !reflect.DeepEqual(got, test.want) || err != nil {
			t.Errorf("parseFunc() = %v, %v; want %v, %v", got, err, test.want, nil)
		}
	}
}
