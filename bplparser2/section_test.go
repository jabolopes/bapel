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
)

func newNameType() ir.IrType {
	typ := ir.NewNameType("i32")
	// TODO: Fix Line field.
	typ.Pos = ir.Pos{"testfile", 2, 2, ""}
	return typ
}

func newTermDecl() ir.IrDecl {
	decl := ir.NewTermDecl("f", newNameType())
	// TODO: Fix Line field.
	decl.Pos = ir.Pos{"testfile", 2, 2, ""}
	return decl
}

func newExports() bplparser.Source {
	source := bplparser.NewExportsSource([]ir.IrDecl{newTermDecl()})
	// TODO: Fix Line field.
	source.Pos = ir.Pos{"testfile", 1, 3, ""}
	return source
}

func TestParseSection(t *testing.T) {
	tests := []struct {
		input string
		want  bplparser.Source
	}{
		{testExports, newExports()},
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
