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
