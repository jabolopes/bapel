package stlc

import "fmt"

type Symbol int

const (
	// Symbol is declared, either via an import, or an `impl`, or
	// declared in the current source file.
	DeclSymbol Symbol = iota
	// Symbol is defined in the current source file.
	DefSymbol
)

func (s Symbol) String() string {
	switch s {
	case DeclSymbol:
		return "declaration symbol"
	case DefSymbol:
		return "definition symbol"
	default:
		panic(fmt.Errorf("unhandled Symbol %d", s))
	}
}
