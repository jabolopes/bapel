package bplparser

import (
	"math"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseArrayType(named bool) (ir.IrType, error) {
	if err := p.shiftToken("["); err != nil {
		return ir.IrType{}, err
	}

	typ, err := p.ParseType(named)
	if err != nil {
		return ir.IrType{}, err
	}

	length, err := shiftInteger[int](p)
	if err != nil {
		length = math.MaxInt
	}

	if err := p.shiftToken("]"); err != nil {
		return ir.IrType{}, err
	}

	return ir.NewArrayType(typ, length), nil
}

func (p *Parser) ParseArrayType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseArrayType(named)
		return err
	})
	return
}
