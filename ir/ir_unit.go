package ir

import (
	"cmp"
	"fmt"

	"golang.org/x/exp/slices"
)

type IrUnitCase int

const (
	BaseUnit IrUnitCase = iota
	ImplUnit
)

func (t IrUnitCase) Format(f fmt.State, verb rune) {
	switch t {
	case BaseUnit:
		fmt.Fprint(f, "base")
	case ImplUnit:
		fmt.Fprint(f, "implementation")
	default:
		panic(fmt.Errorf("unhandled %T %d", t, t))
	}
}

// TODO: Replace string with ModuleID to retain Pos. Requires moving
// ModuleID to the ir package.
type IrImport struct {
	ModuleID string
}

func NewImport(moduleID string) IrImport {
	return IrImport{moduleID}
}

func CompareImport(i1, i2 IrImport) int {
	return cmp.Compare(i1.ModuleID, i2.ModuleID)
}

func EqualsImport(i1, i2 IrImport) bool {
	return CompareImport(i1, i2) == 0
}

func CleanImports(imports []IrImport) []IrImport {
	slices.SortFunc(imports, CompareImport)
	return slices.CompactFunc(imports, EqualsImport)
}

// TODO: Replace string with Filename to retain Pos. Requires moving
// Filename to the ir package.
type IrImpl struct {
	RelativeFilename string
}

func NewImpl(relativeFilename string) IrImpl {
	return IrImpl{relativeFilename}
}

type IrUnit struct {
	Case        IrUnitCase
	ModuleID    string
	Filename    string
	Imports     []IrImport
	Impls       []IrImpl
	ImportDecls []IrDecl
	ImplDecls   []IrDecl
	Decls       []IrDecl
	Functions   []IrFunction
}

func (t IrUnit) Format(f fmt.State, verb rune) {
	{
		fmt.Fprintln(f, "unit {")
		fmt.Fprintf(f, "  module %s\n", t.ModuleID)
		fmt.Fprintf(f, "  filename %q\n", t.Filename)

		fmt.Fprint(f, "  ")
		fmt.Fprintf(f, fmt.FormatString(f, 's'), t.Case)
		fmt.Fprintln(f)

		fmt.Fprintln(f, "}")
	}

	if len(t.Imports) > 0 {
		fmt.Fprintln(f, "\nimports {")
		for _, imp := range t.Imports {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), imp.ModuleID)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}\n")
	}

	if len(t.Impls) > 0 {
		fmt.Fprintln(f, "\nimpls {")
		for _, impl := range t.Impls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 'q'), impl.RelativeFilename)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}\n")
	}

	if len(t.ImportDecls) > 0 {
		fmt.Fprintln(f, "\nimportDecls {")
		for _, decl := range t.ImportDecls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), decl)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}\n")
	}

	if len(t.ImplDecls) > 0 {
		fmt.Fprintln(f, "\nimplDecls {")
		for _, decl := range t.ImplDecls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), decl)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}\n")
	}

	if len(t.Decls) > 0 {
		fmt.Fprintln(f, "\ndecls {")
		for _, decl := range t.Decls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), decl)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}\n")
	}

	if len(t.Functions) > 0 {
		fmt.Fprintln(f)
		fmt.Fprintf(f, fmt.FormatString(f, 's'), t.Functions[0])
		for _, function := range t.Functions[1:] {
			fmt.Fprint(f, "\n\n")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), function)
		}
	}
}
