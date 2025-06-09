package stlc_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/stlc"
)

var regen bool

func init() {
	flag.BoolVar(&regen, "regen", false, "Whether to regenerate test output files.")
}

func TestInferTerm(t *testing.T) {
	parser, err := bplparser2.New()
	if err != nil {
		t.Fatal(err)
	}

	context := stlc.NewContext()
	binds := []stlc.Bind{
		stlc.NewTermBind("print", ir.NewForallType("a", ir.NewTypeKind(), ir.NewFunctionType(ir.Types(ir.NewVarType("a")), ir.Types())), stlc.DefSymbol),
		stlc.NewConstBind("i8", ir.NewTypeKind(), stlc.DefSymbol),
		stlc.NewConstBind("i16", ir.NewTypeKind(), stlc.DefSymbol),
		stlc.NewTermBind("+",
			ir.Forall(
				"a", ir.NewTypeKind(), ir.NewFunctionType(ir.NewTupleType([]ir.IrType{ir.NewVarType("a"), ir.NewVarType("a")}), ir.NewVarType("a"))),
			stlc.ImportSymbol),
	}

	for _, bind := range binds {
		var err error
		context, err = context.AddBind(bind)
		if err != nil {
			t.Fatal(err)
		}
	}

	cases := 0

	for i := 1; ; i++ {
		inFile := fmt.Sprintf("inferencer_test%d.in", i)
		wantFile := fmt.Sprintf("inferencer_test%d.out", i)

		in, err := os.Open(inFile)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		defer in.Close()

		cases++

		module, err := bplparser2.ParseWith(parser, in.Name(), in)
		if err != nil {
			t.Fatal(err)
		}

		for i, source := range module.Body {
			switch source.Case {
			case ast.ComponentSource:
				t.Fatal("ComponentSource not yet supported")

			case ast.FunctionSource:
				typechecker := stlc.NewTypechecker(context)
				if err := typechecker.InferFunction(source.Function); err != nil {
					t.Fatal(err)
				}
				module.Body[i] = source

			case ast.TypeDefSource:
				context, err = context.AddAliasBind(source.TypeDef.Decl)
				if err != nil {
					t.Fatal(err)
				}
			}
		}

		got := fmt.Sprintf("%+s\n", module)

		if regen {
			if err := os.WriteFile(wantFile, []byte(got), 0644); err != nil {
				t.Fatal(err)
			}
		}

		want, err := os.ReadFile(wantFile)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(string(want), got); len(diff) > 0 {
			t.Errorf("Infer() diff = %s", diff)
		}

		if cases == 0 {
			t.Fatal("Found no tests")
		}
	}
}
