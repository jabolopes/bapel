package query

import (
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
)

func QueryExports(inputFile *os.File) ([]ir.IrDecl, error) {
	decls := []ir.IrDecl{}
	insideExports := false

	parser := bplparser.NewParser(nil /* compiler */)
	parser.Open(inputFile)

loop:
	for parser.Scan() {
		source, err := parser.ParseAny()
		if err != nil {
			return nil, err
		}

		switch {
		case !insideExports && source.Case == bplparser.SectionSource && source.Section == "exports":
			insideExports = true
		case insideExports && source.Case == bplparser.DeclSource:
			decls = append(decls, *source.Decl)
		case insideExports && source.Case == bplparser.EndSource:
			break loop
		}
	}

	if err := parser.ScanErr(); err != nil {
		return nil, err
	}

	return decls, nil
}
