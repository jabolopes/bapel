package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type DelimiterCase int

const (
	Parens = DelimiterCase(iota)
	Brackets
)

func (p *Parser) parseTuple(named bool, delimiter DelimiterCase) ([]ir.IrDecl, error) {
	left := "("
	if delimiter == Brackets {
		left = "{"
	}

	right := ")"
	if delimiter == Brackets {
		right = "}"
	}

	if err := p.shiftToken(left); err != nil {
		return nil, err
	}

	if err := p.shiftToken(right); err == nil {
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

		typ, err := p.ParseType(named)
		if err != nil {
			return nil, err
		}

		decls = append(decls, ir.NewVarDecl(id, typ))

		if err := p.shiftToken(","); err == nil {
			continue
		}

		if err = p.shiftToken(right); err != nil {
			return nil, err
		}

		break
	}

	return decls, nil
}

func (p *Parser) ParseTuple(named bool, delimiter DelimiterCase) (result []ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTuple(named, delimiter)
		return err
	})
	return
}

func (p *Parser) parseTupleArrow(named bool) ([]ir.IrDecl, []ir.IrDecl, error) {
	argTuple, err := p.ParseTuple(named, Parens)
	if err != nil {
		return nil, nil, fmt.Errorf("in argument list: %v", err)
	}

	if err := p.shiftToken("->"); err != nil {
		return nil, nil, err
	}

	retTuple, err := p.ParseTuple(named, Parens)
	if err != nil {
		return nil, nil, fmt.Errorf("in return list: %v", err)
	}

	return argTuple, retTuple, nil
}

func (p *Parser) ParseTupleArrow(named bool) (r1, r2 []ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		r1, r2, err = p.parseTupleArrow(named)
		return err
	})
	return
}
