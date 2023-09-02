package bplparser

import (
	"os"
	"reflect"
	"testing"

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
		{"a b", ir.NewTupleTerm([]ir.IrTerm{newID("a"), newID("b")})},
		{"Index.get a 1", ir.NewIndexGetTerm(newID("a"), newNumber(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(newID("a"), newNumber(1), newNumber(10))},
		{"- a", ir.NewOpUnaryTerm("-", newID("a"))},
		{"a + b", ir.NewOpBinaryTerm("+", newID("a"), newID("b"))},
		{"widen a", ir.NewWidenTerm(newID("a"))},
	}

	parser := NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		parser.SetLine(test.input)
		if got, err := parser.parseCall(); !reflect.DeepEqual(got, test.want) || err != nil {
			t.Errorf("parseCall(%q) = %v, %v; want %v, %v",
				test.input, got, err, test.want, nil)
		}
	}
}
