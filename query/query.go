package query

import (
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func QueryExports(inputFile *os.File) ([]ir.IrDecl, error) {
	decls := []ir.IrDecl{}
	insideExports := false

	p := bplparser.NewParser(nil /* compiler */)
	p.Open(inputFile)
	for p.Scan() {
		args := p.Words()
		if section, _, err := p.ParseSection(args); err == nil {
			if section != "exports" {
				continue
			}

			insideExports = true
			continue
		}

		if !insideExports {
			continue
		}

		if _, err := parser.ShiftToken(args, "}"); err == nil {
			break
		}

		if decl, _, err := p.ParseDecl(args, false /* named */); err == nil {
			decls = append(decls, decl)
		}
	}

	if err := p.ScanErr(); err != nil {
		return nil, err
	}

	return decls, nil
}
