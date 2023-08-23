package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func ParseFunctionType(args []string, named bool) (ir.IrFunctionType, []string, error) {
	orig := args

	argTuple, retTuple, args, err := ParseTupleArrow(args, named)
	if err != nil {
		return ir.IrFunctionType{}, orig, err
	}

	argTypes := make([]ir.IrType, len(argTuple))
	for i := range argTuple {
		argTypes[i] = argTuple[i].Type
	}

	retTypes := make([]ir.IrType, len(retTuple))
	for i := range retTuple {
		retTypes[i] = retTuple[i].Type
	}

	return ir.IrFunctionType{argTypes, retTypes}, args, nil
}
