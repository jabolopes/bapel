package stlc

import "fmt"

type Symbol int

const (
	// Symbol exported by another module and imported into the current
	// module.
	ImportSymbol Symbol = iota
	// Symbol defined in an impl file of the current module. It doesn't
	// say whether that impl file exports the symbol or it doesn't.
	ImplSymbol
	// Symbol defined and exported in the current impl file or module file.
	ExportSymbol
	// Symbol declared in the current impl file or module file.
	DeclSymbol
	// Symbol defined in the current impl file or module file.
	DefSymbol
)

func (s Symbol) String() string {
	switch s {
	case ImportSymbol:
		return "import symbol"
	case ImplSymbol:
		return "impl symbol"
	case ExportSymbol:
		return "export symbol"
	case DeclSymbol:
		return "declaration symbol"
	case DefSymbol:
		return "definition symbol"
	default:
		panic(fmt.Errorf("unhandled Symbol %d", s))
	}
}
