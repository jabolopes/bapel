package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseBlock() (ir.IrTerm, error) {
	var terms []ir.IrTerm
	for p.Scan() {
		if p.peek("}") {
			break
		}

		statement, err := p.parseStatement()
		if err != nil {
			return ir.IrTerm{}, err
		}

		terms = append(terms, statement)
	}

	if err := p.shiftLiteral("}"); err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewBlockTerm(terms), nil
}

func (p *Parser) parseIfImpl() (ir.IrTerm, error) {
	if err := p.shiftLiteral("if"); err != nil {
		return ir.IrTerm{}, err
	}

	negate := false
	if p.shiftLiteral("not") == nil {
		negate = true
	}

	if err := p.shiftLiteralEnd("{"); err != nil {
		return ir.IrTerm{}, err
	}

	condition, err := p.parseCall()
	if err != nil {
		return ir.IrTerm{}, err
	}

	then, err := p.parseBlock()
	if err != nil {
		return ir.IrTerm{}, err
	}

	var elseTerm *ir.IrTerm
	if p.shiftLiteral("else") == nil {
		if err := p.shiftLiteral("{"); err != nil {
			return ir.IrTerm{}, err
		}

		term, err := p.parseBlock()
		if err != nil {
			return ir.IrTerm{}, err
		}

		elseTerm = &term
	}

	return ir.NewIfTerm(negate, condition, then, elseTerm), nil
}

func (p *Parser) parseIf() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseIfImpl()
		return err
	})
	return
}
