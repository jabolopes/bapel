package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type ModuleCase int

const (
	TopModule ModuleCase = iota
	ImplModule
)

type Header struct {
	Case ModuleCase
	// This module's name.
	Name string
	// If this Header belongs to a TopModule, this is always empty. Otherwise,
	// this must be the name of the TopModule that this implements.
	TopName ID
}

func (s Header) Format(f fmt.State, verb rune) {
	if len(s.TopName.Value) == 0 {
		return
	}

	fmt.Fprint(f, "implements ")
	s.TopName.Format(f, verb)
	fmt.Fprintln(f)
}

type Imports struct {
	IDs []ID
	Pos ir.Pos
}

func (s Imports) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintln(f, "imports {")
	for _, id := range s.IDs {
		fmt.Fprint(f, "  ")
		id.Format(f, verb)
		fmt.Fprint(f, "\n")
	}
	fmt.Fprint(f, "}")
}

type Exports struct {
	Decls []ir.IrDecl
	Pos   ir.Pos
}

func (s Exports) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintln(f, "exports {")
	for _, decl := range s.Decls {
		fmt.Fprint(f, "  ")
		decl.Format(f, verb)
		fmt.Fprint(f, "\n")
	}
	fmt.Fprint(f, "}")
}

type Impls struct {
	IDs []ID
	Pos ir.Pos
}

func (s Impls) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintln(f, "impls {")
	for _, id := range s.IDs {
		fmt.Fprint(f, "  ")
		id.Format(f, verb)
		fmt.Fprint(f, "\n")
	}
	fmt.Fprint(f, "}")
}

type Module struct {
	Header  Header
	Imports Imports
	Exports Exports
	Impls   Impls
	Body    []Source
}

func (m Module) Format(f fmt.State, verb rune) {
	empty := true

	newline := func() {
		if !empty {
			fmt.Fprintln(f)
			fmt.Fprintln(f)
		}
		empty = false
	}

	if m.Header.Case == ImplModule {
		empty = false
		m.Header.Format(f, verb)
	}

	if len(m.Imports.IDs) > 0 {
		newline()
		m.Imports.Format(f, verb)
	}

	if len(m.Exports.Decls) > 0 {
		newline()
		m.Exports.Format(f, verb)
	}

	if len(m.Impls.IDs) > 0 {
		newline()
		m.Impls.Format(f, verb)
	}

	if len(m.Body) > 0 {
		newline()
		m.Body[0].Format(f, verb)
		for _, source := range m.Body[1:] {
			fmt.Fprintln(f)
			source.Format(f, verb)
		}
	}
}

func NewImports(ids []ID, pos ir.Pos) Imports {
	return Imports{ids, pos}
}

func NewExports(decls []ir.IrDecl, pos ir.Pos) Exports {
	return Exports{decls, pos}
}

func NewImpls(ids []ID, pos ir.Pos) Impls {
	return Impls{ids, pos}
}
