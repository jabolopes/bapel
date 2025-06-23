package bplparser2_test

import (
	"math"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func newFunction(tvars []ir.VarKind, args []ir.IrDecl, retType ir.IrType, body ir.IrTerm) ast.Source {
	return ast.NewFunctionSource(ir.NewFunction(false /* export */, "f", tvars, args, retType, body))
}

func TestParseFunction(t *testing.T) {
	body := ir.NewBlockTerm([]ir.IrTerm{ir.Terms()})
	i32 := ir.NewNameType("i32")
	i64 := ir.NewNameType("i64")
	unit := ir.NewTupleType(nil)

	tests := []struct {
		input string
		want  ast.Source
	}{
		{"fn f() -> () { (); }", newFunction(nil, nil, unit, body)},
		{"fn f(a: i32) -> () { (); }",
			newFunction(
				nil,
				[]ir.IrDecl{ir.NewTermDecl("a", i32, false /* export */)},
				unit,
				body),
		},
		{"fn f() -> i64 { (); }", newFunction(nil, nil, i64, body)},
		{"fn f(a: [i32], b: i64) -> () { (); }",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(i32, math.MaxInt), false /* export */),
					ir.NewTermDecl("b", i64, false /* export */),
				},
				unit,
				body),
		},
		{"fn f(a: [i32], b: i64) -> (i32, [i64]) { (); }",
			newFunction(
				nil,
				[]ir.IrDecl{
					ir.NewTermDecl("a", ir.NewArrayType(i32, math.MaxInt), false /* export */),
					ir.NewTermDecl("b", i64, false /* export */),
				},
				ir.Types(i32, ir.NewArrayType(i64, math.MaxInt)),
				body),
		},
		{"fn f['a](x: 'a) -> 'a { (); }",
			newFunction(
				[]ir.VarKind{{"a", ir.NewTypeKind()}},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.Tvar("a"), false /* export */),
				},
				ir.Tvar("a"),
				body),
		},
		{"fn f['a, 'b](x: 'a, y: 'b) -> ('a, 'b) { (); }",
			newFunction(
				[]ir.VarKind{{"a", ir.NewTypeKind()}, {"b", ir.NewTypeKind()}},
				[]ir.IrDecl{
					ir.NewTermDecl("x", ir.Tvar("a"), false /* export */),
					ir.NewTermDecl("y", ir.Tvar("b"), false /* export */),
				},
				ir.Types(ir.Tvar("a"), ir.Tvar("b")),
				body),
		},
	}

	parser, err := bplparser2.NewWithSymbol("Function")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		parser.Open("testfile", strings.NewReader(test.input))
		got, err := bplparser2.Parse[ast.Source](parser)

		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty(), cmpopts.IgnoreFields(ast.Source{}, "Pos"), cmpopts.IgnoreFields(ir.IrFunction{}, "Pos"), cmpopts.IgnoreFields(ir.IrDecl{}, "Pos"), cmpopts.IgnoreFields(ir.IrTerm{}, "Pos"), cmpopts.IgnoreFields(ir.IrType{}, "Pos")) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
