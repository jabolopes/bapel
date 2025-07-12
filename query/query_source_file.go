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

func queryAnnotationNonBplFile(inputFilename string, input io.Reader, filter filter) (SourceFileQuery, error) {
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
				return SourceFileQuery{}, err
			}
		}

		parser.Open(inputFilename, strings.NewReader(line))

		decl, err := parse.Parse[ir.IrDecl](parser)
		if err != nil {
			return SourceFileQuery{}, err
		}

		decls = append(decls, decl)
	}

	if err := scanner.Err(); err != nil {
		return SourceFileQuery{}, err
	}

	return SourceFileQuery{
		nil, /* Imports */
		nil, /* Impls */
		nil, /* flags */
		decls,
	}, nil
}

func queryDeclsBplFile(inputFilename string) (SourceFileQuery, error) {
	sourceFile, err := parse.ParseSourceFile(inputFilename)
	if err != nil {
		return SourceFileQuery{}, err
	}

	var decls []ir.IrDecl
	for _, source := range sourceFile.Body {
		switch {
		case source.Is(ast.DeclSource):
			decls = append(decls, source.Decl.Decl)
		case source.Is(ast.FunctionSource):
			decls = append(decls, source.Function.Decl())
		}
	}

	return SourceFileQuery{
		sourceFile.Imports.IDs,
		sourceFile.Impls.Filenames,
		sourceFile.Flags.Filenames,
		decls,
	}, nil
}

// Queries the source file without recursing into the implementation
// files of the `impls` section. Returns the source file metadata, and
// the declarations owned by that file.
//
// The file can be a base file or an implementation file.
//
// To recurse into the `impls` section, `QueryModuleDecls` instead.
func QuerySourceFile(inputFilename string) (SourceFileQuery, error) {
	if path.Ext(inputFilename) == ".bpl" {
		return queryDeclsBplFile(inputFilename)
	}

	input, err := os.Open(inputFilename)
	if err != nil {
		return SourceFileQuery{}, fmt.Errorf("failed to query source file %q: %v", inputFilename, err)
	}
	defer input.Close()

	return queryAnnotationNonBplFile(inputFilename, input, func(line string) (decl string, ok bool) {
		decl = strings.TrimPrefix(line, bplDeclAnnotation)
		ok = len(decl) != len(line)
		return
	})
}
