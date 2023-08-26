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
		args := parser.Words()
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

		if _, err := parser.ParseEnd(args); err == nil {
			break
		}

		if decl, _, err := parser.ParseDecl(args, false /* named */); err == nil {
			decls = append(decls, decl)
		}
	}

	if err := parser.ScanErr(); err != nil {
		return nil, err
	}

	return decls, nil
}
