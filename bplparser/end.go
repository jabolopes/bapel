package bplparser

func (p *Parser) parseEnd() error {
	if err := p.shiftToken("}"); err != nil {
		return err
	}

	return p.eol()
}

func (p *Parser) ParseEnd() error {
	return p.withCheckpoint(p.parseEnd)
}
