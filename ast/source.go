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

// TODO: Consolidate on Format.
func (s *declSource) String() string {
	return s.Decl.String()
}

type Source struct {
	Case     SourceCase
	Decl     *declSource
	Function *ir.IrFunction
}

func (s Source) Format(f fmt.State, verb rune) {
	if s.Case == 0 && s.Decl == nil {
		return
	}

	switch s.Case {
	case DeclSource:
		if addMetadata := f.Flag('+'); addMetadata {
			s.Decl.Decl.Pos.Format(f, verb)
		}
		fmt.Fprint(f, s.Decl.String())
	case FunctionSource:
		fmt.Fprintf(f, fmt.FormatString(f, 's'), s.Function)
	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func (s Source) Is(c SourceCase) bool {
	return s.Case == c
}

func NewDeclSource(decl ir.IrDecl) Source {
	return Source{
		Case: DeclSource,
		Decl: &declSource{decl},
	}
}

func NewFunctionSource(function ir.IrFunction) Source {
	return Source{
		Case:     FunctionSource,
		Function: &function,
	}
}
