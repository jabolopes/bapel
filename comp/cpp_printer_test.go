package comp_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/build"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/tests"
)

func TestCppPrinter(t *testing.T) {
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
		if path.Base(inFile) == "order.in" {
			continue
		}

		t.Run(inFile, func(t *testing.T) {
			wantFile := strings.Replace(parse.ReplaceExtension(inFile, ".ccm"), "/in/", "/cpp/", 1)
			gotFile, err := os.CreateTemp("", path.Base(inFile))
			if err != nil {
				t.Fatal(err)
			}

			gotFilename := gotFile.Name()
			gotFile.Close()

			defer func() {
				os.Remove(gotFilename)
			}()

			if err := comp.CompileBPLToCCM(querier, inFile, gotFilename); err != nil {
				got := fmt.Sprintf("%s\n", err)
				if err := os.WriteFile(gotFilename, []byte(got), 0660); err != nil {
					t.Fatal(err)
				}
			}

			diff, err := tests.DiffOutRegenFile(gotFilename, wantFile)
			if err != nil {
				t.Fatal(err)
			}
			if len(diff) > 0 {
				t.Errorf("diff = %s", diff)
			}
		})
	}
}

func TestCppPrinterIsValidCpp(t *testing.T) {
	matches, err := tests.Glob("testdata/cpp/*.ccm")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		switch path.Base(inFile) {
		case "context2.ccm":
			// This test fails typechecking, so it doesn't generate C++.
			continue
		case "array.ccm", "context1.ccm", "polymorphism.ccm":
			// TODO: These tests import 'bapel.core'. Figure out a way to
			// make these tests pass.
			continue
		}

		t.Run(inFile, func(t *testing.T) {
			t.Parallel()

			tmpFile, err := os.CreateTemp("", "*.pcm")
			if err != nil {
				t.Fatal(err)
			}
			wantFile := tmpFile.Name()
			tmpFile.Close()
			func() {
				os.Remove(wantFile)
			}()

			cmd, err := build.CompileCCMToPCMCommand(inFile, nil /* flags */, wantFile)
			if err != nil {
				t.Fatal(err)
			}

			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("failed to run %s: %s", cmd, output)
			}
		})
	}
}
