package bplparser_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/slices"
)

func TestParseTuple(t *testing.T) {
	tests := []struct {
		input   string
		varType ir.IrVarType
		want    []ir.IrVar
	}{
		{"()", ir.ArgVar, nil},
		{"(a i32)", ir.ArgVar, []ir.IrVar{
			ir.NewVar("a", ir.ArgVar, ir.NewIntType(ir.I32)),
		}},
		{"(r i64)", ir.RetVar, []ir.IrVar{
			ir.NewVar("r", ir.RetVar, ir.NewIntType(ir.I64)),
		}},
		{"(a [i32])", ir.ArgVar, []ir.IrVar{
			ir.NewVar("a", ir.ArgVar, ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
		}},
		{"(a [i64], b i32)", ir.ArgVar, []ir.IrVar{
			ir.NewVar("a", ir.ArgVar, ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I64), math.MaxInt})),
			ir.NewVar("b", ir.ArgVar, ir.NewIntType(ir.I32)),
		}},
	}

	for _, test := range tests {
		vars, args, err := bplparser.ParseTuple(parser.Words(test.input), test.varType, true /* named */, bplparser.Parens)
		if !reflect.DeepEqual(vars, test.want) || !slices.Equal(args, nil) || err != nil {
			t.Errorf("ParseTuple(%q) = %v, %v, %v; want %v, %v, %v",
				test.input, vars, args, err, test.want, nil, nil)
		}
	}
}
