package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"golang.org/x/exp/slices"
)

func ParseType(args []string) (ir.IrType, error) {
	if slices.Contains(args, "->") {
		typ, err := ParseFunctionType(args)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewFunctionType(typ), nil
	}

	if len(args) == 1 {
		typ, err := ir.ParseIntType(args[0])
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewIntType(typ), nil
	}

	return ir.IrType{}, fmt.Errorf("expected type; got %v", args)
}
