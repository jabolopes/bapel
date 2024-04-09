package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

type DelimiterCase int

const (
	Parens DelimiterCase = iota
	Brackets
)

func (p *Parser) parseTupleImpl(named bool, delimiter DelimiterCase) ([]ir.IrDecl, error) {
	left := "("
	if delimiter == Brackets {
		left = "{"
	}

	right := ")"
	if delimiter == Brackets {
		right = "}"
	}

	if err := p.shiftLiteral(left); err != nil {
		return nil, err
	}

	if p.shiftLiteral(right) == nil {
		return nil, nil
	}

	var decls []ir.IrDecl
	for {
		var id string
		if named {
			var err error
			if id, err = p.shiftID(); err != nil {
				return nil, err
			}
		}

		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}

		decls = append(decls, ir.NewTermDecl(id, typ))

		if p.shiftLiteral(",") == nil {
			continue
		}

		if err = p.shiftLiteral(right); err != nil {
			return nil, err
		}

		break
	}

	return decls, nil
}

func (p *Parser) parseTuple(named bool, delimiter DelimiterCase) (result []ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTupleImpl(named, delimiter)
		return err
	})
	return
}
