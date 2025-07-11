package ir

import "fmt"

type IrUnitCase int

const (
	BaseUnit IrUnitCase = iota
	ImplUnit
)

// TODO: Replace string with ModuleID to retain Pos. Requires moving
// ModuleID to the ir package.
type IrImport struct {
	ModuleID string
}

func NewImport(moduleID string) IrImport {
	return IrImport{moduleID}
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
	fmt.Fprintf(f, "module %s\n", t.ModuleID)

	fmt.Fprintf(f, "\nfilename %q\n", t.Filename)

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
		for _, function := range t.Functions {
			fmt.Fprintf(f, "\n\n")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), function)
		}
	}
}
