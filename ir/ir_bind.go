package ir

import "fmt"

type IrBindCase int

const (
	MarkerBind = IrBindCase(iota)
	TermBind
	TypeBind
)

type IrBind struct {
	Case   IrBindCase
	Marker *string
	Term   *IrSymbol
	Type   *struct {
		Type     IrType
		Solution *IrType
	}
}

func (b IrBind) ID() string {
	switch b.Case {
	case MarkerBind:
		return ""
	case TermBind:
		return b.Term.ID
	case TypeBind:
		return b.Type.Type.TypeID()
	default:
		panic(fmt.Errorf("unhandled IrBindCase %d", b.Case))
	}
}

func NewMarkerBind(id string) IrBind {
	b := IrBind{}
	b.Case = MarkerBind
	b.Marker = &id
	return b
}

func NewTermBind(symbol IrSymbol) IrBind {
	b := IrBind{}
	b.Case = TermBind
	b.Term = &symbol
	return b
}

func NewTypeBind(typ IrType, solution *IrType) IrBind {
	b := IrBind{}
	b.Case = TypeBind
	b.Type = &struct {
		Type     IrType
		Solution *IrType
	}{typ, solution}
	return b
}
