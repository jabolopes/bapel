package bplparser

import (
	"fmt"
	"strings"
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

	var r rune
	for _, r = range token {
		break
	}

	if r == '\'' {
		return ir.NewVarType(strings.TrimPrefix(token, "'")), nil
	}

	if unicode.IsLetter(r) {
		return ir.NewNameType(token), nil
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

func (p *Parser) parseQuantifiedType(named bool) (result ir.IrType, err error) {
	typ, err := p.parseType(named)
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.QuantifyType(typ), nil
}
