package bplparser

import (
	"fmt"
	"unicode"

	"github.com/jabolopes/bapel/ir"
)

func ParseType(args []string, named bool) (ir.IrType, []string, error) {
	orig := args

	if len(args) <= 0 {
		return ir.IrType{}, nil, fmt.Errorf("expected type; got %v", args)
	}

	if args[0] == "(" {
		typ, args, err := ParseFunctionType(args, named)
		if err != nil {
			return ir.IrType{}, orig, err
		}

		return ir.NewFunctionType(typ), args, nil
	}

	if args[0] == "{" {
		typ, args, err := ParseStructType(args, true /* named */)
		if err != nil {
			return ir.IrType{}, orig, err
		}

		return ir.NewStructType(typ), args, nil
	}

	if args[0] == "[" {
		typ, args, err := ParseArrayType(args, named)
		if err != nil {
			return ir.IrType{}, orig, err
		}

		return ir.NewArrayType(typ), args, nil
	}

	if args[0][0] == 'i' {
		typ, err := ir.ParseIntType(args[0])
		if err != nil {
			return ir.IrType{}, orig, err
		}

		return ir.NewIntType(typ), args[1:], nil
	}

	{
		var c rune
		for _, c = range args[0] {
			break
		}

		if unicode.IsLetter(c) {
			return ir.NewIDType(args[0]), args[1:], nil
		}
	}

	return ir.IrType{}, args, fmt.Errorf("expected type; got %v", args)
}
