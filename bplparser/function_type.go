package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func ParseFunctionType(args []string, named bool) (ir.IrFunctionType, []string, error) {
	orig := args

	vars, args, err := ParseTupleArrow(args, named)
	if err != nil {
		return ir.IrFunctionType{}, orig, err
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
