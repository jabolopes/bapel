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
		wantArgs []parser.Token
		wantRets []string
	}{
		{"call f a1", []parser.Token{parser.NewIDToken("call"), parser.NewIDToken("f"), parser.NewIDToken("a1")}, nil},
		{"r1 r2 <- f a1 a2", []parser.Token{parser.NewIDToken("f"), parser.NewIDToken("a1"), parser.NewIDToken("a2")}, []string{"r1", "r2"}},
	}

	for _, test := range tests {
		args, rets, err := bplparser.ParseCallAssign(parser.Words(test.input))
		if !slices.Equal(args, test.wantArgs) || !slices.Equal(rets, test.wantRets) || err != nil {
			t.Errorf("ParseCallAssign(%q) = %v, %v, %v; want %v, %v, %v",
				test.input, args, rets, err, test.wantArgs, test.wantRets, nil)
		}
	}
}
