package bplparser2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jabolopes/bapel/bplparser2"
	"github.com/kylelemons/godebug/diff"
)

func TestParsePos(t *testing.T) {
	parser, err := bplparser2.New()
	if err != nil {
		t.Fatal(err)
	}

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

		got, err := bplparser2.ParseWith(parser, in.Name(), in)
		if err != nil {
			t.Fatal(err)
		}

		want, err := os.ReadFile(wantFile)
		if err != nil {
			t.Fatal(err)
		}

		if diff := diff.Diff(string(want), fmt.Sprintf("%+s\n", got)); len(diff) > 0 {
			t.Fatalf("Diff(%q, %q) =\n%s", inFile, wantFile, diff)
		}
	}
}
