package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseStatementImpl() (Source, error) {
	term, err := p.parseCallAssign()
	if err != nil {
		return Source{}, err
	}

	return NewTermSource(ir.NewStatementTerm(term)), nil
}

func (p *Parser) parseStatement() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStatementImpl()
		return err
	})
	return
}
