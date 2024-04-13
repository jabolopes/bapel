package bplparser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	SectionSource SourceCase = iota
	EntitySource
	FunctionSource
	TermSource
	TypeDefSource
)

type section struct {
	ID    string
	Decls []ir.IrDecl
}

func (s section) String() string {
	var b strings.Builder
	b.WriteString(s.ID)
	b.WriteString(" {\n")
	for _, decl := range s.Decls {
		b.WriteString("  ")
		b.WriteString(decl.String())
		b.WriteString("\n")
	}
	b.WriteString("}")
	return b.String()
}

type typeDef struct {
	Type ir.IrType
}

type Source struct {
	Case     SourceCase
	Section  *section
	Entity   *ir.IrEntity
	Function *ir.IrFunction
	Term     *ir.IrTerm
	TypeDef  *typeDef
}

func (s Source) String() string {
	if s.Case == 0 && s.Section == nil {
		return ""
	}

	switch s.Case {
	case SectionSource:
		return s.Section.String()
	case EntitySource:
		return s.Entity.String()
	case FunctionSource:
		return s.Function.String()
	case TermSource:
		return s.Term.String()
	case TypeDefSource:
		return s.TypeDef.Type.String()

	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func NewSectionSource(id string, decls []ir.IrDecl) Source {
	return Source{
		Case:    SectionSource,
		Section: &section{id, decls},
	}
}

func NewEntitySource(entity ir.IrEntity) Source {
	s := Source{}
	s.Case = EntitySource
	s.Entity = &entity
	return s
}

func NewFunctionSource(function ir.IrFunction) Source {
	return Source{
		Case:     FunctionSource,
		Function: &function,
	}
}

func NewTermSource(term ir.IrTerm) Source {
	return Source{
		Case: TermSource,
		Term: &term,
	}
}

func NewTypeDefSource(typ ir.IrType) Source {
	return Source{
		Case:    TypeDefSource,
		TypeDef: &typeDef{typ},
	}
}

func (p *Parser) parseAnyImpl() (Source, error) {
	if source, err := p.parseSection(); err == nil {
		return source, nil
	}

	if p.peek("func") {
		return p.parseFunc()
	}

	if p.peek("struct") {
		return p.parseStruct()
	}

	if p.peek("entity") {
		return p.parseEntity()
	}

	return Source{}, fmt.Errorf("unknown syntax")
}

func (p *Parser) parseAny() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseAnyImpl()
		return err
	})
	return
}

func (p *Parser) ParseAny() (result Source, err error) {
	return p.parseAny()
}
