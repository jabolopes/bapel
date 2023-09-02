package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseStructTypeImpl(named bool) (ir.IrType, error) {
	tuple, err := p.parseTuple(named, Brackets)
	if err != nil {
		return ir.IrType{}, err
	}

	fields := make([]ir.StructField, len(tuple))
	for i, decl := range tuple {
		fields[i] = ir.StructField{decl.ID, decl.Type}
	}

	return ir.NewStructType(ir.IrStructType{fields}), nil
}

func (p *Parser) parseStructType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStructTypeImpl(named)
		return err
	})
	return
}
