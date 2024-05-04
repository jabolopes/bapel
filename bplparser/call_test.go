package bplparser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
)

func TestParseCall(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"x", ir.ID("x")},
		{"f0 ()", ir.Call("f0")},
		{"f1 x", ir.Call("f1", ir.ID("x"))},
		{"f2 (x, y)", ir.Call("f2", ir.ID("x"), ir.ID("y"))},
		{"Index.get a 1", ir.NewIndexGetTerm(ir.ID("a"), ir.Number(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(ir.ID("a"), ir.Number(1), ir.Number(10))},
		{"- a", ir.Call("-", ir.Number(0), ir.ID("a"))},
		{"a + b", ir.Call("+", ir.ID("a"), ir.ID("b"))},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))
		parser.Scan()
		if got, err := parser.parseCall(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("parseCall(%q) = %v, %v; want %v, %v",
				test.input, got, err, test.want, nil)
		}
	}
}
