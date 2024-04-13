package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) parseNumberTermImpl() (ir.IrTerm, error) {
	number, err := shiftInteger[int64](p)
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewTokenTerm(parser.NewNumberToken(number)), nil
}

func (p *Parser) parseNumberTerm() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseNumberTermImpl()
		return err
	})
	return
}
