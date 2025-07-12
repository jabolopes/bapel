package ast

import (
	"fmt"
	"slices"

	"github.com/jabolopes/bapel/ir"
)

type SourceFileCase int

const (
	BaseSourceFile SourceFileCase = iota
	ImplSourceFile
)

type SourceFileHeader struct {
	Case SourceFileCase
	// The module ID the current source file belongs to.
	ModuleID ir.ModuleID
	// The current source file's filename, e.g., from where it was read / parsed.
	//
	// TODO: Replace string with ir.Filename.
	Filename string
}

func (s SourceFileHeader) Format(f fmt.State, verb rune) {
	switch s.Case {
	case BaseSourceFile:
		fmt.Fprint(f, "module ")
		s.ModuleID.Format(f, verb)
	case ImplSourceFile:
		fmt.Fprint(f, "implements ")
		s.ModuleID.Format(f, verb)
	}
}

func (s SourceFileHeader) Is(c SourceFileCase) bool {
	return s.Case == c
}

func NewBaseSourceFileHeader(moduleID ir.ModuleID) SourceFileHeader {
	return SourceFileHeader{Case: BaseSourceFile, ModuleID: moduleID}
}

func NewImplSourceFileHeader(moduleID ir.ModuleID) SourceFileHeader {
	return SourceFileHeader{Case: ImplSourceFile, ModuleID: moduleID}
}

type Imports struct {
	IDs []ir.ModuleID
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
	Filenames []ir.Filename
	Pos       ir.Pos
}

func (s Impls) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintln(f, "impls {")
	for _, id := range s.Filenames {
		fmt.Fprint(f, "  ")
		fmt.Fprintf(f, fmt.FormatString(f, 'q'), id)
		fmt.Fprint(f, "\n")
	}
	fmt.Fprint(f, "}")
}

type Flags struct {
	Filenames []ir.Filename
	Pos       ir.Pos
}

func (s Flags) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintln(f, "flags {")
	for _, id := range s.Filenames {
		fmt.Fprint(f, "  ")
		id.Format(f, verb)
		fmt.Fprint(f, "\n")
	}
	fmt.Fprint(f, "}")
}

type SourceFile struct {
	Header  SourceFileHeader
	Imports Imports
	Impls   Impls
	Flags   Flags
	Body    []Source
}

func (m SourceFile) Format(f fmt.State, verb rune) {
	newline := func() {
		fmt.Fprintln(f)
		fmt.Fprintln(f)
	}

	m.Header.Format(f, verb)

	if len(m.Imports.IDs) > 0 {
		newline()
		m.Imports.Format(f, verb)
	}

	if len(m.Impls.Filenames) > 0 {
		newline()
		m.Impls.Format(f, verb)
	}

	if len(m.Flags.Filenames) > 0 {
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

func NewImports(ids []ir.ModuleID, pos ir.Pos) Imports {
	return Imports{ids, pos}
}

func NewImpls(filenames []ir.Filename, pos ir.Pos) Impls {
	return Impls{filenames, pos}
}

func NewFlags(filenames []ir.Filename, pos ir.Pos) Flags {
	return Flags{filenames, pos}
}

func ValidateSourceFile(sourceFile *SourceFile) ir.Validation {
	var validation ir.Validation

	if err := ir.ValidateModuleID(sourceFile.Header.ModuleID); err != nil {
		validation.AddErr(sourceFile.Header.ModuleID.Pos, err)
	}

	{
		// Validate imports.
		if !slices.IsSortedFunc(sourceFile.Imports.IDs, ir.CompareModuleID) {
			validation.AddErrorf(
				sourceFile.Imports.Pos,
				"source file %q has an 'imports' section that is not sorted", sourceFile.Header.ModuleID)
		}

		size := len(sourceFile.Imports.IDs)
		imports := slices.CompactFunc(sourceFile.Imports.IDs, ir.EqualsModuleID)
		if len(imports) != size {
			validation.AddErrorf(
				sourceFile.Imports.Pos,
				"source file %q has an 'imports' section that contains duplicated imports", sourceFile.Header.ModuleID)
		}

		for _, id := range sourceFile.Imports.IDs {
			if err := ir.ValidateModuleID(id); err != nil {
				validation.AddErr(id.Pos, err)
			}
		}
	}

	if sourceFile.Header.Is(ImplSourceFile) {
		if len(sourceFile.Impls.Filenames) > 0 {
			validation.AddErrorf(
				sourceFile.Impls.Pos,
				"implementation file %q has an 'impls' section. The 'impls' section can only be used in base files", sourceFile.Header.ModuleID)
		}
	}

	{
		// Validate impls.
		size := len(sourceFile.Impls.Filenames)

		impls := slices.SortedFunc(slices.Values(sourceFile.Impls.Filenames), ir.CompareFilename)
		impls = slices.CompactFunc(impls, ir.EqualsFilename)

		if len(impls) != size {
			validation.AddErrorf(
				sourceFile.Impls.Pos,
				"file %q has an 'impls' section that contains duplicated implementation files", sourceFile.Header.ModuleID)
		}
	}

	{
		// Validate flags.
		for _, filename := range sourceFile.Flags.Filenames {
			if err := ir.ValidateFilename(filename); err != nil {
				validation.AddErr(filename.Pos, err)
			}
		}
	}

	return validation
}
