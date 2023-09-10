package ir

type IrSymbolCase int

const (
	ImportSymbol = IrSymbolCase(iota)
	ExportSymbol
	DeclSymbol
	DefSymbol
)

type IrSymbol struct {
	Case IrSymbolCase
	Decl IrDecl
}

func NewSymbolFromDecl(c IrSymbolCase, decl IrDecl) IrSymbol {
	return IrSymbol{c, decl}
}
