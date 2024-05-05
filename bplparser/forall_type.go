package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseForallTypeImpl() (ir.IrType, error) {
	if err := p.shiftLiteral("forall"); err != nil {
		return ir.IrType{}, err
	}

	tvars, err := p.parseTypeAbstraction()
	if err != nil {
		return ir.IrType{}, err
	}

	subType, err := p.parseType()
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.ForallVars(tvars, subType), nil
}

func (p *Parser) parseForallType() (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseForallTypeImpl()
		return err
	})
	return
}
