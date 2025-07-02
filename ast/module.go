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
	// This module's ID.
	ModuleID ModuleID
	// This module's filename, e.g., from where it was read / parsed.
	Filename string
}

func (s Header) Format(f fmt.State, verb rune) {
	switch s.Case {
	case BaseFile:
		fmt.Fprint(f, "module ")
		s.ModuleID.Format(f, verb)
	case ImplementationFile:
		fmt.Fprint(f, "implements ")
		s.ModuleID.Format(f, verb)
	}
}

func (s Header) Is(c ModuleFileCase) bool {
	return s.Case == c
}

func NewBaseFileHeader(moduleID ModuleID) Header {
	return Header{Case: BaseFile, ModuleID: moduleID}
}

func NewImplementationFileHeader(moduleID ModuleID) Header {
	return Header{Case: ImplementationFile, ModuleID: moduleID}
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

	// TODO: This should be called for all module files, not just
	// implementation files.
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
	// TODO: Finish.
	//
	// if err := ValidateModuleID(module.Header.ModuleID); err != nil {
	// 	return err
	// }

	{
		// Validate imports.
		if !slices.IsSortedFunc(module.Imports.IDs, CompareModuleID) {
			module.AddError(
				module.Imports.Pos,
				"module %q has an 'imports' section that is not sorted", module.Header.ModuleID)
		}

		size := len(module.Imports.IDs)
		if imports := slices.CompactFunc(module.Imports.IDs, func(id1, id2 ModuleID) bool { return CompareModuleID(id1, id2) == 0 }); len(imports) != size {
			module.AddError(
				module.Imports.Pos,
				"module %q has an 'imports' section that contains duplicated imports", module.Header.ModuleID)
		}

		for _, id := range module.Imports.IDs {
			if err := ValidateModuleID(id); err != nil {
				module.AddError(id.Pos, err.Error())
			}
		}
	}

	if module.Header.Is(ImplementationFile) {
		if len(module.Impls.IDs) > 0 {
			module.AddError(
				module.Impls.Pos,
				"implementation file %q has an 'impls' section. The 'impls' section can only be used in module base files", module.Header.ModuleID)
		}
	}

	{
		// Validate impls.
		if !slices.IsSortedFunc(module.Impls.IDs, CompareID) {
			module.AddError(
				module.Impls.Pos,
				"file %q has an 'impls' section that is not sorted", module.Header.ModuleID)
		}

		size := len(module.Impls.IDs)
		if impls := slices.CompactFunc(module.Impls.IDs, func(id1, id2 ID) bool { return CompareID(id1, id2) == 0 }); len(impls) != size {
			module.AddError(
				module.Impls.Pos,
				"file %q has an 'impls' section that contains duplicated module implementation files", module.Header.ModuleID)
		}
	}
}
