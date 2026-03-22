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

type functionSource struct {
	Function
}

type Source struct {
	Case     SourceCase
	Decl     *declSource
	Function *functionSource
}

func (s Source) Format(f fmt.State, verb rune) {
	if s.Case == 0 && s.Decl == nil {
		return
	}

	switch s.Case {
	case DeclSource:
		fmt.Fprintf(f, fmt.FormatString(f, 's'), s.Decl.Decl)
	case FunctionSource:
		fmt.Fprintf(f, fmt.FormatString(f, 's'), s.Function.Function)
	default:
		panic(fmt.Errorf("unhandled %T %d", s.Case, s.Case))
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

func NewFunctionSource(function Function) Source {
	return Source{
		Case:     FunctionSource,
		Function: &functionSource{function},
	}
}
