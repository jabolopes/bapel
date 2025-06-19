package ast

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	ComponentSource SourceCase = iota
	FunctionSource
	DefSymbolSource
)

type defSymbolSource struct {
	Export bool
	// Whether this is a declaration instead of a definition.
	//
	// If this is true, the this is a symbol declaration rather than a symbol definition.
	//
	// For example, a declaration:
	//   id: forall ['a] a -> a
	//
	// And a definition:
	//   type Point = struct {...}
	IsDecl bool
	Decl   ir.IrDecl
}

func (s *defSymbolSource) String() string {
	var b strings.Builder
	if s.Export {
		b.WriteString("export ")
	}
	b.WriteString(s.Decl.String())
	return b.String()
}

type Source struct {
	Case      SourceCase
	Component *ir.IrComponent
	Function  *ir.IrFunction
	DefSymbol *defSymbolSource
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
	case FunctionSource:
		return s.Function.String()
	case DefSymbolSource:
		return s.DefSymbol.String()
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

func NewFunctionSource(function ir.IrFunction) Source {
	return Source{
		Case:     FunctionSource,
		Function: &function,
	}
}

func NewDefSymbolSource(export, isDecl bool, decl ir.IrDecl) Source {
	return Source{
		Case:      DefSymbolSource,
		DefSymbol: &defSymbolSource{export, isDecl, decl},
	}
}
