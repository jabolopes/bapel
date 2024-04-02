package ir

import "fmt"

type IrBindCase int

const (
	MarkerBind IrBindCase = iota
	DeclBind
)

type IrSymbol int

const (
	ImportSymbol IrSymbol = iota
	ExportSymbol
	DeclSymbol
	DefSymbol
)

func (s IrSymbol) String() string {
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
		panic(fmt.Errorf("unhandled IrSymbol %d", s))
	}
}

type IrBind struct {
	Case   IrBindCase
	Symbol IrSymbol
	Marker *string
	Decl   *IrDecl
}

func (b IrBind) String() string {
	switch b.Case {
	case MarkerBind:
		return fmt.Sprintf("<|%s", *b.Marker)
	case DeclBind:
		return b.Decl.String()
	default:
		panic(fmt.Errorf("unhandled IrBindCase %d", b.Case))
	}
}

func (b IrBind) ID() (string, bool) {
	switch b.Case {
	case MarkerBind:
		return "", false
	case DeclBind:
		switch b.Decl.Case {
		case TypeDecl:
			return b.Decl.Type().ID()
		case TermDecl:
			return b.Decl.Term.ID, true
		default:
			panic(fmt.Errorf("unhandled IrDeclCase %d", b.Decl.Case))
		}
	default:
		panic(fmt.Errorf("unhandled IrBindCase %d", b.Case))
	}
}

func NewMarkerBind(id string) IrBind {
	b := IrBind{}
	b.Symbol = DefSymbol
	b.Case = MarkerBind
	b.Marker = &id
	return b
}

func NewDeclBind(symbol IrSymbol, decl IrDecl) IrBind {
	return IrBind{
		Case:   DeclBind,
		Symbol: symbol,
		Decl:   &decl,
	}
}
