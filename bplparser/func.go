package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseFunc() (string, []ir.IrDecl, []ir.IrDecl, error) {
	if err := p.shiftToken("func"); err != nil {
		return "", nil, nil, err
	}

	id, err := p.shiftID()
	if err != nil {
		return "", nil, nil, err
	}

	argTuple, retTuple, err := p.ParseTupleArrow(true /* named */)
	if err != nil {
		return "", nil, nil, err
	}

	if err = p.shiftToken("{"); err != nil {
		return "", nil, nil, err
	}

	if err := p.eol(); err != nil {
		return "", nil, nil, err
	}

	return id, argTuple, retTuple, err
}

func (p *Parser) ParseFunc() (r1 string, r2, r3 []ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		r1, r2, r3, err = p.parseFunc()
		return err
	})
	return
}
