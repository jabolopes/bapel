package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) ParseLet(args []string) (ir.IrDecl, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "let")
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	id, args, err := parser.ShiftID(args)
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	p.words = args
	typ, err := p.ParseType(false /* named */)
	if err != nil {
		return ir.IrDecl{}, orig, err
	}
	args = p.words

	if err := parser.EOL(args); err != nil {
		return ir.IrDecl{}, orig, err
	}

	// TODO: Could also be constant instead of var, or have 2 syntaxes for mutable
	// and immutable identifiers.
	return ir.NewVarDecl(id, typ), nil, nil
}
