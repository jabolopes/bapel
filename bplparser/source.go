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
	ImplsSource
	ImportsSource
	TypeDefSource
)

type implsSource struct {
	IDs []string
}

func (s implsSource) String() string {
	var b strings.Builder
	b.WriteString("impls {\n")
	for _, id := range s.IDs {
		b.WriteString("  ")
		b.WriteString(id)
		b.WriteString("\n")
	}
	b.WriteString("}")
	return b.String()
}

type importsSource struct {
	IDs []string
}

func (s importsSource) String() string {
	var b strings.Builder
	b.WriteString("imports {\n")
	for _, id := range s.IDs {
		b.WriteString("  ")
		b.WriteString(id)
		b.WriteString("\n")
	}
	b.WriteString("}")
	return b.String()
}

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

type typeDefSource struct {
	Export bool
	Decl   ir.IrDecl
}

func (s *typeDefSource) String() string {
	var b strings.Builder
	if s.Export {
		b.WriteString("export ")
	}
	b.WriteString(s.Decl.String())
	return b.String()
}

type Source struct {
	Case      SourceCase
	Impls     *implsSource
	Imports   *importsSource
	Section   *section
	Component *ir.IrComponent
	Function  *ir.IrFunction
	Term      *ir.IrTerm
	TypeDef   *typeDefSource

	// Position in source file.
	Pos ir.Pos
}

func (s Source) String() string {
	if s.Case == 0 && s.Section == nil {
		return ""
	}

	switch s.Case {
	case ImplsSource:
		return s.Impls.String()
	case ImportsSource:
		return s.Imports.String()
	case SectionSource:
		return s.Section.String()
	case ComponentSource:
		return s.Component.String()
	case FunctionSource:
		return s.Function.String()
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

func NewImplsSource(ids []string) Source {
	return Source{
		Case:  ImplsSource,
		Impls: &implsSource{ids},
	}
}

func NewImportsSource(ids []string) Source {
	return Source{
		Case:    ImportsSource,
		Imports: &importsSource{ids},
	}
}

func NewTypeDefSource(export bool, decl ir.IrDecl) Source {
	return Source{
		Case:    TypeDefSource,
		TypeDef: &typeDefSource{export, decl},
	}
}
