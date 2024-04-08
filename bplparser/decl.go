package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseDeclImpl(named bool) (ir.IrDecl, error) {
	isType := false
	if p.shiftLiteral("type") == nil {
		isType = true
	}

	id, err := p.shiftID()
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.shiftLiteral(":"); err != nil {
		return ir.IrDecl{}, err
	}

	typ, err := p.parseQuantifiedType(named)
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.eol(); err != nil {
		return ir.IrDecl{}, err
	}

	if isType {
		panic("not yet implemented")
	}

	return ir.NewTermDecl(id, typ), nil
}

func (p *Parser) parseDecl(named bool) (result ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseDeclImpl(named)
		return err
	})
	return
}
