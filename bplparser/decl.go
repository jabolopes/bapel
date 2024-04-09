package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseDeclImpl() (ir.IrDecl, error) {
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

	typ, err := p.parseQuantifiedType()
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

func (p *Parser) parseDecl() (result ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseDeclImpl()
		return err
	})
	return
}
