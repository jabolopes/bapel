package ast

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	ComponentSource SourceCase = iota
	ExportsSource
	FunctionSource
	ImplsSource
	ImportsSource
	TypeDefSource
)

type exportsSource struct {
	Decls []ir.IrDecl
}

func (s exportsSource) String() string {
	var b strings.Builder
	b.WriteString("exports {\n")
	for _, decl := range s.Decls {
		b.WriteString("  ")
		b.WriteString(decl.String())
		b.WriteString("\n")
	}
	b.WriteString("}")
	return b.String()
}

type implsSource struct {
	IDs []ID
}

func (s implsSource) String() string {
	var b strings.Builder
	b.WriteString("impls {\n")
	for _, id := range s.IDs {
		b.WriteString("  ")
		b.WriteString(id.Value)
		b.WriteString("\n")
	}
	b.WriteString("}")
	return b.String()
}

type importsSource struct {
	IDs []ID
}

func (s importsSource) String() string {
	var b strings.Builder
	b.WriteString("imports {\n")
	for _, id := range s.IDs {
		b.WriteString("  ")
		b.WriteString(id.Value)
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
	Component *ir.IrComponent
	Exports   *exportsSource
	Impls     *implsSource
	Imports   *importsSource
	Function  *ir.IrFunction
	Term      *ir.IrTerm
	TypeDef   *typeDefSource

	// Position in source file.
	Pos ir.Pos
}

func (s Source) String() string {
	if s.Case == 0 && s.Component == nil {
		return ""
	}

	switch s.Case {
	case ComponentSource:
		return s.Component.String()
	case ExportsSource:
		return s.Exports.String()
	case ImplsSource:
		return s.Impls.String()
	case ImportsSource:
		return s.Imports.String()
	case FunctionSource:
		return s.Function.String()
	case TypeDefSource:
		return s.TypeDef.String()
	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func (s Source) Format(f fmt.State, verb rune) {
	if s.Case == 0 && s.Component == nil {
		return
	}

	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	} else {
		fmt.Fprint(f, s.String())
		return
	}

	switch s.Case {
	case ExportsSource:
		fmt.Fprintln(f, "exports {")
		for _, decl := range s.Exports.Decls {
			fmt.Fprintf(f, "  %+s\n", decl)
		}
		fmt.Fprint(f, "}")

	case ImplsSource:
		fmt.Fprintln(f, "impls {")
		for _, id := range s.Impls.IDs {
			fmt.Fprintf(f, "  %+s\n", id)
		}
		fmt.Fprintf(f, "}")

	case ImportsSource:
		fmt.Fprintln(f, "imports {")
		for _, id := range s.Imports.IDs {
			fmt.Fprintf(f, "  %+s\n", id)
		}
		fmt.Fprint(f, "}")

	case FunctionSource:
		fmt.Fprint(f, s.Function.String())

	case TypeDefSource:
		fmt.Fprint(f, s.TypeDef.String())

	default:
		fmt.Fprint(f, s.String())
	}
}

func (s Source) Is(c SourceCase) bool {
	return s.Case == c
}

func NewComponentSource(component ir.IrComponent) Source {
	return Source{
		Case:      ComponentSource,
		Component: &component,
	}
}

func NewExportsSource(decls []ir.IrDecl) Source {
	return Source{
		Case:    ExportsSource,
		Exports: &exportsSource{decls},
	}
}

func NewFunctionSource(function ir.IrFunction) Source {
	return Source{
		Case:     FunctionSource,
		Function: &function,
	}
}

func NewImplsSource(ids []ID) Source {
	return Source{
		Case:  ImplsSource,
		Impls: &implsSource{ids},
	}
}

func NewImportsSource(ids []ID) Source {
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
