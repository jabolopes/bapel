package bplparser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	SectionSource SourceCase = iota
	ComponentSource
	FunctionSource
	ImportSource
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
	Decl ir.IrDecl
}

type Source struct {
	Case      SourceCase
	Section   *section
	Component *ir.IrComponent
	Function  *ir.IrFunction
	Import    *string
	Term      *ir.IrTerm
	TypeDef   *typeDef
}

func (s Source) String() string {
	if s.Case == 0 && s.Section == nil {
		return ""
	}

	switch s.Case {
	case SectionSource:
		return s.Section.String()
	case ComponentSource:
		return s.Component.String()
	case FunctionSource:
		return s.Function.String()
	case ImportSource:
		return *s.Import
	case TermSource:
		return s.Term.String()
	case TypeDefSource:
		return s.TypeDef.Decl.String()

	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func (s Source) Is(c SourceCase) bool {
	return s.Case == c
}

func NewSectionSource(id string, decls []ir.IrDecl) Source {
	return Source{
		Case:    SectionSource,
		Section: &section{id, decls},
	}
}

func NewComponentSource(component ir.IrComponent) Source {
	return Source{
		Case:      ComponentSource,
		Component: &component,
	}
}

func NewFunctionSource(function ir.IrFunction) Source {
	return Source{
		Case:     FunctionSource,
		Function: &function,
	}
}

func NewImportSource(id string) Source {
	return Source{
		Case:   ImportSource,
		Import: &id,
	}
}

func NewTermSource(term ir.IrTerm) Source {
	return Source{
		Case: TermSource,
		Term: &term,
	}
}

func NewTypeDefSource(decl ir.IrDecl) Source {
	return Source{
		Case:    TypeDefSource,
		TypeDef: &typeDef{decl},
	}
}

func (p *Parser) parseAnyImpl() (Source, error) {
	if p.peek("import") {
		return p.parseImport()
	}

	if source, err := p.parseSection(); err == nil {
		return source, nil
	}

	if p.peek("func") {
		return p.parseFunc()
	}

	if p.peek("struct") {
		return p.parseStruct()
	}

	if p.peek("component") {
		return p.parseComponent()
	}

	if p.peek("export") {
		// TODO: Generalize to other syntaxes.
		return p.parseFunc()
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
