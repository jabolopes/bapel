package bplparser

import (
	"fmt"
	"io"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseTuple(args []string, varType ir.IrVarType, named bool) ([]ir.IrVar, []string, error) {
	args, remainder := parser.ShiftBalancedParens(args)

	args, err := parser.ShiftIf(args, "(", fmt.Errorf("expected token '('; got %v", args))
	if err != nil {
		return nil, nil, err
	}

	if _, err := parser.ShiftIf(args, ")", io.EOF); err == nil {
		return nil, remainder, nil
	}

	var vars []ir.IrVar
	for {
		var id string
		if named {
			id, args, err = parser.Shift(args, fmt.Errorf("expected identifier; got %v", args))
			if err != nil {
				return nil, nil, err
			}
		}

		var typ ir.IrType
		typ, args, err = ParseType(args)
		if err != nil {
			return nil, nil, err
		}

		vars = append(vars, ir.IrVar{Id: id, VarType: varType, Type: typ})

		if args, err = parser.ShiftIf(args, ",", io.EOF); err == nil {
			continue
		}

		args, err = parser.ShiftIf(args, ")", fmt.Errorf("expected token ')'; got %v", args))
		if err != nil {
			return nil, nil, err
		}

		break
	}

	return vars, remainder, nil
}

func ParseTupleArrow(args []string, named bool) ([]ir.IrVar, []string, error) {
	argVars, args, err := ParseTuple(args, ir.ArgVar, named)
	if err != nil {
		return nil, nil, fmt.Errorf("in argument list: %v", err)
	}

	args, err = parser.ShiftIf(args, "->", fmt.Errorf("expected token '->' in return list; got %v", args))
	if err != nil {
		return nil, nil, err
	}

	retVars, args, err := ParseTuple(args, ir.RetVar, named)
	if err != nil {
		return nil, nil, fmt.Errorf("in return list: %v", err)
	}

	return append(argVars, retVars...), args, nil
}

func ParseFunctionType(args []string) (ir.IrFunctionType, []string, error) {
	vars, args, err := ParseTupleArrow(args, false /* named */)
	if err != nil {
		return ir.IrFunctionType{}, nil, err
	}

	typ := ir.IrFunctionType{}
	for _, irvar := range vars {
		switch irvar.VarType {
		case ir.ArgVar:
			typ.Args = append(typ.Args, irvar.Type)
		default:
			typ.Rets = append(typ.Rets, irvar.Type)
		}
	}

	return typ, args, nil
}
