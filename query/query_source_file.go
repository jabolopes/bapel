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
	var imports []ir.ModuleID
	var decls []ir.IrDecl
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		isModuleImport := strings.HasPrefix(line, "import ") &&
			!strings.HasPrefix(line, "import :") && // Exclude module partition imports.
			strings.HasSuffix(line, ";")

		if isModuleImport {
			line = strings.TrimPrefix(line, "import ")
			line = strings.TrimSuffix(line, ";")

			moduleID := ir.NewModuleID(line, ir.Pos{})
			if err := ir.ValidateModuleID(moduleID); err != nil {
				return SourceFileQuery{}, err
			}

			imports = append(imports, moduleID)
		}

		// TODO: Inline filter.
		line, ok := filter(scanner.Text())
		if !ok {
			continue
		}

		decl, err := parse.ParseSymbol[ir.IrDecl]("Decl", inputFilename, strings.NewReader(line))
		if err != nil {
			return SourceFileQuery{}, err
		}

		decls = append(decls, decl)
	}

	if err := scanner.Err(); err != nil {
		return SourceFileQuery{}, err
	}

	return SourceFileQuery{
		imports,
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
