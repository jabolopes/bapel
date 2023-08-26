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
		input        string
		wantArgTerms []ir.IrTerm
		wantRetTerms []ir.IrTerm
	}{
		{"call f a1",
			[]ir.IrTerm{ir.NewTokenTerm(parser.NewIDToken("call")), ir.NewTokenTerm(parser.NewIDToken("f")), ir.NewTokenTerm(parser.NewIDToken("a1"))},
			nil,
		},
		{
			"r1 r2 <- f a1 a2",
			[]ir.IrTerm{ir.NewTokenTerm(parser.NewIDToken("f")), ir.NewTokenTerm(parser.NewIDToken("a1")), ir.NewTokenTerm(parser.NewIDToken("a2"))},
			[]ir.IrTerm{ir.NewTokenTerm(parser.NewIDToken("r1")), ir.NewTokenTerm(parser.NewIDToken("r2"))},
		},
	}

	for _, test := range tests {
		argTerms, retTerms, args, err := bplparser.ParseCallAssign(parser.Words(test.input))
		if !reflect.DeepEqual(argTerms, test.wantArgTerms) || !reflect.DeepEqual(retTerms, test.wantRetTerms) || !slices.Equal(args, nil) || err != nil {
			t.Errorf("ParseCallAssign(%q) = %v, %v, %v, %v; want %v, %v, %v, %v",
				test.input, argTerms, retTerms, args, err, test.wantArgTerms, test.wantRetTerms, nil, nil)
		}
	}
}
