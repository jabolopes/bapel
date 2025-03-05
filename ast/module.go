package ast

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/jabolopes/bapel/ir"
)

type ModuleCase int

const (
	BaseModule ModuleCase = iota
	ImplModule
)

type Header struct {
	Case ModuleCase
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
	Exports Exports
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

func NewImports(ids []ID, pos ir.Pos) Imports {
	return Imports{ids, pos}
}

func NewExports(decls []ir.IrDecl, pos ir.Pos) Exports {
	return Exports{decls, pos}
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
		if !slices.IsSortedFunc(module.Imports.IDs, func(id1, id2 ID) int { return cmp.Compare(id1.Value, id2.Value) }) {
			module.AddError(
				module.Imports.Pos,
				"module %q has an 'imports' section that is not sorted", module.Header.Name)
		}

		size := len(module.Imports.IDs)
		if imports := slices.Compact(module.Imports.IDs); len(imports) != size {
			module.AddError(
				module.Imports.Pos,
				"module %q has an 'imports' section that contains duplicated imports", module.Header.Name)
		}
	}

	switch module.Header.Case {
	case BaseModule:
		if len(module.Header.BaseModuleName.Value) != 0 {
			module.AddError(
				module.Header.BaseModuleName.Pos,
				"base module %q has an 'implements' line. The 'implements' line can only be used in implementation modules", module.Header.Name)
		}

	case ImplModule:
		if len(module.Header.BaseModuleName.Value) == 0 {
			module.AddError(
				module.Header.BaseModuleName.Pos,
				"implementation module %q is missing an 'implements' line at the top of the file. The 'implements' line must be present in implementation modules and it must name the base module that the implementation belongs to", module.Header.Name)
		}

		if len(module.Impls.IDs) > 0 {
			module.AddError(
				module.Impls.Pos,
				"implementation module %q has an 'impls' section. The 'impls' section can only be used in base modules", module.Header.Name)
		}
	}

	{
		// Validate impls.
		if !slices.IsSortedFunc(module.Impls.IDs, func(id1, id2 ID) int { return cmp.Compare(id1.Value, id2.Value) }) {
			module.AddError(
				module.Impls.Pos,
				"module %q has an 'impls' section that is not sorted", module.Header.Name)
		}

		size := len(module.Impls.IDs)
		if impls := slices.Compact(module.Impls.IDs); len(impls) != size {
			module.AddError(
				module.Impls.Pos,
				"module %q has an 'impls' section that contains duplicated implementation modules", module.Header.Name)
		}
	}
}
