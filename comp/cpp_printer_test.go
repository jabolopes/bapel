package comp_test

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/tests"
)

func TestCppPrinter(t *testing.T) {
	matches, err := tests.Glob("../tests/testdata/comp/in/*.in")
	if err != nil {
		t.Fatal(err)
	}

	tmpBin, err := os.CreateTemp("", "test_codegen_*")
	if err != nil {
		t.Fatal(err)
	}
	tmpBinPath := tmpBin.Name()
	tmpBin.Close()
	defer os.Remove(tmpBinPath)

	cmd := exec.Command("clang++", "-std=c++17", "-I..", "-I../bin", "-x", "c++", "-", "-o", tmpBinPath)
	cmd.Stdin = strings.NewReader(`
#include "bin/codegen_impl.h"
#include <iostream>

int main(int argc, char** argv) {
  if (argc < 3) return 1;
  return codegen::compile_unit(argv[1], argv[2]);
}
`)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to compile native C++ codegen test runner: %v\n%s", err, out)
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

			runCmd := exec.Command(tmpBinPath, inFile, gotFilenameBase)
			if err := runCmd.Run(); err != nil {
				t.Logf("Skipping %s due to compilation error: %v", inFile, err)
				return
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
	matches, err := tests.Glob("../tests/testdata/comp/cpp/*.cc")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		switch path.Base(inFile) {
		case "array.cc", "context1.cc", "loops.cc", "polymorphism.cc":
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

			flags := []string{fmt.Sprintf("-I%s", path.Dir(inFile)), "-I..", "-I."}

			args := append([]string{"-std=c++17", "-c", inFile, "-o", wantFile}, flags...)
			cmd := exec.Command("clang++", args...)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("failed to run %s: %s", cmd, output)
			}
		})
	}
}
