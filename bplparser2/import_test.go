package bplparser2_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/bplparser2"
)

const (
	testImport = `import c.bpl`
)

func TestParseImport(t *testing.T) {
	tests := []struct {
		input string
		want  []bplparser.Source
	}{
		{testImport, []bplparser.Source{bplparser.NewImportSource("c.bpl")}},
	}

	parser := bplparser2.NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))

		got, err := bplparser2.Parse[[]bplparser.Source](parser)
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, test.want, cmpopts.EquateEmpty()))
		}
	}
}
