package bplparser

import (
	"math"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseArrayTypeImpl() (ir.IrType, error) {
	if err := p.shiftLiteral("["); err != nil {
		return ir.IrType{}, err
	}

	typ, err := p.parseType()
	if err != nil {
		return ir.IrType{}, err
	}

	length, err := shiftInteger[int](p)
	if err != nil {
		length = math.MaxInt
	}

	if err := p.shiftLiteral("]"); err != nil {
		return ir.IrType{}, err
	}

	return ir.NewArrayType(typ, length), nil
}

func (p *Parser) parseArrayType() (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseArrayTypeImpl()
		return err
	})
	return
}
