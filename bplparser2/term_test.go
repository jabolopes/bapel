package bplparser2_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func newThen(term ir.IrTerm) ir.IrTerm {
	return ir.NewBlockTerm([]ir.IrTerm{term})
}

func newElse(term ir.IrTerm) *ir.IrTerm {
	x := ir.NewBlockTerm([]ir.IrTerm{term})
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
		{"x <- 1 ;", ir.NewAssignTerm(one, x)},

		// If.
		{`if x {
0
} else {
1
}`, ir.NewIfTerm(false /* negated */, nil /* types */, x, newThen(zero), newElse(one))},

		{`if !x {
0
} else {
1
}`, ir.NewIfTerm(false /* negated */, nil /* types */, ir.Call("!", x), newThen(zero), newElse(one))},

		{`if [i8] x {
0
} else {
1
}`, ir.NewIfTerm(false /* negated */, []ir.IrType{ir.NewNameType("i8")}, x, newThen(zero), newElse(one))},

		// Tuple.
		{"() ;", tupleTerm0},
		{"(x, x) ;", tupleTerm2},

		// ID.
		{"x ;", x},
	}

	parser := bplparser2.NewParser()
	parser.SetInitialSymbol("Term")
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))

		got, err := bplparser2.Parse[ir.IrTerm](parser)
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
