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
	testImport = `import c.bpl`
)

func makePos(lineNum int, line string) ir.Pos {
	return ir.Pos{"testfile", lineNum, lineNum, line}
}

func newImportSource(pos ir.Pos, id string) bplparser.Source {
	source := bplparser.NewImportSource(id)
	source.Pos = pos
	return source
}

func TestParseImport(t *testing.T) {
	tests := []struct {
		input string
		want  bplparser.Source
	}{
		{testImport, newImportSource(makePos(1, testImport), "c.bpl")},
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
