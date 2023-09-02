package bplparser

import (
	"os"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func TestParseCallAssign(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"a b", ir.NewTupleTerm([]ir.IrTerm{newID("a"), newID("b")})},
		{"Index.get a 1", ir.NewIndexGetTerm(newID("a"), newNumber(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(newID("a"), newNumber(1), newNumber(10))},
		{"- a", ir.NewOpUnaryTerm("-", newID("a"))},
		{"a + b", ir.NewOpBinaryTerm("+", newID("a"), newID("b"))},
		{"widen a", ir.NewWidenTerm(newID("a"))},
		{"r <- a", ir.NewAssignTerm(newID("a"), newID("r"))},
		{"r1 r2 <- a1 a2", ir.NewAssignTerm(ir.NewTupleTerm([]ir.IrTerm{newID("a1"), newID("a2")}), ir.NewTupleTerm([]ir.IrTerm{newID("r1"), newID("r2")}))},
	}

	parser := NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		parser.SetLine(test.input)
		if got, err := parser.parseCallAssign(); !reflect.DeepEqual(got, test.want) || err != nil {
			t.Errorf("ParseCallAssign(%q) = %v, %v; want %v, %v",
				test.input, got, err, test.want, nil)
		}
	}
}
