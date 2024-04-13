package bplparser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
)

const (
	testImports = `imports {
  f : i32
}
`

	emptyImports = `imports {
}
`

	testExports = `exports {
  f : i32
}
`

	emptyExports = `exports {
}
`

	testDecls = `decls {
  f : i32
}
`

	emptyDecls = `decls {
}
`
)

func newSection(id string) Source {
	return NewSectionSource(id, []ir.IrDecl{ir.NewTermDecl("f", ir.NewNameType("i32"))})
}

func newEmptySection(id string) Source {
	return NewSectionSource(id, nil)
}

func TestParseSection(t *testing.T) {
	tests := []struct {
		input string
		want  Source
	}{
		{testImports, newSection("imports")},
		{emptyImports, newEmptySection("imports")},
		{testExports, newSection("exports")},
		{emptyExports, newEmptySection("exports")},
		{testDecls, newSection("decls")},
		{emptyDecls, newEmptySection("decls")},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))
		parser.Scan()
		if got, err := parser.parseSection(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("parseSection(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
		}
	}
}
