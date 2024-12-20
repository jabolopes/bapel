package bplparser2_test

import (
	"math"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func newFunction(tvars []ir.VarKind, args, rets []ir.IrDecl, body ir.IrTerm) bplparser.Source {
	return bplparser.NewFunctionSource(ir.NewFunction(false /* export */, "f", tvars, args, rets, body))
}

func TestParseFunction(t *testing.T) {
	body := ir.NewBlockTerm(nil)
	i32 := ir.NewNameType("i32")
	i64 := ir.NewNameType("i64")

	tests := []struct {
		input string
		want  bplparser.Source
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

	parser := bplparser2.NewParser()
	parser.SetInitialSymbol("Function")
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))

		got, err := bplparser2.Parse[bplparser.Source](parser)
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
