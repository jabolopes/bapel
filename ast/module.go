package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

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
	Imports Imports
	Exports Exports
	Impls   Impls
	Body    []Source
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
