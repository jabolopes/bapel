package stlc

import "fmt"

type Symbol int

const (
	ImportSymbol Symbol = iota
	ExportSymbol
	DeclSymbol
	DefSymbol
)

func (s Symbol) String() string {
	switch s {
	case ImportSymbol:
		return "import symbol"
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
