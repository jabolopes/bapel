package bplparser2_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func newThen(term ir.IrTerm) ir.IrTerm {
	return ir.NewBlockTerm([]ir.IrTerm{term})
}

func newElse(term ir.IrTerm) *ir.IrTerm {
	x := ir.NewBlockTerm([]ir.IrTerm{term})
	return &x
}

func TestParseTerm(t *testing.T) {
	x := ir.ID("x")
	zero := ir.Number(0)
	one := ir.Number(1)
	tupleTerm0 := ir.NewTupleTerm(nil)
	tupleTerm2 := ir.NewTupleTerm([]ir.IrTerm{x, x})
	i8 := ir.Const("i8")

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
}`, ir.NewIfTerm(x, newThen(zero), newElse(one))},

		{`if !x {
0
} else {
1
}`, ir.NewIfTerm(ir.Call(ir.ID("!"), x), newThen(zero), newElse(one))},

		{`if x [i8] {
0
} else {
1
}`, ir.NewIfTerm(ir.CallPF(x, ir.TypesA(i8)), newThen(zero), newElse(one))},

		// Tuple.
		{"() ;", tupleTerm0},
		{"(x, x) ;", tupleTerm2},

		// ID.
		{"x ;", x},
	}

	parser, err := bplparser2.NewWithSymbol("Term")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		parser.Open("testfile", strings.NewReader(test.input))

		got, err := bplparser2.Parse[ir.IrTerm](parser)
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty(), cmpopts.IgnoreFields(ir.IrTerm{}, "Pos"), cmpopts.IgnoreFields(ir.IrType{}, "Pos")) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
