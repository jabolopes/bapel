package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseTupleArrowImpl(named bool) ([]ir.IrDecl, []ir.IrDecl, error) {
	argTuple, err := p.parseTuple(named, Parens)
	if err != nil {
		return nil, nil, fmt.Errorf("in argument list: %v", err)
	}

	if err := p.shiftLiteral("->"); err != nil {
		return nil, nil, err
	}

	retTuple, err := p.parseTuple(named, Parens)
	if err != nil {
		return nil, nil, fmt.Errorf("in return list: %v", err)
	}

	return argTuple, retTuple, nil
}

func (p *Parser) parseTupleArrow(named bool) (r1, r2 []ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		r1, r2, err = p.parseTupleArrowImpl(named)
		return err
	})
	return
}
