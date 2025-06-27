package comp_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/tests"
)

func TestResolver(t *testing.T) {
	matches, err := tests.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		wantFile := fmt.Sprintf("%s.out", strings.TrimSuffix(inFile, ".in"))

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

		diff, err := tests.DiffOutRegen(got, wantFile)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
		if len(diff) > 0 {
			t.Errorf("in test %s: diff = %s", inFile, diff)
		}
	}
}
