package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseDecl(named bool) (ir.IrDecl, error) {
	isType := false
	if err := p.shiftToken("type"); err == nil {
		isType = true
	}

	id, err := p.shiftID()
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.shiftToken(":"); err != nil {
		return ir.IrDecl{}, err
	}

	typ, err := p.ParseType(named)
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.eol(); err != nil {
		return ir.IrDecl{}, err
	}

	if isType {
		return ir.NewTypeDecl(id, typ), nil
	}

	return ir.NewTermDecl(id, typ), nil
}

func (p *Parser) ParseDecl(named bool) (result ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseDecl(named)
		return err
	})
	return
}
