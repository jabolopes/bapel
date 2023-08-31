package bplparser_test

import (
	"os"
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
)

func TestParseSection(t *testing.T) {
	tests := []struct {
		input string
		want  bplparser.Source
	}{
		{"imports {", bplparser.NewSectionSource("imports")},
		{"decls {", bplparser.NewSectionSource("decls")},
		{"exports {", bplparser.NewSectionSource("exports")},
	}

	p := bplparser.NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		p.SetLine(test.input)
		if section, err := p.ParseSection(); section != test.want || err != nil {
			t.Errorf("ParseSection(%q) = %v, %v; want %v, %v", test.input, section, err, test.want, nil)
		}
	}
}
