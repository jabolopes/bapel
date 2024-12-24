package bplparser2_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"x", ir.ID("x")},
		{"f0 ()", ir.CallID("f0")},
		{"f1 x", ir.CallID("f1", ir.ID("x"))},
		{"f2 (x, y)", ir.CallID("f2", ir.ID("x"), ir.ID("y"))},
		{"Index.get a 1", ir.NewIndexGetTerm(ir.ID("a"), ir.Number(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(ir.ID("a"), ir.Number(1), ir.Number(10))},
		{"- a", ir.CallID("-", ir.Number(0), ir.ID("a"))},
		{"a + b", ir.CallID("+", ir.ID("a"), ir.ID("b"))},
		{"! a", ir.CallID("!", ir.ID("a"))},
		{"1 [i8]", ir.CallPF(ir.Number(1), []ir.IrType{ir.Const("i8")})},
	}

	parser := bplparser2.NewParser()
	parser.SetInitialSymbol("Expression")
	for _, test := range tests {
		parser.Open("testfile", strings.NewReader(test.input))

		got, err := bplparser2.Parse[ir.IrTerm](parser)
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty(), cmpopts.IgnoreFields(ir.IrTerm{}, "Pos"), cmpopts.IgnoreFields(ir.IrType{}, "Pos")) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
