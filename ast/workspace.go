package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type PackageCase int

const (
	ModulePackage PackageCase = iota
	PrefixPackage
)

type modulePackage struct {
	ModuleID ModuleID
}

func (s *modulePackage) Format(f fmt.State, verb rune) {
	fmt.Fprint(f, "module ")
	s.ModuleID.Format(f, verb)
}

type prefixPackage struct {
	Prefix ModuleID
}

func (s *prefixPackage) Format(f fmt.State, verb rune) {
	fmt.Fprint(f, "prefix ")
	s.Prefix.Format(f, verb)
}

type Package struct {
	Case     PackageCase
	Module   *modulePackage
	Prefix   *prefixPackage
	Filename Filename
	Pos      ir.Pos
}

func (s Package) Is(c PackageCase) bool {
	return s.Case == c
}

func (s Package) Format(f fmt.State, verb rune) {
	if s.Case == 0 && s.Module == nil {
		return
	}

	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	switch s.Case {
	case ModulePackage:
		s.Module.Format(f, verb)
	case PrefixPackage:
		s.Prefix.Format(f, verb)
	default:
		panic(fmt.Sprintf("unhandled %T %d", s.Case, s.Case))
	}

	fmt.Fprintf(f, " in %q", s.Filename)
}

func NewModulePackage(moduleID ModuleID, filename Filename, pos ir.Pos) Package {
	return Package{
		Case:     ModulePackage,
		Module:   &modulePackage{moduleID},
		Filename: filename,
		Pos:      pos,
	}
}

func NewPrefixPackage(prefix ModuleID, filename Filename, pos ir.Pos) Package {
	return Package{
		Case:     PrefixPackage,
		Prefix:   &prefixPackage{prefix},
		Filename: filename,
		Pos:      pos,
	}
}

func ValidatePackage(pkg Package) ir.Validation {
	var validation ir.Validation

	switch pkg.Case {
	case ModulePackage:
		c := pkg.Module

		if err := ValidateModuleID(c.ModuleID); err != nil {
			validation.AddErr(c.ModuleID.Pos, err)
		}
	case PrefixPackage:
		c := pkg.Prefix

		if err := ValidateModuleID(c.Prefix); err != nil {
			validation.AddErr(c.Prefix.Pos, err)
		}
	default:
		panic(fmt.Sprintf("unhandled %T %d", pkg.Case, pkg.Case))
	}

	if err := ValidateFilename(pkg.Filename); err != nil {
		validation.AddErr(pkg.Filename.Pos, err)
	}

	return validation
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

func ValidatePackages(packages Packages) ir.Validation {
	var validation ir.Validation

	for _, pkg := range packages.Packages {
		validation.Join(ValidatePackage(pkg))
	}

	return validation
}

type Workspace struct {
	Packages Packages
}

// TODO: Fix indentation: `s.Packages` should be further indented.
func (s Workspace) Format(f fmt.State, verb rune) {
	fmt.Fprintln(f, "workspace {")
	s.Packages.Format(f, verb)
	fmt.Fprint(f, "}")
}

func NewWorkspace(packages Packages) Workspace {
	return Workspace{packages}
}

func ValidateWorkspace(workspace Workspace) ir.Validation {
	return ValidatePackages(workspace.Packages)
}
