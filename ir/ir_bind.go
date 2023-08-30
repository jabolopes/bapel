package ir

type IrBindCase int

const (
	SymbolBind = IrBindCase(iota)
	ScopeBind
)

type IrBind struct {
	Case   IrBindCase
	Symbol *IrSymbol
	Scope  *string
}

func NewSymbolBind(symbol IrSymbol) IrBind {
	b := IrBind{}
	b.Case = SymbolBind
	b.Symbol = &symbol
	return b
}

func NewScopeBind(id string) IrBind {
	b := IrBind{}
	b.Case = ScopeBind
	b.Scope = &id
	return b
}
