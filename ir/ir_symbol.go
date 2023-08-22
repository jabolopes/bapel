package ir

type IrSymbolCase int

const (
	ImportSymbol = IrSymbolCase(iota)
	ExportSymbol
	DeclSymbol
	FunctionSymbol
	ReferenceSymbol
	StructSymbol
)

type IrSymbol struct {
	Case IrSymbolCase
	Decl IrDecl
}

func NewSymbol(c IrSymbolCase, decl IrDecl) IrSymbol {
	return IrSymbol{c, decl}
}
