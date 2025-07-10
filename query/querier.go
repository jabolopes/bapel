package query

import (
	"slices"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
)

type Querier struct {
	finder moduleFinder
}

func (q Querier) SourceFileBaseSourceFilename(moduleID ast.ModuleID) ast.Filename {
	return q.finder.sourceFileBaseSourceFilename(moduleID)
}

func (q Querier) SourceFileImplFilename(baseFilename ast.Filename, relativeImplFilename ast.Filename) ast.Filename {
	return q.finder.sourceFileImplFilename(baseFilename, relativeImplFilename)
}

func (q Querier) QueryModuleDecls(moduleID ast.ModuleID) ([]ir.IrDecl, error) {
	baseFilename := q.finder.sourceFileBaseSourceFilename(moduleID)

	sourceFile, decls, err := queryDeclsBplFile(baseFilename.Value)
	if err != nil {
		return nil, err
	}

	for _, relativeImplFilename := range sourceFile.Impls.Filenames {
		implFilename := q.finder.sourceFileImplFilename(baseFilename, relativeImplFilename)

		implDecls, err := QuerySourceFileDecls(implFilename.Value)
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
func (q Querier) QueryModuleExports(moduleID ast.ModuleID) ([]ir.IrDecl, error) {
	decls, err := q.QueryModuleDecls(moduleID)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(decls, func(decl ir.IrDecl) bool { return !decl.Export }), nil
}

// Queries module metadata (e.g. imports, impls, flags, etc),recursing into the
// implementation files defined in the `impls` section to discover all the
// imports, all the flags, etc.
//
// The module body is not populated in the ast.SourceFile because this only returns
// module metadata.
func (q Querier) QueryModuleMetadata(moduleID ast.ModuleID) (ast.SourceFile, error) {
	baseFilename := q.finder.sourceFileBaseSourceFilename(moduleID)

	sourceFile, err := parseSourceFileNoBody(baseFilename.Value)
	if err != nil {
		return ast.SourceFile{}, err
	}

	for _, relativeImplFilename := range sourceFile.Impls.Filenames {
		if !strings.HasSuffix(relativeImplFilename.Value, ".bpl") {
			continue
		}

		implFilename := q.finder.sourceFileImplFilename(baseFilename, relativeImplFilename)

		implSourceFile, err := parseSourceFileNoBody(implFilename.Value)
		if err != nil {
			return ast.SourceFile{}, err
		}

		sourceFile.Imports.IDs = append(sourceFile.Imports.IDs, implSourceFile.Imports.IDs...)
		sourceFile.Flags.Filenames = append(sourceFile.Flags.Filenames, implSourceFile.Flags.Filenames...)
		sourceFile.Validation.Join(implSourceFile.Validation)
	}

	slices.SortFunc(sourceFile.Imports.IDs, ast.CompareModuleID)
	sourceFile.Imports.IDs = slices.CompactFunc(sourceFile.Imports.IDs, func(id1, id2 ast.ModuleID) bool {
		return ast.CompareModuleID(id1, id2) == 0
	})

	slices.SortFunc(sourceFile.Flags.Filenames, ast.CompareFilename)
	sourceFile.Flags.Filenames = slices.Compact(sourceFile.Flags.Filenames)

	return sourceFile, nil
}

func New() (Querier, error) {
	finder, err := newModuleFinder(nil)
	if err != nil {
		return Querier{}, err
	}

	return Querier{finder}, nil
}

func NewWithWorkspace(workspace ast.Workspace) (Querier, error) {
	finder, err := newModuleFinder(&workspace)
	if err != nil {
		return Querier{}, err
	}

	return Querier{finder}, nil
}
