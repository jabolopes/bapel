package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseStructImpl() (Source, error) {
	if err := p.shiftLiteral("struct"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	var tvars []ir.VarKind
	if p.peek("[") {
		var err error
		tvars, err = p.parseTypeAbstraction()
		if err != nil {
			return Source{}, err
		}
	}

	structType, err := p.parseStructType()
	if err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	lambdaType := ir.LambdaVars(tvars, structType)
	return NewTypeDefSource(ir.NewAliasDecl(id, lambdaType)), nil
}

func (p *Parser) parseStruct() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStructImpl()
		return err
	})
	return
}
