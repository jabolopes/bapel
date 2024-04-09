package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseTupleTypeImpl(named bool) (ir.IrType, error) {
	tuple, err := p.parseTuple(named, Parens)
	if err != nil {
		return ir.IrType{}, err
	}

	types := make([]ir.IrType, len(tuple))
	for i, decl := range tuple {
		types[i] = decl.Type()
	}

	return ir.NewTupleType(types), nil
}

func (p *Parser) parseTupleType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTupleTypeImpl(named)
		return err
	})
	return
}
