package bplparser2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func TestParseAssign(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"r <- 1", ir.NewAssignTerm(ir.Number(1), ir.ID("r"))},
		{"r <- x", ir.NewAssignTerm(ir.ID("x"), ir.ID("r"))},
		{"(r1, r2) <- (a1, a2)", ir.NewAssignTerm(ir.Terms(ir.ID("a1"), ir.ID("a2")), ir.Terms(ir.ID("r1"), ir.ID("r2")))},
		{"r <- f0 ()", ir.NewAssignTerm(ir.CallID("f0"), ir.ID("r"))},
		{"r <- f1 x", ir.NewAssignTerm(ir.CallID("f1", ir.ID("x")), ir.ID("r"))},
		{"r <- f2 (x, y)", ir.NewAssignTerm(ir.CallID("f2", ir.ID("x"), ir.ID("y")), ir.ID("r"))},
		{"r <- x->1", ir.NewAssignTerm(ir.NewProjectionTerm(ir.ID("x"), "1"), ir.ID("r"))},
		{"r <- x->y", ir.NewAssignTerm(ir.NewProjectionTerm(ir.ID("x"), "y"), ir.ID("r"))},
		{"r <- - a", ir.NewAssignTerm(ir.CallID("-", ir.Number(0), ir.ID("a")), ir.ID("r"))},
		{"r <- a + b", ir.NewAssignTerm(ir.CallID("+", ir.ID("a"), ir.ID("b")), ir.ID("r"))},
	}

	parser, err := bplparser2.NewWithSymbol("AssignTerm")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		input := fmt.Sprintf("%s ;", test.input)
		parser.Open("testfile", strings.NewReader(input))

		got, err := bplparser2.Parse[ir.IrTerm](parser)
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty(), cmpopts.IgnoreFields(ir.IrTerm{}, "Pos"), cmpopts.IgnoreFields(ir.IrLiteral{}, "Pos")) || err != nil {
			t.Errorf("ParseAssign(%q) = %v, %v; want %v, %v",
				test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
