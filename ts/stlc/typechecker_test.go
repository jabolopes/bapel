package stlc_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/stlc"
)

func newContext(t *testing.T) stlc.Context {
	t.Helper()

	context := stlc.NewContext()
	{
		binds := []stlc.Bind{
			stlc.NewConstBind("i8", ir.NewTypeKind(), stlc.ImportSymbol),
		}
		for _, bind := range binds {
			var err error
			if context, err = context.AddBind(bind); err != nil {
				t.Errorf("failed to add bind %v: %v", bind, err)
			}
		}
	}

	return context
}

func typed(term ir.IrTerm, typ ir.IrType) ir.IrTerm {
	term.Type = &typ
	return term
}

func TestTypechecker(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"1 [i8]", typed(
			ir.CallPF(
				typed(ir.Number(1), ir.Forall("a", ir.NewTypeKind(), ir.Tvar("a"))),
				[]ir.IrType{ir.Const("i8")}),
			ir.Const("i8")),
		},
	}

	parser := bplparser2.NewParser()
	parser.SetInitialSymbol("Expression")
	for _, test := range tests {
		parser.Open("testfile", strings.NewReader(test.input))

		got, err := bplparser2.Parse[ir.IrTerm](parser)
		if err != nil {
			t.Fatal(err)
		}

		typechecker := stlc.NewTypechecker(newContext(t))
		if err := typechecker.TypecheckTerm(&got); err != nil {
			t.Fatalf("TypecheckTerm(%v) err = %v; want %v", got, err, nil)
		}

		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty(), cmpopts.IgnoreFields(ir.IrTerm{}, "Pos"), cmpopts.IgnoreFields(ir.IrTerm{}, "Pos"), cmpopts.IgnoreFields(ir.IrType{}, "Pos")) || err != nil {
			t.Errorf("TypecheckTerm(%q) = %v; want %v", test.input, got, test.want)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
