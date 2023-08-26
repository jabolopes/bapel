package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseStructType(named bool) (ir.IrStructType, error) {
	tuple, err := p.ParseTuple(named, Brackets)
	if err != nil {
		return ir.IrStructType{}, err
	}

	fields := make([]ir.StructField, len(tuple))
	for i, decl := range tuple {
		fields[i] = ir.StructField{decl.ID, decl.Type}
	}

	return ir.IrStructType{fields}, nil
}

func (p *Parser) ParseStructType(named bool) (result ir.IrStructType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStructType(named)
		return err
	})
	return
}
