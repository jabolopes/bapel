package query

import (
	"bufio"
	"os"
	"strings"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func QueryExports(inputFile *os.File) ([]ir.IrDecl, error) {
	decls := []ir.IrDecl{}
	insideExports := false

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		args := parser.Words(line)
		if section, _, err := bplparser.ParseSection(args); err == nil {
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

		if decl, _, err := bplparser.ParseDecl(args, false /* named */); err == nil {
			decls = append(decls, decl)
		}
	}

	return decls, nil
}
