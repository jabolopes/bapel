package bplparser

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func newID(id string) ir.IrTerm {
	return ir.NewTokenTerm(parser.NewIDToken(id))
}

func newNumber(value int64) ir.IrTerm {
	return ir.NewTokenTerm(parser.NewNumberToken(value))
}

func TestParseCall(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"f0", newID("f0")},
		{"f0 ()", ir.NewCallTerm("f0", nil /* types */, ir.NewTupleTerm(nil))},
		{"f1 a", ir.NewCallTerm("f1", nil /* types */, ir.NewTupleTerm([]ir.IrTerm{newID("a")}))},
		{"a b", ir.NewCallTerm("a", nil /* types */, newID("b"))},
		{"Index.get a 1", ir.NewIndexGetTerm(newID("a"), newNumber(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(newID("a"), newNumber(1), newNumber(10))},
		{"- a", ir.NewCallTerm("-", nil /* types */, ir.NewTupleTerm([]ir.IrTerm{newNumber(0), newID("a")}))},
		{"a + b", ir.NewCallTerm("+", nil /* types */, ir.NewTupleTerm([]ir.IrTerm{newID("a"), newID("b")}))},
		{"widen a", ir.NewWidenTerm(newID("a"))},
	}

	compiler := ir.NewCompiler(os.Stdout)
	compiler.Section("imports", []ir.IrDecl{
		ir.NewTermDecl("f0", ir.NewFunctionType(ir.NewTupleType(nil), ir.NewTupleType(nil))),
		ir.NewTermDecl("f1", ir.NewFunctionType(ir.NewNameType("int"), ir.NewTupleType(nil))),
	})

	parser := NewParser(compiler)
	for _, test := range tests {
		parser.SetLine(test.input)
		if got, err := parser.parseCall(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("parseCall(%q) = %v, %v; want %v, %v",
				test.input, got, err, test.want, nil)
		}
	}
}
