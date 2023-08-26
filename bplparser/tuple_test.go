package bplparser_test

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/slices"
)

func TestParseTuple(t *testing.T) {
	tests := []struct {
		input string
		want  []ir.IrDecl
	}{
		{"()", nil},
		{"(a i32)", []ir.IrDecl{
			ir.NewVarDecl("a", ir.NewIntType(ir.I32)),
		}},
		{"(r i64)", []ir.IrDecl{
			ir.NewVarDecl("r", ir.NewIntType(ir.I64)),
		}},
		{"(a [i32])", []ir.IrDecl{
			ir.NewVarDecl("a", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
		}},
		{"(a [i64], b i32)", []ir.IrDecl{
			ir.NewVarDecl("a", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I64), math.MaxInt})),
			ir.NewVarDecl("b", ir.NewIntType(ir.I32)),
		}},
	}

	p := bplparser.NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		vars, args, err := p.ParseTuple(parser.Words(test.input), true /* named */, bplparser.Parens)
		if !reflect.DeepEqual(vars, test.want) || !slices.Equal(args, nil) || err != nil {
			t.Errorf("ParseTuple(%q) = %v, %v, %v; want %v, %v, %v",
				test.input, vars, args, err, test.want, nil, nil)
		}
	}
}
