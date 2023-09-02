package bplparser

import (
	"os"
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func TestParseSection(t *testing.T) {
	tests := []struct {
		input string
		want  Source
	}{
		{"imports {", NewSectionSource("imports")},
		{"decls {", NewSectionSource("decls")},
		{"exports {", NewSectionSource("exports")},
	}

	p := NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		p.SetLine(test.input)
		if section, err := p.parseSection(); section != test.want || err != nil {
			t.Errorf("parseSection(%q) = %v, %v; want %v, %v", test.input, section, err, test.want, nil)
		}
	}
}
