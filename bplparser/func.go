package bplparser

func (p *Parser) parseFuncImpl() (Source, error) {
	if err := p.shiftToken("func"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	argTuple, retTuple, err := p.parseTupleArrow(true /* named */)
	if err != nil {
		return Source{}, err
	}

	if err = p.shiftToken("{"); err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewFunctionSource(id, argTuple, retTuple), err
}

func (p *Parser) parseFunc() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseFuncImpl()
		return err
	})
	return
}
