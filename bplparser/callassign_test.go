package bplparser_test

import (
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/slices"
)

func TestParseCallAssign(t *testing.T) {
	tests := []struct {
		input    string
		wantArgs []string
		wantRets []string
	}{
		{"call f a1", []string{"call", "f", "a1"}, nil},
		{"r1 r2 <- f a1 a2", []string{"f", "a1", "a2"}, []string{"r1", "r2"}},
	}

	for _, test := range tests {
		args, rets, err := bplparser.ParseCallAssign(parser.Words(test.input))
		if !slices.Equal(args, test.wantArgs) || !slices.Equal(rets, test.wantRets) || err != nil {
			t.Errorf("ParseCallAssign(%q) = %v, %v, %v; want %v, %v, %v",
				test.input, args, rets, err, test.wantArgs, test.wantRets, nil)
		}
	}
}
