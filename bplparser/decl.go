package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseDecl(named bool) (ir.IrDecl, error) {
	id, err := p.shiftID()
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.shiftToken(":"); err != nil {
		return ir.IrDecl{}, err
	}

	typ, err := p.ParseType(named)
	if err != nil {
		return ir.IrDecl{}, err
	}

	if err := p.eol(); err != nil {
		return ir.IrDecl{}, err
	}

	// TODO: Finish. The following is technically wrong.
	if typ.Is(ir.StructType) {
		return ir.NewTypeDecl(id, typ), nil
	}

	if typ.Is(ir.FunType) {
		return ir.NewConstantDecl(id, typ), nil
	}

	return ir.NewVarDecl(id, typ), nil
}

func (p *Parser) ParseDecl(named bool) (result ir.IrDecl, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseDecl(named)
		return err
	})
	return
}
