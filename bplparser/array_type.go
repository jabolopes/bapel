package bplparser

import (
	"fmt"
	"math"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseArrayType(args []string, named bool) (ir.IrArrayType, []string, error) {
	args, err := parser.ShiftIf(args, "[", fmt.Errorf("expected token '['; got %v", args))
	if err != nil {
		return ir.IrArrayType{}, nil, err
	}

	typ, args, err := ParseType(args, named)
	if err != nil {
		return ir.IrArrayType{}, nil, err
	}

	args, err = parser.ShiftIf(args, "]", fmt.Errorf("expected token ']'; got %v", args))
	if err != nil {
		return ir.IrArrayType{}, nil, err
	}

	return ir.IrArrayType{typ, math.MaxInt}, args, nil
}
