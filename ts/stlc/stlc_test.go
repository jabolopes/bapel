package stlc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/tests"
	"github.com/jabolopes/bapel/ts/stlc"
)

func checkModule(filename string, typecheck bool) (ast.Module, error) {
	context := stlc.NewContext()

	module, err := bplparser2.ParseModuleFile(filename)
	if err != nil {
		return ast.Module{}, err
	}

	for i := range module.Body {
		source := &module.Body[i]

		switch source.Case {
		case ast.DeclSource:
			context, err = context.AddSymbol(source.Decl.Decl, stlc.DeclSymbol)
			if err != nil {
				return ast.Module{}, err
			}

		case ast.FunctionSource:
			typechecker := stlc.NewTypechecker(context)

			if !typecheck {
				var err error
				context, err = typechecker.InferFunction(source.Function)
				if err != nil {
					return ast.Module{}, err
				}
			} else {
				var err error
				if _, err = typechecker.InferFunction(source.Function); err != nil {
					return ast.Module{}, err
				}

				context, err = typechecker.TypecheckFunction(source.Function)
				if err != nil {
					return ast.Module{}, err
				}
			}
		}
	}

	return module, nil
}

func TestInferTerm(t *testing.T) {
	matches, err := tests.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		wantFile := fmt.Sprintf("%s.out", strings.TrimSuffix(inFile, ".in"))

		module, err := checkModule(inFile, false /* typecheck */)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		got := fmt.Sprintf("%+s\n", module)

		diff, err := tests.DiffOutRegen(got, wantFile)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
		if len(diff) > 0 {
			t.Errorf("in test %s: diff = %s", inFile, diff)
		}
	}
}

func TestTypecheckTerm(t *testing.T) {
	matches, err := tests.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		if _, err := checkModule(inFile, true /* typecheck */); err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
	}
}
