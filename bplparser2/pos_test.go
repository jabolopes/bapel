package bplparser2_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/jabolopes/bapel/bplparser2"
	"github.com/kylelemons/godebug/diff"
)

var regen bool

func init() {
	flag.BoolVar(&regen, "regen", false, "Whether to regenerate test output files.")
}

func TestParsePos(t *testing.T) {
	parser, err := bplparser2.New()
	if err != nil {
		t.Fatal(err)
	}

	cases := 0

	for i := 1; ; i++ {
		inFile := fmt.Sprintf("pos_test%d.in", i)
		wantFile := fmt.Sprintf("pos_test%d.out", i)

		in, err := os.Open(inFile)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		defer in.Close()

		cases++

		module, err := bplparser2.ParseWith(parser, in.Name(), in)
		if err != nil {
			t.Fatal(err)
		}

		got := fmt.Sprintf("%+s\n", module)

		if regen {
			if err := os.WriteFile(wantFile, []byte(got), 0644); err != nil {
				t.Fatal(err)
			}
		}

		want, err := os.ReadFile(wantFile)
		if err != nil {
			t.Fatal(err)
		}

		if diff := diff.Diff(string(want), got); len(diff) > 0 {
			t.Fatalf("Diff(%q, %q) =\n%s", inFile, wantFile, diff)
		}
	}

	if cases == 0 {
		t.Fatal("Found no tests")
	}
}
