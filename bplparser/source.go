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
	Export bool
	Decl   ir.IrDecl
}

func (s *typeDef) String() string {
	var b strings.Builder
	if s.Export {
		b.WriteString("export ")
	}
	b.WriteString(s.Decl.String())
	return b.String()
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
		return fmt.Sprintf("import %s", *s.Import)
	case TermSource:
		return s.Term.String()
	case TypeDefSource:
		return s.TypeDef.String()

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

func NewTypeDefSource(export bool, decl ir.IrDecl) Source {
	return Source{
		Case:    TypeDefSource,
		TypeDef: &typeDef{export, decl},
	}
}
