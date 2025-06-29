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

func parseModuleNoBody(filename string) (ast.Module, error) {
	input, err := os.Open(filename)
	if err != nil {
		return ast.Module{}, fmt.Errorf("failed to parse module metadata: %v", err)
	}
	defer input.Close()

	module, err := bplparser2.ParseFile(filename, input)
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
	input, err := os.Open(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to query module file declarations: %v", err)
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

// Queries all the declarations of a module, recursing into the
// implementation files of the `impls` section.
//
// moduleID: identifier of the module, e.g., 'core'.
func QueryModuleDecls(moduleID ast.ModuleID) ([]ir.IrDecl, error) {
	baseFilename := ast.ModuleBaseFilename(moduleID)

	input, err := os.Open(baseFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to query module declarations: %v", err)
	}
	defer input.Close()

	module, decls, err := queryDeclsBplFile(baseFilename, input)
	if err != nil {
		return nil, err
	}

	for _, relativeImplFilename := range module.Impls.IDs {
		implFilename := ast.ModuleImplFilename(baseFilename, relativeImplFilename)

		implDecls, err := QueryFileDecls(implFilename)
		if err != nil {
			return nil, err
		}

		decls = append(decls, implDecls...)
	}

	return decls, nil
}

// Queries all the exports of a module, recursing into the implementation files
// of the `impls` section.
//
// moduleID: identifier of the module, e.g., 'core'.
func QueryModuleExports(moduleID ast.ModuleID) ([]ir.IrDecl, error) {
	decls, err := QueryModuleDecls(moduleID)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(decls, func(decl ir.IrDecl) bool { return !decl.Export }), nil
}

// Queries module metadata (e.g. imports, impls, flags, etc),recursing into the
// implementation files defined in the `impls` section to discover all the
// imports, all the flags, etc.
//
// The module body is not populated in the ast.Module because this only returns
// module metadata.
func QueryModuleMetadata(moduleID ast.ModuleID) (ast.Module, error) {
	baseFilename := ast.ModuleBaseFilename(moduleID)

	module, err := parseModuleNoBody(baseFilename)
	if err != nil {
		return ast.Module{}, err
	}

	for _, relativeImplFilename := range module.Impls.IDs {
		if !strings.HasSuffix(relativeImplFilename.Value, ".bpl") {
			continue
		}

		implFilename := ast.ModuleImplFilename(baseFilename, relativeImplFilename)

		implModule, err := parseModuleNoBody(implFilename)
		if err != nil {
			return ast.Module{}, err
		}

		module.Imports.IDs = append(module.Imports.IDs, implModule.Imports.IDs...)
		module.Flags.IDs = append(module.Flags.IDs, implModule.Flags.IDs...)
		module.Errors = append(module.Errors, implModule.Errors...)
	}

	slices.SortFunc(module.Imports.IDs, ast.CompareModuleID)
	module.Imports.IDs = slices.CompactFunc(module.Imports.IDs, func(id1, id2 ast.ModuleID) bool {
		return ast.CompareModuleID(id1, id2) == 0
	})

	slices.SortFunc(module.Flags.IDs, ast.CompareID)
	module.Flags.IDs = slices.Compact(module.Flags.IDs)

	return module, nil
}
