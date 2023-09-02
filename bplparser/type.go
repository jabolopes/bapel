package bplparser

import (
	"fmt"
	"unicode"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseTypeImpl(named bool) (ir.IrType, error) {
	if p.peek("(") {
		return p.parseFunctionType(named)
	}

	if p.peek("{") {
		return p.parseStructType(true /* named */)
	}

	if p.peek("[") {
		return p.parseArrayType(named)
	}

	token, err := p.shiftID()
	if err != nil {
		return ir.IrType{}, err
	}

	typ, err := ir.ParseIntType(token)
	if err == nil {
		return ir.NewIntType(typ), nil
	}

	var r rune
	for _, r = range token {
		break
	}

	if unicode.IsLetter(r) {
		return ir.NewIDType(token), nil
	}

	return ir.IrType{}, fmt.Errorf("expected type; got %q", token)
}

func (p *Parser) parseType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTypeImpl(named)
		return err
	})
	return
}
