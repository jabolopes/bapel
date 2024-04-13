package query

import (
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
)

func QueryExports(inputFile *os.File) ([]ir.IrDecl, error) {
	decls := []ir.IrDecl{}

	parser := bplparser.NewParser()
	parser.Open(inputFile)

	for parser.Scan() {
		source, err := parser.ParseAny()
		if err != nil {
			return nil, err
		}

		if source.Case == bplparser.SectionSource && source.Section.ID == "exports" {
			decls = append(decls, source.Section.Decls...)
		}
	}

	if err := parser.ScanErr(); err != nil {
		return nil, err
	}

	return decls, nil
}
