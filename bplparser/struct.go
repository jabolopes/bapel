package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) ParseStruct() (ir.IrDecl, error) {
	if err := p.shiftToken("struct"); err != nil {
		return ir.IrDecl{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return ir.IrDecl{}, err
	}

	typ, err := p.ParseStructType(true /* named */)
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.eol(); err != nil {
		return ir.IrDecl{}, err
	}

	return ir.NewTypeDecl(id, typ), err
}
