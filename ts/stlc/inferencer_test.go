package stlc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/stlc"
)

func TestInferTerm(t *testing.T) {
	context := stlc.NewContext()
	binds := []stlc.Bind{
		stlc.NewTermBind("print", ir.NewForallType("a", ir.NewTypeKind(), ir.NewFunctionType(ir.Types(ir.NewVarType("a")), ir.Types())), stlc.DefSymbol),
		stlc.NewConstBind("i8", ir.NewTypeKind(), stlc.DefSymbol),
	}

	for _, bind := range binds {
		var err error
		context, err = context.AddBind(bind)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 1; i < 7; i++ {
		inFile := fmt.Sprintf("inferencer_test%d.in", i)
		wantFile := fmt.Sprintf("inferencer_test%d.out", i)

		in, err := os.Open(inFile)
		if err != nil {
			t.Fatal(err)
		}
		defer in.Close()

		want, err := os.ReadFile(wantFile)
		if err != nil {
			t.Fatal(err)
		}

		module, err := bplparser2.ParseFile(in.Name(), in)
		if err != nil {
			t.Fatal(err)
		}

		var inFunction *ir.IrFunction
		for _, source := range module.Body {
			if !source.Is(ast.FunctionSource) {
				continue
			}

			inFunction = source.Function
			break
		}

		if inFunction == nil {
			t.Fatal("Missing in function")
		}

		typechecker := stlc.NewTypechecker(context)
		if err := typechecker.InferFunction(inFunction); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(fmt.Sprintf("%v\n", inFunction), string(want)); len(diff) > 0 {
			t.Errorf("Infer() diff = %s", diff)
		}
	}
}
