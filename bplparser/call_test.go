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

func newID(id string) ir.IrTerm {
	return ir.NewTokenTerm(parser.NewIDToken(id))
}

func newNumber(value int64) ir.IrTerm {
	return ir.NewTokenTerm(parser.NewNumberToken(value))
}

func TestParseCall(t *testing.T) {
	tests := []struct {
		input string
		want  ir.IrTerm
	}{
		{"a b", ir.NewTupleTerm([]ir.IrTerm{newID("a"), newID("b")})},
		{"Index.get a 1", ir.NewIndexGetTerm(newID("a"), newNumber(1))},
		{"Index.set a 1 10", ir.NewIndexSetTerm(newID("a"), newNumber(1), newNumber(10))},
		{"- a", ir.NewOpUnaryTerm("-", newID("a"))},
		{"a + b", ir.NewOpBinaryTerm("+", newID("a"), newID("b"))},
		{"widen a", ir.NewWidenTerm(newID("a"))},
		// {
		// 	"r1 r2 <- f a1 a2",
		// 	[]ir.IrTerm{ir.NewTokenTerm(parser.NewIDToken("f")), ir.NewTokenTerm(parser.NewIDToken("a1")), ir.NewTokenTerm(parser.NewIDToken("a2"))},
		// 	[]ir.IrTerm{ir.NewTokenTerm(parser.NewIDToken("r1")), ir.NewTokenTerm(parser.NewIDToken("r2"))},
		// },
	}

	p := bplparser.NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		got, args, err := p.ParseCall(parser.Words(test.input))
		if !reflect.DeepEqual(got, test.want) || !slices.Equal(args, nil) || err != nil {
			t.Errorf("ParseCall(%q) = %v, %v, %v; want %v, %v, %v",
				test.input, got, args, err, test.want, nil, nil)
		}
	}
}
