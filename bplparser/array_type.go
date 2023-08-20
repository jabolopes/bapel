package bplparser

import (
	"math"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseArrayType(args []string, named bool) (ir.IrArrayType, []string, error) {
	args, err := parser.ShiftToken(args, "[")
	if err != nil {
		return ir.IrArrayType{}, nil, err
	}

	typ, args, err := ParseType(args, named)
	if err != nil {
		return ir.IrArrayType{}, nil, err
	}

	length, args, err := parser.ShiftNumber[int](args)
	if err != nil {
		length = math.MaxInt
	}

	args, err = parser.ShiftToken(args, "]")
	if err != nil {
		return ir.IrArrayType{}, nil, err
	}

	return ir.IrArrayType{typ, length}, args, nil
}
