package bplparser

import (
	"strings"
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

func newCall(id string, terms ...ir.IrTerm) ir.IrTerm {
	return ir.NewCallTerm(id, nil /* types */, ir.NewTupleTerm(terms))
}

func TestParseCall(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"x", newID("x")},
		{"f0 ()", newCall("f0")},
		{"f1 x", newCall("f1", newID("x"))},
		{"f2 (x, y)", newCall("f2", newID("x"), newID("y"))},
		{"Index.get a 1", ir.NewIndexGetTerm(newID("a"), newNumber(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(newID("a"), newNumber(1), newNumber(10))},
		{"- a", newCall("-", newNumber(0), newID("a"))},
		{"a + b", newCall("+", newID("a"), newID("b"))},
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
