package ir

import "fmt"

type IrBindCase int

const (
	MarkerBind = IrBindCase(iota)
	TermBind
	TypeBind
)

type IrSymbol int

const (
	ImportSymbol = IrSymbol(iota)
	ExportSymbol
	DeclSymbol
	DefSymbol
)

type IrBind struct {
	Case   IrBindCase
	Symbol IrSymbol
	Marker *string
	Term   *struct {
		Decl IrDecl
	}
	Type *struct {
		Type     IrType
		Solution *IrType
	}
}

func (b IrBind) String() string {
	switch b.Case {
	case MarkerBind:
		return fmt.Sprintf("<|%s", *b.Marker)
	case TermBind:
		return fmt.Sprintf("%s:%s", b.Term.Decl.ID, b.Term.Decl.Type)
	case TypeBind:
		if b.Type.Solution == nil {
			return fmt.Sprintf("type %s", b.Type.Type)
		}
		return fmt.Sprintf("type %s = %s", b.Type.Type, *b.Type.Solution)
	default:
		panic(fmt.Errorf("unhandled IrBindCase %d", b.Case))
	}
}

func (b IrBind) ID() (string, bool) {
	switch b.Case {
	case MarkerBind:
		return "", false
	case TermBind:
		return b.Term.Decl.ID, true
	case TypeBind:
		return b.Type.Type.ID()
	default:
		panic(fmt.Errorf("unhandled IrBindCase %d", b.Case))
	}
}

func (b IrBind) Decl() (IrDecl, bool) {
	id, ok := b.ID()
	if !ok {
		return IrDecl{}, false
	}

	switch b.Case {
	case TermBind:
		return NewTermDecl(id, b.Term.Decl.Type), true

	case TypeBind:
		if b.Type.Solution == nil {
			return NewTypeDecl(id, b.Type.Type), true
		}
		return NewTypeDecl(id, *b.Type.Solution), true

	default:
		return IrDecl{}, false
	}
}

func NewMarkerBind(id string) IrBind {
	b := IrBind{}
	b.Symbol = DefSymbol
	b.Case = MarkerBind
	b.Marker = &id
	return b
}

func NewTermBind(symbolCase IrSymbol, decl IrDecl) IrBind {
	b := IrBind{}
	b.Symbol = symbolCase
	b.Case = TermBind
	b.Term = &struct {
		Decl IrDecl
	}{decl}
	return b
}

func NewTypeBind(symbol IrSymbol, typ IrType, solution *IrType) IrBind {
	b := IrBind{}
	b.Symbol = symbol
	b.Case = TypeBind
	b.Type = &struct {
		Type     IrType
		Solution *IrType
	}{typ, solution}
	return b
}

func NewBindFromDecl(symbol IrSymbol, decl IrDecl) IrBind {
	switch decl.Case {
	case TypeDecl:
		return NewTypeBind(symbol, NewNameType(decl.ID), &decl.Type)
	case TermDecl:
		return NewTermBind(symbol, decl)
	default:
		panic(fmt.Sprintf("unhandled IrDeclCase %d", decl.Case))
	}
}
