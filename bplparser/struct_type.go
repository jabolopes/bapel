package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseStructTypeImpl() (ir.IrType, error) {
	if err := p.shiftLiteral("{"); err != nil {
		return ir.IrType{}, err
	}

	if p.shiftLiteral("}") == nil {
		return ir.NewStructType(nil), nil
	}

	var fields []ir.StructField
	for {
		id, err := p.shiftID()
		if err != nil {
			return ir.IrType{}, err
		}

		typ, err := p.parseType()
		if err != nil {
			return ir.IrType{}, err
		}

		fields = append(fields, ir.StructField{id, typ})

		if p.shiftLiteral(",") == nil {
			continue
		}

		break
	}

	if err := p.shiftLiteral("}"); err != nil {
		return ir.IrType{}, err
	}

	return ir.NewStructType(fields), nil
}

func (p *Parser) parseStructType() (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStructTypeImpl()
		return err
	})
	return
}
