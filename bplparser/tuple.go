package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

type DelimiterCase int

const (
	Parens = DelimiterCase(iota)
	Brackets
)

func (p *Parser) ParseTuple(args []string, named bool, delimiter DelimiterCase) ([]ir.IrDecl, []string, error) {
	orig := args

	args, remainder := parser.ShiftBalancedParens(args)

	left := "("
	if delimiter == Brackets {
		left = "{"
	}

	right := ")"
	if delimiter == Brackets {
		right = "}"
	}

	args, err := parser.ShiftToken(args, left)
	if err != nil {
		return nil, orig, err
	}

	if _, err := parser.ShiftToken(args, right); err == nil {
		return nil, remainder, nil
	}

	var decls []ir.IrDecl
	for {
		var id string
		if named {
			id, args, err = parser.ShiftID(args)
			if err != nil {
				return nil, orig, err
			}
		}

		var typ ir.IrType
		typ, args, err = p.ParseType(args, named)
		if err != nil {
			return nil, orig, err
		}

		decls = append(decls, ir.NewVarDecl(id, typ))

		if args, err = parser.ShiftToken(args, ","); err == nil {
			continue
		}

		args, err = parser.ShiftToken(args, right)
		if err != nil {
			return nil, orig, err
		}

		break
	}

	return decls, remainder, nil
}

func (p *Parser) ParseTupleArrow(args []string, named bool) ([]ir.IrDecl, []ir.IrDecl, []string, error) {
	orig := args

	argTuple, args, err := p.ParseTuple(args, named, Parens)
	if err != nil {
		return nil, nil, orig, fmt.Errorf("in argument list: %v", err)
	}

	args, err = parser.ShiftToken(args, "->")
	if err != nil {
		return nil, nil, orig, err
	}

	retTuple, args, err := p.ParseTuple(args, named, Parens)
	if err != nil {
		return nil, nil, orig, fmt.Errorf("in return list: %v", err)
	}

	return argTuple, retTuple, args, nil
}
