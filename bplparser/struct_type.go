package bplparser

import "github.com/jabolopes/bapel/ir"

func ParseStructType(args []string, named bool) (ir.IrStructType, []string, error) {
	vars, args, err := ParseTuple(args, ir.ArgVar, named, Brackets)
	if err != nil {
		return ir.IrStructType{}, args, err
	}

	fields := make([]ir.StructField, len(vars))
	for i, irvar := range vars {
		fields[i] = ir.StructField{irvar.Id, irvar.Type}
	}

	return ir.IrStructType{fields}, args, nil
}
