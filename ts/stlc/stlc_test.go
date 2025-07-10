package stlc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/tests"
	"github.com/jabolopes/bapel/ts/stlc"
)

func checkSourceFile(filename string, typecheck bool) (ast.SourceFile, error) {
	context := stlc.NewContext()

	sourceFile, err := parse.ParseSourceFile(filename)
	if err != nil {
		return ast.SourceFile{}, err
	}

	if !sourceFile.Valid() {
		return ast.SourceFile{}, sourceFile.Error()
	}

	for i := range sourceFile.Body {
		source := &sourceFile.Body[i]

		switch source.Case {
		case ast.DeclSource:
			context, err = context.AddSymbol(source.Decl.Decl)
			if err != nil {
				return ast.SourceFile{}, err
			}

		case ast.FunctionSource:
			typechecker := stlc.NewTypechecker(context)

			if !typecheck {
				var err error
				context, err = typechecker.InferFunction(source.Function)
				if err != nil {
					return ast.SourceFile{}, err
				}
			} else {
				var err error
				if _, err = typechecker.InferFunction(source.Function); err != nil {
					return ast.SourceFile{}, err
				}

				context, err = typechecker.TypecheckFunction(source.Function)
				if err != nil {
					return ast.SourceFile{}, err
				}
			}
		}
	}

	return sourceFile, nil
}

func TestInferTerm(t *testing.T) {
	matches, err := tests.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		wantFile := fmt.Sprintf("%s.out", strings.TrimSuffix(inFile, ".in"))

		sourceFile, err := checkSourceFile(inFile, false /* typecheck */)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		got := fmt.Sprintf("%+s\n", sourceFile)

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
		if _, err := checkSourceFile(inFile, true /* typecheck */); err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
	}
}
