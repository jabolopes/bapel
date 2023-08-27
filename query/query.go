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
	for parser.Scan() {
		if section, err := parser.ParseSection(); err == nil {
			if section != "exports" {
				continue
			}

			insideExports = true
			continue
		}

		if !insideExports {
			continue
		}

		if err := parser.ParseEnd(); err == nil {
			break
		}

		if decl, err := parser.ParseDecl(false /* named */); err == nil {
			decls = append(decls, decl)
		}
	}

	if err := parser.ScanErr(); err != nil {
		return nil, err
	}

	return decls, nil
}
