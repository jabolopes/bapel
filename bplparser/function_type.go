package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func ParseFunctionType(args []string, named bool) (ir.IrFunctionType, []string, error) {
	vars, args, err := ParseTupleArrow(args, named)
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
