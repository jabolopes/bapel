package stlc

import "fmt"

type Symbol int

const (
	// Symbol is imported into the current module file.
	//
	// It could be imported via an `imports` or `impls` section.
	ImportSymbol Symbol = iota
	// Symbol declared in the current impl file or module file.
	//
	// Can shadow imported symbols. This is necessary to avoid the
	// problem of some imported module defining a new toplevel symbol
	// that then would break the current module because it shadows a
	// symbol that the current module already declared / defined.
	//
	// If a symbol is declared, it must precede its definition.
	DeclSymbol
	// Symbol defined in the current impl file or module file.
	//
	// Can shadow imported symbols. This is necessary to avoid the
	// problem of some imported module defining a new toplevel symbol
	// that then would break the current module because it shadows a
	// symbol that the current module already declared / defined.
	//
	// All declared symbols must be defined.
	DefSymbol
)

func (s Symbol) String() string {
	switch s {
	case ImportSymbol:
		return "import symbol"
	case DeclSymbol:
		return "declaration symbol"
	case DefSymbol:
		return "definition symbol"
	default:
		panic(fmt.Errorf("unhandled Symbol %d", s))
	}
}
