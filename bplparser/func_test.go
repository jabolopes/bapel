package bplparser

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func newFunction(tvars []string, args, rets []ir.IrDecl, body ir.IrTerm) Source {
	return NewFunctionSource(ir.NewFunction("f", tvars, args, rets, body))
}

func TestParseFunc(t *testing.T) {
	body := ir.NewBlockTerm(nil)

	tests := []struct {
		input string
		want  Source
	}{
		{"func f() -> () {\n}", newFunction(nil, nil, nil, body)},
		{"func f(a i32) -> () {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{ir.NewTermDecl("a", ir.NewNameType("i32"))},
				nil,
				body),
		},
		{"func f() -> (r i64) {\n}",
			newFunction(
				nil,
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("r", ir.NewNameType("i64")),
				},
				body),
		},
		{"func f(a [i32], b i64) -> () {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(ir.NewNameType("i32"), math.MaxInt)),
					ir.NewTermDecl("b", ir.NewNameType("i64")),
				},
				nil,
				body),
		},
		{"func f(a [i32], b i64) -> (r1 i32, r2 [i64]) {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(ir.NewNameType("i32"), math.MaxInt)),
					ir.NewTermDecl("b", ir.NewNameType("i64")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r1", ir.NewNameType("i32")),
					ir.NewTermDecl("r2", ir.NewArrayType(ir.NewNameType("i64"), math.MaxInt)),
				},
				body),
		},
		{"func f['a](x 'a) -> (r 'a) {\n}",
			newFunction(
				[]string{"a"},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.NewVarType("a")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r", ir.NewVarType("a")),
				},
				body),
		},
		{"func f['a, 'b](x 'a, y 'b) -> (r1 'a, r2 'b) {\n}",
			newFunction(
				[]string{"a", "b"},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.NewVarType("a")),
					ir.NewTermDecl("y", ir.NewVarType("b")),
				},
				[]ir.IrDecl{
					ir.NewTermDecl("r1", ir.NewVarType("a")),
					ir.NewTermDecl("r2", ir.NewVarType("b")),
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
