package comp_test

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/tests"
)

var regen bool

func init() {
	flag.BoolVar(&regen, "regen", false, "Whether to regenerate test output files.")
}

func TestResolver(t *testing.T) {
	matches, err := tests.Glob("../testdata/*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		wantFile := fmt.Sprintf("%s.resolver.out", strings.TrimSuffix(inFile, ".in"))

		in, err := os.Open(inFile)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
		defer in.Close()

		module, err := bplparser2.ParseFile(inFile, in)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		if err := comp.ResolveModule(&module); err != nil {
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
