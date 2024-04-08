package bplparser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseForallTypeImpl(named bool) (ir.IrType, error) {
	if err := p.shiftLiteral("forall"); err != nil {
		return ir.IrType{}, err
	}

	if err := p.shiftLiteral("["); err != nil {
		return ir.IrType{}, err
	}

	var typeVars []string
	for {
		token, err := p.shiftID()
		if err != nil {
			return ir.IrType{}, err
		}

		if !strings.HasPrefix(token, "'") {
			return ir.IrType{}, fmt.Errorf("expected type variable; got %q", token)
		}

		typeVars = append(typeVars, strings.TrimPrefix(token, "'"))

		if p.shiftLiteral(",") == nil {
			continue
		}

		if err := p.shiftLiteral("]"); err != nil {
			return ir.IrType{}, err
		}

		break
	}

	subType, err := p.parseType(named)
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.NewForallType(typeVars, subType), nil
}

func (p *Parser) parseForallType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseForallTypeImpl(named)
		return err
	})
	return
}
