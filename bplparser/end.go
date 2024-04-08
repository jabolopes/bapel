package bplparser

func (p *Parser) parseEndImpl() error {
	if err := p.shiftLiteral("}"); err != nil {
		return err
	}

	return p.eol()
}

func (p *Parser) parseEnd() error {
	return p.withCheckpoint(p.parseEndImpl)
}
