package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

type DelimiterCase int

const (
	Parens = DelimiterCase(iota)
	Brackets
)

func ParseTuple(args []string, varType ir.IrVarType, named bool, delimiter DelimiterCase) ([]ir.IrVar, []string, error) {
	orig := args

	args, remainder := parser.ShiftBalancedParens(args)

	left := "("
	if delimiter == Brackets {
		left = "{"
	}

	right := ")"
	if delimiter == Brackets {
		right = "}"
	}

	args, err := parser.ShiftToken(args, left)
	if err != nil {
		return nil, orig, err
	}

	if _, err := parser.ShiftToken(args, right); err == nil {
		return nil, remainder, nil
	}

	var vars []ir.IrVar
	for {
		var id string
		if named {
			id, args, err = parser.Shift(args, fmt.Errorf("expected identifier; got %v", args))
			if err != nil {
				return nil, orig, err
			}
		}

		var typ ir.IrType
		typ, args, err = ParseType(args, named)
		if err != nil {
			return nil, orig, err
		}

		vars = append(vars, ir.IrVar{Id: id, VarType: varType, Type: typ})

		if args, err = parser.ShiftToken(args, ","); err == nil {
			continue
		}

		args, err = parser.ShiftToken(args, right)
		if err != nil {
			return nil, orig, err
		}

		break
	}

	return vars, remainder, nil
}

func ParseTupleArrow(args []string, named bool) ([]ir.IrVar, []string, error) {
	orig := args

	argVars, args, err := ParseTuple(args, ir.ArgVar, named, Parens)
	if err != nil {
		return nil, orig, fmt.Errorf("in argument list: %v", err)
	}

	args, err = parser.ShiftToken(args, "->")
	if err != nil {
		return nil, orig, err
	}

	retVars, args, err := ParseTuple(args, ir.RetVar, named, Parens)
	if err != nil {
		return nil, orig, fmt.Errorf("in return list: %v", err)
	}

	return append(argVars, retVars...), args, nil
}
