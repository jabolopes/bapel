package bplparser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	testImport = `import c.bpl`
)

func TestImport(t *testing.T) {
	tests := []struct {
		input string
		want  Source
	}{
		{testImport, NewImportSource("c.bpl")},
	}

	parser := NewParser()
	for _, test := range tests {
		parser.Open(strings.NewReader(test.input))
		parser.Scan()
		if got, err := parser.parseAny(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("parseSection(%q) = %v, %v; want %v, %v", test.input, got, err, test.want, nil)
		}
	}
}
