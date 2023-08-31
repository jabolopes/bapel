package bplparser

import (
	"fmt"
	"unicode"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseType(named bool) (ir.IrType, error) {
	if p.peekToken("(") {
		return p.ParseFunctionType(named)
	}

	if p.peekToken("{") {
		return p.ParseStructType(true /* named */)
	}

	if p.peekToken("[") {
		return p.ParseArrayType(named)
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

func (p *Parser) ParseType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseType(named)
		return err
	})
	return
}
