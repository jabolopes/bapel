package comp_test

import (
	"fmt"
	"testing"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/tests"
)

func TestTypecheckSourceFile(t *testing.T) {
	matches, err := tests.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		wantFile := parse.ReplaceExtension(inFile, ".out")

		workspace := ast.NewWorkspace(ast.NewPackages([]ast.Package{
			ast.NewPrefixPackage(ir.NewModuleID("", ir.Pos{}), ir.NewFilename("../", ir.Pos{}), ir.Pos{}),
		}, ir.Pos{}))

		querier, err := query.NewWithWorkspace(workspace)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		sourceFile, err := comp.TypecheckSourceFile(querier, comp.TypecheckOptions{}, inFile)

		var got string
		if err == nil {
			got = fmt.Sprintf("%+s\n", sourceFile)
		} else {
			got = fmt.Sprintf("%s\n", err)
		}

		diff, err := tests.DiffOutRegen(got, wantFile)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
		if len(diff) > 0 {
			t.Errorf("in test %s: diff = %s", inFile, diff)
		}
	}
}
