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

func (q Querier) BaseSourceFilename(moduleID ast.ModuleID) ast.Filename {
	return q.finder.baseSourceFilename(moduleID)
}

func (q Querier) ImplSourceFilename(baseFilename ast.Filename, relativeImplFilename ast.Filename) ast.Filename {
	return q.finder.implSourceFilename(baseFilename, relativeImplFilename)
}

func (q Querier) QueryModule(moduleID ast.ModuleID) (ModuleQuery, error) {
	baseFilename := q.finder.baseSourceFilename(moduleID)

	moduleQuery := ModuleQuery{}
	var implFilenames []ast.Filename
	{
		sourceFileQuery, err := queryDeclsBplFile(baseFilename.Value)
		if err != nil {
			return ModuleQuery{}, err
		}

		moduleQuery.Imports = append(moduleQuery.Imports, sourceFileQuery.Imports...)
		moduleQuery.Decls = append(moduleQuery.Decls, sourceFileQuery.Decls...)

		implFilenames = sourceFileQuery.Impls
	}

	for _, relativeImplFilename := range implFilenames {
		implFilename := q.finder.implSourceFilename(baseFilename, relativeImplFilename)

		implFileQuery, err := QuerySourceFile(implFilename.Value)
		if err != nil {
			return ModuleQuery{}, err
		}

		moduleQuery.Imports = append(moduleQuery.Imports, implFileQuery.Imports...)
		moduleQuery.Decls = append(moduleQuery.Decls, implFileQuery.Decls...)
	}

	slices.SortFunc(moduleQuery.Imports, ast.CompareModuleID)
	moduleQuery.Imports = slices.CompactFunc(moduleQuery.Imports, ast.EqualsModuleID)

	return moduleQuery, nil
}

// Queries all the exports of a module, recursing into the implementation files
// of the `impls` section.
//
// moduleID: identifier of the module, e.g., 'core'.
func (q Querier) QueryModuleExports(moduleID ast.ModuleID) (ModuleQuery, error) {
	moduleQuery, err := q.QueryModule(moduleID)
	if err != nil {
		return ModuleQuery{}, err
	}

	moduleQuery.Decls = slices.DeleteFunc(moduleQuery.Decls, func(decl ir.IrDecl) bool { return !decl.Export })
	return moduleQuery, nil
}

// Queries module metadata (e.g. imports, impls, flags, etc),recursing into the
// implementation files defined in the `impls` section to discover all the
// imports, all the flags, etc.
//
// The module body is not populated in the ast.SourceFile because this only returns
// module metadata.
func (q Querier) QueryModuleMetadata(moduleID ast.ModuleID) (ast.SourceFile, error) {
	baseFilename := q.finder.baseSourceFilename(moduleID)

	sourceFile, err := parseSourceFileNoBody(baseFilename.Value)
	if err != nil {
		return ast.SourceFile{}, err
	}

	for _, relativeImplFilename := range sourceFile.Impls.Filenames {
		if !strings.HasSuffix(relativeImplFilename.Value, ".bpl") {
			continue
		}

		implFilename := q.finder.implSourceFilename(baseFilename, relativeImplFilename)

		implSourceFile, err := parseSourceFileNoBody(implFilename.Value)
		if err != nil {
			return ast.SourceFile{}, err
		}

		sourceFile.Imports.IDs = append(sourceFile.Imports.IDs, implSourceFile.Imports.IDs...)
		sourceFile.Flags.Filenames = append(sourceFile.Flags.Filenames, implSourceFile.Flags.Filenames...)
	}

	slices.SortFunc(sourceFile.Imports.IDs, ast.CompareModuleID)
	sourceFile.Imports.IDs = slices.CompactFunc(sourceFile.Imports.IDs, ast.EqualsModuleID)

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
