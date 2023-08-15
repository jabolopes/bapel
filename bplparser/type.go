package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"golang.org/x/exp/slices"
)

func ParseType(args []string) (ir.IrType, []string, error) {
	if slices.Contains(args, "->") {
		typ, args, err := ParseFunctionType(args)
		if err != nil {
			return ir.IrType{}, nil, err
		}

		return ir.NewFunctionType(typ), args, nil
	}

	if len(args) > 0 {
		typ, err := ir.ParseIntType(args[0])
		if err != nil {
			return ir.IrType{}, nil, err
		}

		return ir.NewIntType(typ), args[1:], nil
	}

	return ir.IrType{}, nil, fmt.Errorf("expected type; got %v", args)
}
