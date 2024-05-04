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
		{"r <- 1", ir.NewAssignTerm(ir.Number(1), ir.ID("r"))},
		{"r <- x", ir.NewAssignTerm(ir.ID("x"), ir.ID("r"))},
		{"r1 r2 <- (a1, a2)", ir.NewAssignTerm(ir.Terms(ir.ID("a1"), ir.ID("a2")), ir.Terms(ir.ID("r1"), ir.ID("r2")))},
		{"r <- f0 ()", ir.NewAssignTerm(ir.Call("f0"), ir.ID("r"))},
		{"r <- f1 x", ir.NewAssignTerm(ir.Call("f1", ir.ID("x")), ir.ID("r"))},
		{"r <- f2 (x, y)", ir.NewAssignTerm(ir.Call("f2", ir.ID("x"), ir.ID("y")), ir.ID("r"))},
		{"r <- Index.get x 1", ir.NewAssignTerm(ir.NewIndexGetTerm(ir.ID("x"), ir.Number(1)), ir.ID("r"))},
		{"r <- Index.set x 1 10", ir.NewAssignTerm(ir.NewIndexSetTerm(ir.ID("x"), ir.Number(1), ir.Number(10)), ir.ID("r"))},
		{"r <- - a", ir.NewAssignTerm(ir.Call("-", ir.Number(0), ir.ID("a")), ir.ID("r"))},
		{"r <- a + b", ir.NewAssignTerm(ir.Call("+", ir.ID("a"), ir.ID("b")), ir.ID("r"))},
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
