package query

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
)

const (
	bplDeclAnnotation = "// @bpl: "
)

type filter = func(string) (string, bool)

func queryAnnotationNonBplFile(inputFilename string, input io.Reader, filter filter) ([]ir.IrDecl, error) {
	var parser *parse.Parser

	var decls []ir.IrDecl
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line, ok := filter(scanner.Text())
		if !ok {
			continue
		}

		if parser == nil {
			var err error
			if parser, err = parse.NewWithSymbol("Decl"); err != nil {
				return nil, err
			}
		}

		parser.Open(inputFilename, strings.NewReader(line))

		decl, err := parse.Parse[ir.IrDecl](parser)
		if err != nil {
			return nil, err
		}

		decls = append(decls, decl)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return decls, nil
}

func queryDeclsBplFile(inputFilename string) (ast.Module, []ir.IrDecl, error) {
	module, err := parse.ParseModuleFile(inputFilename)
	if err != nil {
		return ast.Module{}, nil, err
	}

	var decls []ir.IrDecl
	for _, source := range module.Body {
		switch {
		case source.Is(ast.FunctionSource):
			decls = append(decls, source.Function.Decl())
		case source.Is(ast.DeclSource):
			decls = append(decls, source.Decl.Decl)
		}
	}

	return module, decls, nil
}

func parseModuleNoBody(inputFilename string) (ast.Module, error) {
	module, err := parse.ParseModuleFile(inputFilename)
	if err != nil {
		return ast.Module{}, err
	}

	// TODO: At this stage, the builder only cares about the build graph, so we
	// could optimize the build process by not parsing the module body.
	module.Body = nil

	return module, nil
}

// Queries all the declarations of a file, without recursing into the
// implementation files of the `impls` section.
//
// The file can be a base file or an implementation file.
//
// To recurse into the `impls` section, `QueryModuleDecls` instead.
func QueryFileDecls(inputFilename string) ([]ir.IrDecl, error) {
	if path.Ext(inputFilename) == ".bpl" {
		_, decls, err := queryDeclsBplFile(inputFilename)
		return decls, err
	}

	input, err := os.Open(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to query module file declarations: %v", err)
	}
	defer input.Close()

	return queryAnnotationNonBplFile(inputFilename, input, func(line string) (decl string, ok bool) {
		decl = strings.TrimPrefix(line, bplDeclAnnotation)
		ok = len(decl) != len(line)
		return
	})
}
