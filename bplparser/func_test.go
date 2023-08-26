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

func TestParseFunc(t *testing.T) {
	tests := []struct {
		input    string
		wantArgs []ir.IrDecl
		wantRets []ir.IrDecl
	}{
		{"func f() -> () {", nil, nil},
		{"func f(a i32) -> () {", []ir.IrDecl{
			ir.NewVarDecl("a", ir.NewIntType(ir.I32)),
		}, nil},
		{"func f() -> (r i64) {", nil,
			[]ir.IrDecl{
				ir.NewVarDecl("r", ir.NewIntType(ir.I64)),
			}},
		{"func f(a [i32], b i64) -> () {", []ir.IrDecl{
			ir.NewVarDecl("a", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
			ir.NewVarDecl("b", ir.NewIntType(ir.I64)),
		}, nil},
		{"func f(a [i32], b i64) -> (r1 i32, r2 [i64]) {",
			[]ir.IrDecl{
				ir.NewVarDecl("a", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
				ir.NewVarDecl("b", ir.NewIntType(ir.I64)),
			},
			[]ir.IrDecl{
				ir.NewVarDecl("r1", ir.NewIntType(ir.I32)),
				ir.NewVarDecl("r2", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I64), math.MaxInt})),
			},
		},
	}

	p := bplparser.NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		id, argTuple, retTuple, args, err := p.ParseFunc(parser.Words(test.input))
		if id != "f" ||
			!reflect.DeepEqual(argTuple, test.wantArgs) ||
			!reflect.DeepEqual(retTuple, test.wantRets) ||
			!slices.Equal(args, nil) ||
			err != nil {
			t.Errorf("ParseFunc() = %v, %v, %v, %v, %v; want %v, %v, %v, %v, %v",
				id, argTuple, retTuple, args, err, "f", test.wantArgs, test.wantRets, nil, nil)
		}
	}
}
