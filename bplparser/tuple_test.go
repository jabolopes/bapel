package bplparser

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func TestParseTuple(t *testing.T) {
	tests := []struct {
		input string
		want  []ir.IrDecl
	}{
		{"()", nil},
		{"(a i32)", []ir.IrDecl{
			ir.NewTermDecl("a", ir.NewNameType("i32")),
		}},
		{"(r i64)", []ir.IrDecl{
			ir.NewTermDecl("r", ir.NewNameType("i64")),
		}},
		{"(a [i32])", []ir.IrDecl{
			ir.NewTermDecl("a", ir.NewArrayType(ir.NewNameType("i32"), math.MaxInt)),
		}},
		{"(a [i64], b i32)", []ir.IrDecl{
			ir.NewTermDecl("a", ir.NewArrayType(ir.NewNameType("i64"), math.MaxInt)),
			ir.NewTermDecl("b", ir.NewNameType("i32")),
		}},
	}

	parser := NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		parser.SetLine(test.input)
		if vars, err := parser.parseTuple(true /* named */, Parens); !reflect.DeepEqual(vars, test.want) || err != nil {
			t.Errorf("parseTuple(%q) = %v, %v; want %v, %v",
				test.input, vars, err, test.want, nil)
		}
	}
}
