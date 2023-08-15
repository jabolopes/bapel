package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseTuple(args []string, varType ir.IrVarType, named bool) ([]ir.IrVar, error) {
	args, err := parser.ShiftIf(args, "(", fmt.Errorf("expected token '('; got %v", args))
	if err != nil {
		return nil, err
	}

	args, err = parser.ShiftIfEnd(args, ")", fmt.Errorf("expected token ')'; got %v", args))
	if err != nil {
		return nil, err
	}

	var vars []ir.IrVar

	for len(args) > 0 {
		var id string
		if named {
			id, args, err = parser.Shift(args, fmt.Errorf("expected identifier; got %v", args))
			if err != nil {
				return nil, err
			}
		}

		var typStr string
		typStr, args, err = parser.Shift(args, fmt.Errorf("expected type for identifier; got %v", args))
		if err != nil {
			return nil, err
		}

		if len(args) > 0 {
			args, err = parser.ShiftIf(args, ",", fmt.Errorf("expected token ','; got %v", args))
			if err != nil {
				return nil, err
			}
		}

		typ, err := ir.ParseIntType(typStr)
		if err != nil {
			return nil, err
		}

		vars = append(vars, ir.IrVar{Id: id, VarType: varType, Type: ir.NewIntType(typ)})
	}

	return vars, nil
}

func ParseTupleArrow(args []string, named bool) ([]ir.IrVar, error) {
	args, rets := parser.ShiftBalancedParens(args)

	rets, err := parser.ShiftIf(rets, "->", fmt.Errorf("expected token '->' in return list; got %v", rets))
	if err != nil {
		return nil, err
	}

	vars, err := ParseTuple(args, ir.ArgVar, named)
	if err != nil {
		return nil, fmt.Errorf("in argument list: %v", err)
	}

	retVars, err := ParseTuple(rets, ir.RetVar, named)
	if err != nil {
		return nil, fmt.Errorf("in return list: %v", err)
	}

	return append(vars, retVars...), nil
}

func ParseFunctionType(args []string) (ir.IrFunctionType, error) {
	vars, err := ParseTupleArrow(args, false /* named */)
	if err != nil {
		return ir.IrFunctionType{}, err
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

	return typ, nil
}
