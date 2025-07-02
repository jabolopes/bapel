package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type Package struct {
	ModuleID ModuleID
	Filename ID
	Pos      ir.Pos
}

func (s Package) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintf(f, "module %q in %q", s.ModuleID, s.Filename)
}

func NewPackage(moduleID ModuleID, filename ID, pos ir.Pos) Package {
	return Package{moduleID, filename, pos}
}

type Packages struct {
	Packages []Package
	Pos      ir.Pos
}

func (s Packages) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintln(f, "packages {")
	for _, pkg := range s.Packages {
		fmt.Fprint(f, "  ")
		pkg.Format(f, verb)
		fmt.Fprint(f, "\n")
	}
	fmt.Fprintln(f, "}")
}

func NewPackages(packages []Package, pos ir.Pos) Packages {
	return Packages{packages, pos}
}

type Workspace struct {
	Packages Packages
	Errors   []ir.Error
}

// TODO: Fix indentation: `s.Packages` should be further indented.
func (s Workspace) Format(f fmt.State, verb rune) {
	fmt.Fprintln(f, "workspace {")
	s.Packages.Format(f, verb)
	fmt.Fprint(f, "}")
}

func NewWorkspace(packages Packages) Workspace {
	return Workspace{packages, nil /* Errors */}
}

func ValidateWorkspace(workspace *Workspace) {
	// TODO: Finish.
}
