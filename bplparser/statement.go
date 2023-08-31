package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseStatement() (Source, error) {
	term, err := p.ParseCallAssign()
	if err != nil {
		return Source{}, err
	}

	return NewTermSource(ir.NewStatementTerm(term)), nil
}

func (p *Parser) ParseStatement() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStatement()
		return err
	})
	return
}
