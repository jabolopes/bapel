package query

import (
	"io"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func QueryExports(inputFilename string, input io.Reader) ([]ir.IrDecl, error) {
	sources, err := bplparser2.ParseFile(inputFilename, input)
	if err != nil {
		return nil, err
	}

	var decls []ir.IrDecl
	for _, source := range sources {
		switch {
		case source.Is(ast.ExportsSource):
			decls = append(decls, source.Exports.Decls...)
		case source.Is(ast.FunctionSource) && source.Function.Export:
			decls = append(decls, source.Function.Decl())
		case source.Is(ast.TypeDefSource) && source.TypeDef.Export:
			decls = append(decls, source.TypeDef.Decl)
		}
	}

	return decls, nil
}
