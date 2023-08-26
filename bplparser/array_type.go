package bplparser

import (
	"math"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) ParseArrayType(args []string, named bool) (ir.IrArrayType, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "[")
	if err != nil {
		return ir.IrArrayType{}, orig, err
	}

	typ, args, err := p.ParseType(args, named)
	if err != nil {
		return ir.IrArrayType{}, orig, err
	}

	length, args, err := parser.ShiftNumber[int](args)
	if err != nil {
		length = math.MaxInt
	}

	args, err = parser.ShiftToken(args, "]")
	if err != nil {
		return ir.IrArrayType{}, orig, err
	}

	return ir.IrArrayType{typ, length}, args, nil
}
