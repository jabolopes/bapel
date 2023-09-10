package ir

type IrBindCase int

const (
	MarkerBind = IrBindCase(iota)
	TermBind
)

type IrBind struct {
	Case   IrBindCase
	Marker *string
	Term   *IrSymbol
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
