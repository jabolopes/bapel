package bplparser2_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jabolopes/bapel/bplparser2"
)

var regen bool

func init() {
	flag.BoolVar(&regen, "regen", false, "Whether to regenerate test output files.")
}

func TestParser(t *testing.T) {
	matches, err := filepath.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) == 0 {
		t.Fatal("Found no tests")
	}

	for _, inFile := range matches {
		wantFile := fmt.Sprintf("%s.out", strings.TrimSuffix(inFile, ".in"))

		in, err := os.Open(inFile)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		module, err := bplparser2.ParseFile(inFile, in)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		got := fmt.Sprintf("%+s\n", module)

		if regen {
			if err := os.WriteFile(wantFile, []byte(got), 0644); err != nil {
				t.Fatalf("in test %s: %v", inFile, err)
			}
		}

		want, err := os.ReadFile(wantFile)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		if diff := cmp.Diff(string(want), got); len(diff) > 0 {
			t.Errorf("in test %s: diff = %s", inFile, diff)
		}
	}
}
