package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func ParseType(args []string, named bool) (ir.IrType, []string, error) {
	if len(args) <= 0 {
		return ir.IrType{}, nil, fmt.Errorf("expected type; got %v", args)
	}

	if args[0] == "(" {
		typ, args, err := ParseFunctionType(args, named)
		if err != nil {
			return ir.IrType{}, nil, err
		}

		return ir.NewFunctionType(typ), args, nil
	}

	if args[0] == "[" {
		typ, args, err := ParseArrayType(args, named)
		if err != nil {
			return ir.IrType{}, nil, err
		}

		// TODO: Finish. Return ArrayType instead of ElementType.
		return ir.NewArrayType(typ), args, nil
	}

	typ, err := ir.ParseIntType(args[0])
	if err != nil {
		return ir.IrType{}, nil, err
	}

	return ir.NewIntType(typ), args[1:], nil
}
