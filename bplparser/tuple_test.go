package bplparser_test

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
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

	parser := bplparser.NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		parser.SetLine(test.input)
		vars, err := parser.ParseTuple(true /* named */, bplparser.Parens)
		if !reflect.DeepEqual(vars, test.want) || err != nil {
			t.Errorf("ParseTuple(%q) = %v, %v; want %v, %v",
				test.input, vars, err, test.want, nil)
		}
	}
}
