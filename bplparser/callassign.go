package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseAssignImpl() (ir.IrTerm, error) {
	isAssign := false
	var rets []ir.IrTerm
	for {
		if p.shiftLiteral("<-") == nil {
			isAssign = true
			break
		}

		token, err := p.shiftToken()
		if err != nil {
			return ir.IrTerm{}, err
		}

		rets = append(rets, ir.NewTokenTerm(token))
	}

	if !isAssign {
		return ir.IrTerm{}, fmt.Errorf("expected token '<-' in assignment term")
	}

	if len(rets) == 0 {
		return ir.IrTerm{}, fmt.Errorf("expected at least 1 return value before token '<-'")
	}

	callTerm, err := p.parseCall()
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewAssignTerm(callTerm, ir.NewTupleTerm(rets)), nil
}

func (p *Parser) parseAssign() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseAssignImpl()
		return err
	})
	return
}

// ParseCallAssign parses call and assignment.
//
// Note that a call is an assignment without the '<-' and without any return
// values.
func (p *Parser) parseCallAssign() (ir.IrTerm, error) {
	if typ, err := p.parseAssign(); err == nil {
		return typ, nil
	}

	return p.parseCall()
}
