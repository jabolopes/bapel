package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) ParseStruct() (string, ir.IrStructType, error) {
	if err := p.shiftToken("struct"); err != nil {
		return "", ir.IrStructType{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return "", ir.IrStructType{}, err
	}

	typ, err := p.ParseStructType(true /* named */)
	if err != nil {
		return "", ir.IrStructType{}, err
	}

	if err := p.eol(); err != nil {
		return "", ir.IrStructType{}, err
	}

	return id, typ, err
}
