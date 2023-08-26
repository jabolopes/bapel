package bplparser_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/slices"
)

func TestParseCallAssign(t *testing.T) {
	bplparser.Compiler = ir.NewCompiler(os.Stdout)

	tests := []struct {
		input    string
		wantArgs []ir.IrTerm
		wantRets []string
	}{
		{"call f a1", []ir.IrTerm{ir.NewTokenTerm(parser.NewIDToken("call")), ir.NewTokenTerm(parser.NewIDToken("f")), ir.NewTokenTerm(parser.NewIDToken("a1"))}, nil},
		{"r1 r2 <- f a1 a2", []ir.IrTerm{ir.NewTokenTerm(parser.NewIDToken("f")), ir.NewTokenTerm(parser.NewIDToken("a1")), ir.NewTokenTerm(parser.NewIDToken("a2"))}, []string{"r1", "r2"}},
	}

	for _, test := range tests {
		args, rets, err := bplparser.ParseCallAssign(parser.Words(test.input))
		if !reflect.DeepEqual(args, test.wantArgs) || !slices.Equal(rets, test.wantRets) || err != nil {
			t.Errorf("ParseCallAssign(%q) = %v, %v, %v; want %v, %v, %v",
				test.input, args, rets, err, test.wantArgs, test.wantRets, nil)
		}
	}
}
