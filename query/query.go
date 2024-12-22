package query

import (
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func QueryExports(inputFile *os.File) ([]ir.IrDecl, error) {
	sources, err := bplparser2.ParseFile(inputFile.Name(), inputFile)
	if err != nil {
		return nil, err
	}

	var decls []ir.IrDecl
	for _, source := range sources {
		switch {
		case source.Is(bplparser.SectionSource) && source.Section.ID == "exports":
			decls = append(decls, source.Section.Decls...)
		case source.Is(bplparser.FunctionSource) && source.Function.Export:
			decls = append(decls, source.Function.Decl())
		case source.Is(bplparser.TypeDefSource) && source.TypeDef.Export:
			decls = append(decls, source.TypeDef.Decl)
		}
	}

	return decls, nil
}
