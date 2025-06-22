package stlc_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	matches, err := filepath.Glob("inferencer_test_*.in")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) == 0 {
		t.Fatal("Found no tests")
	}

	for _, inFile := range matches {
		context := context

		wantFile := fmt.Sprintf("%s.out", strings.TrimSuffix(inFile, ".in"))

		in, err := os.Open(inFile)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
		defer in.Close()

		module, err := bplparser2.ParseWith(parser, in.Name(), in)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: Deduplicate with compiler.
		//
		// TODO: Finish.
		for i, source := range module.Body {
			switch source.Case {
			case ast.ComponentSource:
				t.Fatalf("in test %s: ComponentSource is not yet supported", inFile)

			case ast.DeclSource:
				context, err = context.AddSymbol(source.Decl.Decl, stlc.DeclSymbol)
				if err != nil {
					t.Fatalf("in test %s: %v", inFile, err)
				}

			case ast.FunctionSource:
				typechecker := stlc.NewTypechecker(context)

				var err error
				context, err = typechecker.InferFunction(source.Function)
				if err != nil {
					t.Fatalf("in test %s: %v", inFile, err)
				}

				module.Body[i] = source

			case ast.DefSymbolSource:
				symbol := stlc.DefSymbol
				if source.DefSymbol.IsDecl {
					symbol = stlc.DeclSymbol
				}
				context, err = context.AddSymbol(source.DefSymbol.Decl, symbol)
				if err != nil {
					t.Fatalf("in test %s: %v", inFile, err)
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
	}
}
