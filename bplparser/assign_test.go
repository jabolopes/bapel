package bplparser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
)

func TestParseAssign(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"r <- 1", ir.NewAssignTerm(newNumber(1), newID("r"))},
		{"r <- x", ir.NewAssignTerm(newID("x"), newID("r"))},
		{"r1 r2 <- (a1, a2)", ir.NewAssignTerm(ir.NewTupleTerm([]ir.IrTerm{newID("a1"), newID("a2")}), ir.NewTupleTerm([]ir.IrTerm{newID("r1"), newID("r2")}))},
		{"r <- f0 ()", ir.NewAssignTerm(newCall("f0"), newID("r"))},
		{"r <- f1 x", ir.NewAssignTerm(newCall("f1", newID("x")), newID("r"))},
		{"r <- f2 (x, y)", ir.NewAssignTerm(newCall("f2", newID("x"), newID("y")), newID("r"))},
		{"r <- Index.get x 1", ir.NewAssignTerm(ir.NewIndexGetTerm(newID("x"), newNumber(1)), newID("r"))},
		{"r <- Index.set x 1 10", ir.NewAssignTerm(ir.NewIndexSetTerm(newID("x"), newNumber(1), newNumber(10)), newID("r"))},
		{"r <- - a", ir.NewAssignTerm(newCall("-", newNumber(0), newID("a")), newID("r"))},
		{"r <- a + b", ir.NewAssignTerm(newCall("+", newID("a"), newID("b")), newID("r"))},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))
		parser.Scan()
		if got, err := parser.parseAssign(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("ParseAssign(%q) = %v, %v; want %v, %v",
				test.input, got, err, test.want, nil)
			t.Errorf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
