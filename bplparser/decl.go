package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseTypeDecl() (ir.IrDecl, error) {
	if err := p.shiftLiteral("type"); err != nil {
		return ir.IrDecl{}, err
	}

	typ, err := p.parseQuantifiedType()
	if err != nil {
		return ir.IrDecl{}, err
	}

	if p.shiftLiteral("=") == nil {
		typ2, err := p.parseQuantifiedType()
		if err != nil {
			return ir.IrDecl{}, err
		}

		typ = ir.NewAliasType(typ, typ2)
	}

	if err := p.eol(); err != nil {
		return ir.IrDecl{}, err
	}

	return ir.NewTypeDecl(typ), nil
}

func (p *Parser) parseTermDecl() (ir.IrDecl, error) {
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

	return ir.NewTermDecl(id, typ), nil
}

func (p *Parser) parseDeclImpl() (ir.IrDecl, error) {
	if p.peek("type") {
		return p.parseTypeDecl()
	}

	return p.parseTermDecl()
}

func (p *Parser) parseDecl() (result ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseDeclImpl()
		return err
	})
	return
}
