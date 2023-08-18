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

func TestParseFunc(t *testing.T) {
	tests := []struct {
		input string
		want  []ir.IrVar
	}{
		{"f() -> ()", nil},
		{"f(a i32) -> ()", []ir.IrVar{
			ir.NewVar("a", ir.ArgVar, ir.NewIntType(ir.I32)),
		}},
		{"f() -> (r i64)", []ir.IrVar{
			ir.NewVar("r", ir.RetVar, ir.NewIntType(ir.I64)),
		}},
		{"f(a [i32], b i64) -> ()", []ir.IrVar{
			ir.NewVar("a", ir.ArgVar, ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
			ir.NewVar("b", ir.ArgVar, ir.NewIntType(ir.I64)),
		}},
		{"f(a [i32], b i64) -> (r1 i32, r2 [i64])", []ir.IrVar{
			ir.NewVar("a", ir.ArgVar, ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
			ir.NewVar("b", ir.ArgVar, ir.NewIntType(ir.I64)),
			ir.NewVar("r1", ir.RetVar, ir.NewIntType(ir.I32)),
			ir.NewVar("r2", ir.RetVar, ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I64), math.MaxInt})),
		}},
	}

	for _, test := range tests {
		id, vars, args, err := bplparser.ParseFunc(parser.Words(test.input))
		if id != "f" || !reflect.DeepEqual(vars, test.want) || !slices.Equal(args, nil) || err != nil {
			t.Errorf("ParseFunc() = %v, %v, %v, %v; want %v, %v, %v, %v",
				id, vars, args, err, "f", test.want, nil, nil)
		}
	}
}
