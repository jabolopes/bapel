package query

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

const (
	bplDeclAnnotation = "// @bpl: "
)

type filter = func(string) (string, bool)

func queryAnnotationNonBplFile(inputFilename string, input io.Reader, filter filter) ([]ir.IrDecl, error) {
	var parser *bplparser2.Parser

	var decls []ir.IrDecl
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line, ok := filter(scanner.Text())
		if !ok {
			continue
		}

		if parser == nil {
			var err error
			if parser, err = bplparser2.NewWithSymbol("Decl"); err != nil {
				return nil, err
			}
		}

		parser.Open(inputFilename, strings.NewReader(line))

		decl, err := bplparser2.Parse[ir.IrDecl](parser)
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

func queryDeclsBplFile(inputFilename string, input io.Reader) (ast.Module, []ir.IrDecl, error) {
	module, err := bplparser2.ParseFile(inputFilename, input)
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

// Queries all the declarations of a file. It can be a base module file or an
// implementation module file.
//
// This does not query all module declarations since it only looks at one file
// and it does not automatically traverse the `impls` section. Use
// `QueryModuleDecls` for that.
func QueryFileDecls(inputFilename string) ([]ir.IrDecl, error) {
	input, err := os.Open(inputFilename)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	if path.Ext(inputFilename) == ".bpl" {
		_, decls, err := queryDeclsBplFile(inputFilename, input)
		return decls, err
	}
	return queryAnnotationNonBplFile(inputFilename, input, func(line string) (decl string, ok bool) {
		decl = strings.TrimPrefix(line, bplDeclAnnotation)
		ok = len(decl) != len(line)

		return
	})
}

// Queries all the declarations of a module, including the `impls` section.
//
// moduleID: identifier of the module, e.g., 'core'.
func QueryModuleDecls(moduleID string) ([]ir.IrDecl, error) {
	inputFilename := fmt.Sprintf("%s.bpl", moduleID)

	input, err := os.Open(inputFilename)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	module, decls, err := queryDeclsBplFile(inputFilename, input)
	if err != nil {
		return nil, err
	}

	for _, filename := range module.Impls.IDs {
		implDecls, err := QueryFileDecls(filename.Value)
		if err != nil {
			return nil, err
		}

		decls = append(decls, implDecls...)
	}

	return decls, nil
}

// Queries all the exports of a module, including the `impls` section.
//
// moduleID: identifier of the module, e.g., 'core'.
func QueryModuleExports(moduleID string) ([]ir.IrDecl, error) {
	decls, err := QueryModuleDecls(moduleID)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(decls, func(decl ir.IrDecl) bool { return !decl.Export }), nil
}
