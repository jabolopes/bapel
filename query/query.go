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
	// TODO: The export should be handled by the grammar.
	bplExportAnnotation = "// @bpl: export "
)

type filter = func(string) (string, bool, bool)

func queryAnnotationNonBplFile(inputFilename string, input io.Reader, filter filter) ([]ir.IrDeclE, error) {
	var parser *bplparser2.Parser

	var decls []ir.IrDeclE
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line, isExport, ok := filter(scanner.Text())
		if !ok {
			continue
		}

		// TODO: Find a way to get rid of this implementation detail.
		//
		// Required by the grammar.
		line += ";"

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

		decls = append(decls, ir.NewDeclE(decl, isExport))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return decls, nil
}

func queryDeclsBplFile(inputFilename string, input io.Reader) (ast.Module, []ir.IrDeclE, error) {
	module, err := bplparser2.ParseFile(inputFilename, input)
	if err != nil {
		return ast.Module{}, nil, err
	}

	var decls []ir.IrDeclE
	for _, decl := range module.Exports.Decls {
		decls = append(decls, ir.NewDeclE(decl, true /* export */))
	}

	for _, source := range module.Body {
		switch {
		case source.Is(ast.FunctionSource):
			decls = append(decls, ir.NewDeclE(source.Function.Decl(), source.Function.Export))
		case source.Is(ast.DeclSource):
			decls = append(decls, ir.NewDeclE(source.Decl.Decl, false /* export */))
		case source.Is(ast.ExportSource):
			decls = append(decls, ir.NewDeclE(source.Export.Decl, true /* export */))
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
func QueryFileDecls(inputFilename string, input io.Reader) ([]ir.IrDeclE, error) {
	if path.Ext(inputFilename) == ".bpl" {
		_, decls, err := queryDeclsBplFile(inputFilename, input)
		return decls, err
	}
	return queryAnnotationNonBplFile(inputFilename, input, func(line string) (decl string, isExport, ok bool) {
		if strings.HasPrefix(line, bplExportAnnotation) {
			decl = strings.TrimPrefix(line, bplExportAnnotation)
			isExport = true
			ok = len(decl) != len(line)
		} else {
			decl = strings.TrimPrefix(line, bplDeclAnnotation)
			isExport = false
			ok = len(decl) != len(line)
		}

		return
	})
}

// Queries all the declarations of a module, including the `impls` section.
//
// moduleID: identifier of the module, e.g., 'core'.
func QueryModuleDecls(moduleID string) ([]ir.IrDeclE, error) {
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
		input, err := os.Open(filename.Value)
		if err != nil {
			return nil, err
		}
		// TODO: Avoid keeping files open during the loop.
		defer input.Close()

		implDecls, err := QueryFileDecls(filename.Value, input)
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
func QueryModuleExports(moduleID string) ([]ir.IrDeclE, error) {
	decls, err := QueryModuleDecls(moduleID)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(decls, func(decl ir.IrDeclE) bool { return !decl.Export }), nil
}
