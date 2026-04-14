package parse_test

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/tests"
)

func TestParser(t *testing.T) {
	t.Parallel()

	matches, err := tests.Glob("testdata/in/*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		t.Run(inFile, func(t *testing.T) {
			t.Parallel()

			wantFile := strings.Replace(parse.ReplaceExtension(inFile, ".bpl"), "/in/", "/parsed/", 1)
			wantErr := strings.HasPrefix(path.Base(inFile), "bad_")

			module, err := parse.ParseSourceFile(inFile)
			if !wantErr && err != nil {
				t.Fatal(err)
			} else if wantErr && err == nil {
				t.Fatal("expected parse error but it succeeded")
			}

			var got string
			if wantErr {
				got = fmt.Sprintf("%s\n", err.Error())
			} else {
				got = fmt.Sprintf("%+s\n", module)
			}

			diff, err := tests.DiffOutRegen(got, wantFile)
			if err != nil {
				t.Fatal(err)
			}
			if len(diff) > 0 {
				t.Errorf("diff = %s", diff)
			}
		})
	}
}
