package bplparser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
)

func TestParseCallAssign(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"f0", newID("f0")},
		{"f0 ()", ir.NewCallTerm("f0", nil /* types */, ir.NewTupleTerm(nil))},
		{"f1 a", ir.NewCallTerm("f1", nil /* types */, newID("a"))},
		{"Index.get a 1", ir.NewIndexGetTerm(newID("a"), newNumber(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(newID("a"), newNumber(1), newNumber(10))},
		{"- a", ir.NewCallTerm("-", nil /* types */, ir.NewTupleTerm([]ir.IrTerm{newNumber(0), newID("a")}))},
		{"a + b", ir.NewCallTerm("+", nil /* types */, ir.NewTupleTerm([]ir.IrTerm{newID("a"), newID("b")}))},
		{"widen a", ir.NewWidenTerm(newID("a"))},
		{"r1 <- f1 a", ir.NewAssignTerm(ir.NewCallTerm("f1", nil /* types */, ir.NewTupleTerm([]ir.IrTerm{newID("a")})), newID("r1"))},
		{"r <- a", ir.NewAssignTerm(newID("a"), newID("r"))},
		{"r1 r2 <- (a1, a2)", ir.NewAssignTerm(ir.NewTupleTerm([]ir.IrTerm{newID("a1"), newID("a2")}), ir.NewTupleTerm([]ir.IrTerm{newID("r1"), newID("r2")}))},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.SetLine(test.input)
		if got, err := parser.parseCallAssign(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("ParseCallAssign(%q) = %v, %v; want %v, %v",
				test.input, got, err, test.want, nil)
		}
	}
}
