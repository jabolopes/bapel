package ast

import (
	"fmt"
	"slices"

	"github.com/jabolopes/bapel/ir"
)

type ModuleFileCase int

const (
	BaseFile ModuleFileCase = iota
	ImplementationFile
)

type Header struct {
	Case ModuleFileCase
	// This module's name.
	Name string
	// Name of the base module this module belongs to. If this module is a
	// BaseModule, then this must be empty, otherwise it must be non-empty.
	BaseModuleName ID
}

func (s Header) Format(f fmt.State, verb rune) {
	if len(s.BaseModuleName.Value) == 0 {
		return
	}

	fmt.Fprint(f, "implements ")
	s.BaseModuleName.Format(f, verb)
}

func NewBaseFileHeader() Header {
	return Header{BaseFile, "", ID{}}
}

func NewImplementationFileHeader(baseModuleID ID) Header {
	return Header{ImplementationFile, "", baseModuleID}
}

type Imports struct {
	IDs []ModuleID
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

type Flags struct {
	IDs []ID
	Pos ir.Pos
}

func (s Flags) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintln(f, "flags {")
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
	// `impls` section of a TopModule. Must be empty for `ImplModule`.
	Impls  Impls
	Flags  Flags
	Body   []Source
	Errors []ir.Error
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

	if m.Header.Case == ImplementationFile {
		empty = false
		m.Header.Format(f, verb)
	}

	if len(m.Imports.IDs) > 0 {
		newline()
		m.Imports.Format(f, verb)
	}

	if len(m.Impls.IDs) > 0 {
		newline()
		m.Impls.Format(f, verb)
	}

	if len(m.Flags.IDs) > 0 {
		newline()
		m.Flags.Format(f, verb)
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

func (m Module) Valid() bool {
	return len(m.Errors) == 0
}

func (m *Module) AddError(pos ir.Pos, format string, args ...any) {
	m.Errors = append(m.Errors, ir.NewError(pos, fmt.Sprintf(format, args...)))
}

func NewImports(ids []ModuleID, pos ir.Pos) Imports {
	return Imports{ids, pos}
}

func NewImpls(ids []ID, pos ir.Pos) Impls {
	return Impls{ids, pos}
}

func NewFlags(ids []ID, pos ir.Pos) Flags {
	return Flags{ids, pos}
}

func ValidateModule(module *Module) {
	if len(module.Header.Name) == 0 {
		module.AddError(ir.Pos{}, "module missing module name (module name is empty)")
	}

	{
		// Validate imports.
		if !slices.IsSortedFunc(module.Imports.IDs, CompareModuleID) {
			module.AddError(
				module.Imports.Pos,
				"module %q has an 'imports' section that is not sorted", module.Header.Name)
		}

		size := len(module.Imports.IDs)
		if imports := slices.CompactFunc(module.Imports.IDs, func(id1, id2 ModuleID) bool { return CompareModuleID(id1, id2) == 0 }); len(imports) != size {
			module.AddError(
				module.Imports.Pos,
				"module %q has an 'imports' section that contains duplicated imports", module.Header.Name)
		}

		for _, id := range module.Imports.IDs {
			if err := ValidateModuleID(id); err != nil {
				module.AddError(id.Pos, err.Error())
			}
		}
	}

	switch module.Header.Case {
	case BaseFile:
		if len(module.Header.BaseModuleName.Value) != 0 {
			module.AddError(
				module.Header.BaseModuleName.Pos,
				"base file %q has an 'implements' line. The 'implements' line can only be used in module implementation files", module.Header.Name)
		}

	case ImplementationFile:
		if len(module.Header.BaseModuleName.Value) == 0 {
			module.AddError(
				module.Header.BaseModuleName.Pos,
				"implementation file %q is missing an 'implements' line at the top of the file. The 'implements' line must be present in module implementation files and it must name the base module that the implementation belongs to", module.Header.Name)
		}

		if len(module.Impls.IDs) > 0 {
			module.AddError(
				module.Impls.Pos,
				"implementation file %q has an 'impls' section. The 'impls' section can only be used in module base files", module.Header.Name)
		}
	}

	{
		// Validate impls.
		if !slices.IsSortedFunc(module.Impls.IDs, CompareID) {
			module.AddError(
				module.Impls.Pos,
				"file %q has an 'impls' section that is not sorted", module.Header.Name)
		}

		size := len(module.Impls.IDs)
		if impls := slices.CompactFunc(module.Impls.IDs, func(id1, id2 ID) bool { return CompareID(id1, id2) == 0 }); len(impls) != size {
			module.AddError(
				module.Impls.Pos,
				"file %q has an 'impls' section that contains duplicated module implementation files", module.Header.Name)
		}
	}
}
