package bplparser_test

import (
	"os"
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func TestParseSection(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"imports {", "imports"},
		{"decls {", "decls"},
		{"exports {", "exports"},
	}

	p := bplparser.NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		section, _, err := p.ParseSection(parser.Words(test.input))
		if section != test.want || err != nil {
			t.Errorf("ParseSection(%q) = %v, %v; want %v, %v", test.input, section, err, test.want, nil)
		}
	}
}
