package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	DeclSource SourceCase = iota
	FunctionSource
)

type declSource struct {
	Decl ir.IrDecl
}

func (s *declSource) String() string {
	return s.Decl.String()
}

type Source struct {
	Case     SourceCase
	Decl     *declSource
	Function *ir.IrFunction
	// Position in source file.
	Pos ir.Pos
}

// TODO: Consolidate on Format() instead of String().
func (s Source) String() string {
	if s.Case == 0 && s.Decl == nil {
		return ""
	}

	switch s.Case {
	case DeclSource:
		return s.Decl.String()
	case FunctionSource:
		return fmt.Sprintf("%s", s.Function)
	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func (s Source) Format(f fmt.State, verb rune) {
	if s.Case == 0 && s.Decl == nil {
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
