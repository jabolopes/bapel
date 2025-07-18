package comp_test

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/tests"
)

func TestTypecheckSourceFile(t *testing.T) {
	matches, err := tests.Glob("testdata/in/*.in")
	if err != nil {
		t.Fatal(err)
	}

	workspace := ast.NewWorkspace(ast.NewPackages([]ast.Package{
		ast.NewPrefixPackage(ir.NewModuleID("", ir.Pos{}), ir.NewFilename("../", ir.Pos{}), ir.Pos{}),
	}, ir.Pos{}))

	querier, err := query.NewWithWorkspace(workspace)
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		wantFile := strings.Replace(parse.ReplaceExtension(inFile, ".out"), "/in/", "/typecheck/", 1)

		options := comp.TypecheckOptions{}
		if path.Base(inFile) == "order.in" {
			options.SkipUndefinedTermChecks = true
		}

		sourceFile, err := comp.TypecheckSourceFile(querier, options, inFile)

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
