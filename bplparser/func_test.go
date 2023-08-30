package bplparser_test

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
)

func TestParseFunc(t *testing.T) {
	tests := []struct {
		input    string
		wantArgs []ir.IrDecl
		wantRets []ir.IrDecl
	}{
		{"func f() -> () {", nil, nil},
		{"func f(a i32) -> () {", []ir.IrDecl{
			ir.NewTermDecl("a", ir.NewIntType(ir.I32)),
		}, nil},
		{"func f() -> (r i64) {", nil,
			[]ir.IrDecl{
				ir.NewTermDecl("r", ir.NewIntType(ir.I64)),
			}},
		{"func f(a [i32], b i64) -> () {", []ir.IrDecl{
			ir.NewTermDecl("a", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
			ir.NewTermDecl("b", ir.NewIntType(ir.I64)),
		}, nil},
		{"func f(a [i32], b i64) -> (r1 i32, r2 [i64]) {",
			[]ir.IrDecl{
				ir.NewTermDecl("a", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I32), math.MaxInt})),
				ir.NewTermDecl("b", ir.NewIntType(ir.I64)),
			},
			[]ir.IrDecl{
				ir.NewTermDecl("r1", ir.NewIntType(ir.I32)),
				ir.NewTermDecl("r2", ir.NewArrayType(ir.IrArrayType{ir.NewIntType(ir.I64), math.MaxInt})),
			},
		},
	}

	parser := bplparser.NewParser(ir.NewCompiler(os.Stdout))
	for _, test := range tests {
		parser.SetLine(test.input)
		id, argTuple, retTuple, err := parser.ParseFunc()
		if id != "f" ||
			!reflect.DeepEqual(argTuple, test.wantArgs) ||
			!reflect.DeepEqual(retTuple, test.wantRets) ||
			err != nil {
			t.Errorf("ParseFunc() = %v, %v, %v, %v; want %v, %v, %v, %v",
				id, argTuple, retTuple, err, "f", test.wantArgs, test.wantRets, nil)
		}
	}
}
