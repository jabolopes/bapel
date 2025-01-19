package bplparser2_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

const (
	testImport = `imports {
    core
  }`

	testImports = `imports {
    core
    vec
  }`
)

func makePos(beginLineNum, endLineNum int, line string) ir.Pos {
	// TODO: Fix Line field.
	return ir.Pos{"testfile", beginLineNum, endLineNum, ""}
}

func newImportSource(pos ir.Pos, ids ...ast.ID) ast.Source {
	source := ast.NewImportsSource(ids)
	source.Pos = pos
	return source
}

func TestParseImport(t *testing.T) {
	core := ast.ID{makePos(2, 2, "core"), "core"}
	vec := ast.ID{makePos(3, 3, "vec"), "vec"}

	tests := []struct {
		input string
		want  ast.Source
	}{
		{testImport, newImportSource(makePos(1, 3, testImport), core)},
		{testImports, newImportSource(makePos(1, 4, testImports), core, vec)},
	}

	parser, err := bplparser2.New()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		parser.Open("testfile", strings.NewReader(test.input))

		want := []ast.Source{test.want}
		got, err := bplparser2.Parse[[]ast.Source](parser)
		if !cmp.Equal(got, want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", test.input, got, err, want, nil)
			t.Fatalf("Diff = %v", cmp.Diff(got, want, cmpopts.EquateEmpty()))
		}
	}
}
