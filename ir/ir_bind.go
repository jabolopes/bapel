package ir

type IrBindCase int

const (
	SymbolBind = IrBindCase(iota)
	MarkerBind
)

type IrBind struct {
	Case   IrBindCase
	Symbol *IrSymbol
	Marker *string
}

func NewSymbolBind(symbol IrSymbol) IrBind {
	b := IrBind{}
	b.Case = SymbolBind
	b.Symbol = &symbol
	return b
}

func NewMarkerBind(id string) IrBind {
	b := IrBind{}
	b.Case = MarkerBind
	b.Marker = &id
	return b
}
