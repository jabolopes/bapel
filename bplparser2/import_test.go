package bplparser2_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

const (
	testImport = `imports {
    c.bpl
  }`

	testImports = `imports {
    c.bpl
    vector.bpl
  }`
)

func makePos(beginLineNum, endLineNum int, line string) ir.Pos {
	// TODO: Fix Line field.
	return ir.Pos{"testfile", beginLineNum, endLineNum, ""}
}

func newImportSource(pos ir.Pos, ids ...string) bplparser.Source {
	source := bplparser.NewImportsSource(ids)
	source.Pos = pos
	return source
}

func TestParseImport(t *testing.T) {
	tests := []struct {
		input string
		want  bplparser.Source
	}{
		{testImport, newImportSource(makePos(1, 3, testImport), "c.bpl")},
		{testImports, newImportSource(makePos(1, 4, testImports), "c.bpl", "vector.bpl")},
	}

	parser := bplparser2.NewParser()
	for _, test := range tests {
		parser.Open("testfile", strings.NewReader(test.input))

		want := []bplparser.Source{test.want}
		got, err := bplparser2.Parse[[]bplparser.Source](parser)
		if !cmp.Equal(got, want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, want, cmpopts.EquateEmpty()))
		}
	}
}
