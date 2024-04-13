package bplparser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func newThen(term ir.IrTerm) ir.IrTerm {
	return ir.NewBlockTerm([]ir.IrTerm{ir.NewStatementTerm(term)})
}

func newElse(term ir.IrTerm) *ir.IrTerm {
	x := ir.NewBlockTerm([]ir.IrTerm{ir.NewStatementTerm(term)})
	return &x
}

func TestParseTerm(t *testing.T) {
	x := ir.NewTokenTerm(parser.NewIDToken("x"))
	zero := ir.NewTokenTerm(parser.NewNumberToken(0))
	one := ir.NewTokenTerm(parser.NewNumberToken(1))
	tupleTerm0 := ir.NewTupleTerm(nil)
	tupleTerm2 := ir.NewTupleTerm([]ir.IrTerm{x, x})

	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		// Assign.
		{"x <- 1", ir.NewStatementTerm(ir.NewAssignTerm(one, x))},
		// If.
		{`if x {
0
} else {
1
}`, ir.NewIfTerm(false /* negated */, x, newThen(zero), newElse(one))},
		// Tuple.
		{"()", ir.NewStatementTerm(tupleTerm0)},
		{"x", ir.NewStatementTerm(x)},
		{"(x, x)", ir.NewStatementTerm(tupleTerm2)},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))
		parser.Scan()
		if got, err := parser.parseTerm(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("parseTerm(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
		}
	}
}
