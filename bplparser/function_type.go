package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseFunctionTypeImpl() (ir.IrType, error) {
	arg, err := p.parseSimpleType()
	if err != nil {
		return ir.IrType{}, err
	}

	if err := p.shiftLiteral("->"); err != nil {
		return ir.IrType{}, err
	}

	ret, err := p.parseType()
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.NewFunctionType(arg, ret), nil
}

func (p *Parser) parseFunctionType() (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseFunctionTypeImpl()
		return err
	})
	return
}
