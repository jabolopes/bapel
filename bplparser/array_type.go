package bplparser

import (
	"math"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseArrayType(named bool) (ir.IrArrayType, error) {
	if err := p.shiftToken("["); err != nil {
		return ir.IrArrayType{}, err
	}

	typ, args, err := p.ParseType(p.words, named)
	if err != nil {
		return ir.IrArrayType{}, err
	}
	p.words = args

	length, err := shiftInteger[int](p)
	if err != nil {
		length = math.MaxInt
	}

	if err := p.shiftToken("]"); err != nil {
		return ir.IrArrayType{}, err
	}

	return ir.IrArrayType{typ, length}, nil
}

func (p *Parser) ParseArrayType(named bool) (result ir.IrArrayType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseArrayType(named)
		return err
	})
	return
}
