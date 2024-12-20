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
	testExports = `exports {
  f : i32
}
`
	testDecls = `decls {
  f : i32
}
`
)

func newSection(id string) bplparser.Source {
	return bplparser.NewSectionSource(id, []ir.IrDecl{ir.NewTermDecl("f", ir.NewNameType("i32"))})
}

func TestParseSection(t *testing.T) {
	tests := []struct {
		input string
		want  []bplparser.Source
	}{
		{testExports, []bplparser.Source{newSection("exports")}},
		{testDecls, []bplparser.Source{newSection("decls")}},
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
