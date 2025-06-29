package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	ComponentSource SourceCase = iota
	DeclSource
	ExportSource
	FunctionSource
	ImportSource
	ImplSource
)

type declSource struct {
	Decl ir.IrDecl
}

func (s *declSource) String() string {
	return s.Decl.String()
}

type importSource struct {
	ModuleID ModuleID
	Decl     ir.IrDecl
}

func (s *importSource) String() string {
	return fmt.Sprintf("import %s %s", s.ModuleID, s.Decl)
}

type implSource struct {
	ModuleFilename string // e.g., 'core_impl.bpl' or 'core_impl.cc'
	Decl           ir.IrDecl
}

func (s *implSource) String() string {
	return fmt.Sprintf("impl %s %s", s.ModuleFilename, s.Decl)
}

type Source struct {
	Case      SourceCase
	Component *ir.IrComponent
	Decl      *declSource
	Function  *ir.IrFunction
	Import    *importSource
	Impl      *implSource
	// Position in source file.
	Pos ir.Pos
}

func (s Source) String() string {
	if s.Case == 0 && s.Component == nil {
		return ""
	}

	switch s.Case {
	case ComponentSource:
		return s.Component.String()
	case DeclSource:
		return s.Decl.String()
	case FunctionSource:
		return s.Function.String()
	case ImportSource:
		return s.Import.String()
	case ImplSource:
		return s.Impl.String()
	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func (s Source) Format(f fmt.State, verb rune) {
	if s.Case == 0 && s.Component == nil {
		return
	}

	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprint(f, s.String())
}

func (s Source) Is(c SourceCase) bool {
	return s.Case == c
}

func NewComponentSource(component ir.IrComponent) Source {
	return Source{
		Case:      ComponentSource,
		Component: &component,
	}
}

func NewDeclSource(decl ir.IrDecl) Source {
	return Source{
		Case: DeclSource,
		Decl: &declSource{decl},
		Pos:  decl.Pos,
	}
}

func NewFunctionSource(function ir.IrFunction) Source {
	return Source{
		Case:     FunctionSource,
		Function: &function,
		Pos:      function.Pos,
	}
}

func NewImportSource(moduleID ModuleID, decl ir.IrDecl) Source {
	return Source{
		Case:   ImportSource,
		Import: &importSource{moduleID, decl},
		Pos:    decl.Pos,
	}
}

func NewImplSource(moduleFilename string, decl ir.IrDecl) Source {
	return Source{
		Case: ImplSource,
		Impl: &implSource{moduleFilename, decl},
		Pos:  decl.Pos,
	}
}
