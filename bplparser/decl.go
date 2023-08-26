package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) ParseDecl(args []string, named bool) (ir.IrDecl, []string, error) {
	orig := args

	id, args, err := parser.ShiftID(args)
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	args, err = parser.ShiftToken(args, ":")
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	p.words = args
	typ, err := p.ParseType(named)
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	if err := p.eol(); err != nil {
		return ir.IrDecl{}, orig, err
	}

	// TODO: Finish. The following is technically wrong.
	if typ.Is(ir.StructType) {
		return ir.NewTypeDecl(id, typ), p.words, nil
	}

	if typ.Is(ir.FunType) {
		return ir.NewConstantDecl(id, typ), p.words, nil
	}

	return ir.NewVarDecl(id, typ), p.words, nil
}
