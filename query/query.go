package query

import (
	"io"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func QueryExports(inputFilename string, input io.Reader) ([]ir.IrDecl, error) {
	module, err := bplparser2.ParseFile(inputFilename, input)
	if err != nil {
		return nil, err
	}

	decls := module.Exports.Decls
	for _, source := range module.Body {
		switch {
		case source.Is(ast.FunctionSource) && source.Function.Export:
			decls = append(decls, source.Function.Decl())
		case source.Is(ast.TypeDefSource) && source.TypeDef.Export:
			decls = append(decls, source.TypeDef.Decl)
		}
	}

	return decls, nil
}

func QueryDecls(inputFilename string, input io.Reader) ([]ir.IrDecl, error) {
	module, err := bplparser2.ParseFile(inputFilename, input)
	if err != nil {
		return nil, err
	}

	decls := module.Exports.Decls
	for _, source := range module.Body {
		switch {
		case source.Is(ast.FunctionSource):
			decls = append(decls, source.Function.Decl())
		case source.Is(ast.TypeDefSource):
			decls = append(decls, source.TypeDef.Decl)
		}
	}

	return decls, nil
}
