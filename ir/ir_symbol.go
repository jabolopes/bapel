package ir

type IrSymbolCase int

const (
	ImportSymbol = IrSymbolCase(iota)
	ExportSymbol
	DeclSymbol
	DefSymbol
)

type IrSymbol struct {
	Case     IrSymbolCase
	DeclCase IrDeclCase
	ID       string
	Type     *IrType
}

func NewSymbolFromDecl(c IrSymbolCase, decl IrDecl) IrSymbol {
	return IrSymbol{c, decl.Case, decl.ID, &decl.Type}
}
