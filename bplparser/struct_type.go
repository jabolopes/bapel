package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) ParseStructType(args []string, named bool) (ir.IrStructType, []string, error) {
	tuple, args, err := p.ParseTuple(args, named, Brackets)
	if err != nil {
		return ir.IrStructType{}, args, err
	}

	fields := make([]ir.StructField, len(tuple))
	for i, decl := range tuple {
		fields[i] = ir.StructField{decl.ID, decl.Type}
	}

	return ir.IrStructType{fields}, args, nil
}
