package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseLet() (ir.IrDecl, error) {
	if err := p.shiftToken("let"); err != nil {
		return ir.IrDecl{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return ir.IrDecl{}, err
	}

	typ, err := p.ParseType(false /* named */)
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.eol(); err != nil {
		return ir.IrDecl{}, err
	}

	return ir.NewTermDecl(id, typ), nil
}

func (p *Parser) ParseLet() (result ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseLet()
		return err
	})
	return
}
