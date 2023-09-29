package ir

import "fmt"

type IrBindCase int

const (
	MarkerBind = IrBindCase(iota)
	TermBind
	TypeBind
)

type IrBind struct {
	Case       IrBindCase
	SymbolCase IrSymbolCase
	Marker     *string
	Term       *struct {
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

func NewMarkerBind(id string) IrBind {
	b := IrBind{}
	b.SymbolCase = DefSymbol
	b.Case = MarkerBind
	b.Marker = &id
	return b
}

func NewTermBind(symbolCase IrSymbolCase, decl IrDecl) IrBind {
	b := IrBind{}
	b.SymbolCase = symbolCase
	b.Case = TermBind
	b.Term = &struct {
		Decl IrDecl
	}{decl}
	return b
}

func NewTypeBind(symbol IrSymbolCase, typ IrType, solution *IrType) IrBind {
	b := IrBind{}
	b.SymbolCase = symbol
	b.Case = TypeBind
	b.Type = &struct {
		Type     IrType
		Solution *IrType
	}{typ, solution}
	return b
}
