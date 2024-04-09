package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseStructTypeImpl(named bool) (ir.IrType, error) {
	if err := p.shiftLiteral("{"); err != nil {
		return ir.IrType{}, err
	}

	if p.shiftLiteral("}") == nil {
		return ir.NewStructType(nil), nil
	}

	var fields []ir.StructField
	for {
		var id string
		if named {
			var err error
			if id, err = p.shiftID(); err != nil {
				return ir.IrType{}, err
			}
		}

		typ, err := p.parseType(named)
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

func (p *Parser) parseStructType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStructTypeImpl(named)
		return err
	})
	return
}
