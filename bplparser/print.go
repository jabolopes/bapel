package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) ParsePrintU() (Source, error) {
	if err := p.shiftToken("printU"); err != nil {
		return Source{}, err
	}

	return NewPrintSource(ir.Unsigned, p.Words()), nil
}

func (p *Parser) ParsePrintS() (Source, error) {
	if err := p.shiftToken("printS"); err != nil {
		return Source{}, err
	}

	return NewPrintSource(ir.Signed, p.Words()), nil
}
