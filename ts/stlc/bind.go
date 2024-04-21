package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type BindCase int

const (
	MarkerBind BindCase = iota
	DeclBind
)

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

type Bind struct {
	Case   BindCase
	Symbol Symbol
	Marker *string
	Decl   *ir.IrDecl
}

func (b Bind) String() string {
	switch b.Case {
	case MarkerBind:
		return fmt.Sprintf("<|%s", *b.Marker)
	case DeclBind:
		return b.Decl.String()
	default:
		panic(fmt.Errorf("unhandled BindCase %d", b.Case))
	}
}

func (b Bind) ID() (string, bool) {
	switch b.Case {
	case MarkerBind:
		return "", false
	case DeclBind:
		switch b.Decl.Case {
		case ir.TypeDecl:
			return b.Decl.Type().ID()
		case ir.TermDecl:
			return b.Decl.Term.ID, true
		default:
			panic(fmt.Errorf("unhandled %T %d", b.Decl.Case, b.Decl.Case))
		}
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func NewMarkerBind(id string) Bind {
	b := Bind{}
	b.Symbol = DefSymbol
	b.Case = MarkerBind
	b.Marker = &id
	return b
}

func NewDeclBind(symbol Symbol, decl ir.IrDecl) Bind {
	return Bind{
		Case:   DeclBind,
		Symbol: symbol,
		Decl:   &decl,
	}
}
