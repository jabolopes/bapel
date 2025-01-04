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

func newFunction(tvars []ir.VarKind, args []ir.IrDecl, retType ir.IrType, body ir.IrTerm) bplparser.Source {
	return bplparser.NewFunctionSource(ir.NewFunction(false /* export */, "f", tvars, args, retType, body))
}

func TestParseFunction(t *testing.T) {
	body := ir.NewBlockTerm(nil)
	i32 := ir.NewNameType("i32")
	i64 := ir.NewNameType("i64")
	unit := ir.NewTupleType(nil)

	tests := []struct {
		input string
		want  bplparser.Source
	}{
		{"fn f() -> () {\n}", newFunction(nil, nil, unit, body)},
		{"fn f(a: i32) -> () {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{ir.NewTermDecl("a", i32)},
				unit,
				body),
		},
		{"fn f() -> i64 {\n}", newFunction(nil, nil, i64, body)},
		{"fn f(a: [i32], b: i64) -> () {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(i32, math.MaxInt)),
					ir.NewTermDecl("b", i64),
				},
				unit,
				body),
		},
		{"fn f(a: [i32], b: i64) -> (i32, [i64]) {\n}",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(i32, math.MaxInt)),
					ir.NewTermDecl("b", i64),
				},
				ir.Types(i32, ir.NewArrayType(i64, math.MaxInt)),
				body),
		},
		{"fn f['a](x: 'a) -> 'a {\n}",
			newFunction(
				[]ir.VarKind{{"a", ir.NewTypeKind()}},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.Tvar("a")),
				},
				ir.Tvar("a"),
				body),
		},
		{"fn f['a, 'b](x: 'a, y: 'b) -> ('a, 'b) {\n}",
			newFunction(
				[]ir.VarKind{{"a", ir.NewTypeKind()}, {"b", ir.NewTypeKind()}},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.Tvar("a")),
					ir.NewTermDecl("y", ir.Tvar("b")),
				},
				ir.Types(ir.Tvar("a"), ir.Tvar("b")),
				body),
		},
	}

	parser := bplparser2.NewParser()
	parser.SetInitialSymbol("Function")
	for _, test := range tests {
		parser.Open("testfile", strings.NewReader(test.input))
		got, err := bplparser2.Parse[bplparser.Source](parser)

		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty(), cmpopts.IgnoreFields(bplparser.Source{}, "Pos"), cmpopts.IgnoreFields(ir.IrFunction{}, "Pos"), cmpopts.IgnoreFields(ir.IrDecl{}, "Pos"), cmpopts.IgnoreFields(ir.IrTerm{}, "Pos"), cmpopts.IgnoreFields(ir.IrType{}, "Pos")) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
