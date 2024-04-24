package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseBlock() (ir.IrTerm, error) {
	if err := p.shiftLiteral("{"); err != nil {
		return ir.IrTerm{}, err
	}

	if err := p.eol(); err != nil {
		return ir.IrTerm{}, err
	}

	var terms []ir.IrTerm
	for p.Scan() {
		if p.peek("}") {
			break
		}

		term, err := p.parseTerm()
		if err != nil {
			return ir.IrTerm{}, err
		}

		terms = append(terms, term)
	}

	if err := p.shiftLiteral("}"); err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewBlockTerm(terms), nil
}
