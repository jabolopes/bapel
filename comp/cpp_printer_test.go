package comp_test

import (
	"fmt"
	"os"
	"os/exec"
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
			gotDir, err := os.MkdirTemp("", "bapel-test-*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(gotDir)

			baseName := parse.TrimExtension(path.Base(inFile))
			gotFilenameBase := path.Join(gotDir, baseName)

			if err := comp.CompileBPL(querier, inFile, gotFilenameBase); err != nil {
				if strings.Contains(err.Error(), "failed to typecheck") {
					// Skip generating C++ for any tests that do not typecheck.
					return
				}
				t.Fatalf("CompileBPL failed: %v", err)
			}

			wantFileH := strings.Replace(parse.ReplaceExtension(inFile, ".h"), "/in/", "/cpp/", 1)
			wantFilePrivH := strings.Replace(parse.ReplaceExtension(inFile, "_private.h"), "/in/", "/cpp/", 1)
			wantFileCc := strings.Replace(parse.ReplaceExtension(inFile, ".cc"), "/in/", "/cpp/", 1)

			if diff, err := tests.DiffOutRegenFile(gotFilenameBase+".h", wantFileH); err != nil {
				t.Fatal(err)
			} else if len(diff) > 0 {
				t.Errorf(".h diff = %s", diff)
			}

			if diff, err := tests.DiffOutRegenFile(gotFilenameBase+"_private.h", wantFilePrivH); err != nil {
				t.Fatal(err)
			} else if len(diff) > 0 {
				t.Errorf("_private.h diff = %s", diff)
			}

			if diff, err := tests.DiffOutRegenFile(gotFilenameBase+".cc", wantFileCc); err != nil {
				t.Fatal(err)
			} else if len(diff) > 0 {
				t.Errorf(".cc diff = %s", diff)
			}
		})
	}
}

func TestCppPrinterIsValidCpp(t *testing.T) {
	matches, err := tests.Glob("testdata/cpp/*.cc")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		switch path.Base(inFile) {
		case "array.cc", "context1.cc", "polymorphism.cc":
			// TODO: These tests import 'bapel.core'. Figure out a way to
			// make these tests pass.
			continue
		}

		t.Run(inFile, func(t *testing.T) {
			t.Parallel()

			tmpFile, err := os.CreateTemp("", "*.o")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			wantFile := tmpFile.Name()
			tmpFile.Close()

			flags := []string{fmt.Sprintf("-I%s", path.Dir(inFile))}

			args := append([]string{"-std=c++17", "-c", inFile, "-o", wantFile}, flags...)
			cmd := exec.Command("clang++", args...)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("failed to run %s: %s", cmd, output)
			}
		})
	}
}
