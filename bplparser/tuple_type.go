package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseTupleTypeImpl(named bool) (ir.IrType, error) {
	if err := p.shiftLiteral("("); err != nil {
		return ir.IrType{}, err
	}

	if p.shiftLiteral(")") == nil {
		return ir.NewTupleType(nil), nil
	}

	var types []ir.IrType
	{
		typ, err := p.parseType(named)
		if err != nil {
			return ir.IrType{}, err
		}

		types = append(types, typ)

		if err := p.shiftLiteral(","); err != nil {
			return ir.IrType{}, err
		}
	}

	for {
		typ, err := p.parseType(named)
		if err != nil {
			return ir.IrType{}, err
		}

		types = append(types, typ)

		if p.shiftLiteral(",") == nil {
			continue
		}

		break
	}

	if err := p.shiftLiteral(")"); err != nil {
		return ir.IrType{}, err
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
