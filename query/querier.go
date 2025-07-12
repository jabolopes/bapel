package query

import (
	"slices"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
)

type Querier struct {
	finder moduleFinder
}

func (q Querier) BaseSourceFilename(moduleID ast.ModuleID) ir.Filename {
	return q.finder.baseSourceFilename(moduleID)
}

func (q Querier) ImplSourceFilename(baseFilename ir.Filename, relativeImplFilename ir.Filename) ir.Filename {
	return q.finder.implSourceFilename(baseFilename, relativeImplFilename)
}

func (q Querier) QueryModule(moduleID ast.ModuleID) (ModuleQuery, error) {
	baseFilename := q.finder.baseSourceFilename(moduleID)

	moduleQuery := ModuleQuery{}
	{
		sourceFileQuery, err := queryDeclsBplFile(baseFilename.Value)
		if err != nil {
			return ModuleQuery{}, err
		}

		moduleQuery.Imports = append(moduleQuery.Imports, sourceFileQuery.Imports...)
		moduleQuery.Impls = sourceFileQuery.Impls
		moduleQuery.Flags = append(moduleQuery.Flags, sourceFileQuery.Flags...)
		moduleQuery.Decls = append(moduleQuery.Decls, sourceFileQuery.Decls...)
	}

	for _, relativeImplFilename := range moduleQuery.Impls {
		implFilename := q.finder.implSourceFilename(baseFilename, relativeImplFilename)

		implFileQuery, err := QuerySourceFile(implFilename.Value)
		if err != nil {
			return ModuleQuery{}, err
		}

		moduleQuery.Imports = append(moduleQuery.Imports, implFileQuery.Imports...)
		moduleQuery.Flags = append(moduleQuery.Flags, implFileQuery.Flags...)
		moduleQuery.Decls = append(moduleQuery.Decls, implFileQuery.Decls...)
	}

	slices.SortFunc(moduleQuery.Imports, ast.CompareModuleID)
	moduleQuery.Imports = slices.CompactFunc(moduleQuery.Imports, ast.EqualsModuleID)

	slices.SortFunc(moduleQuery.Flags, ir.CompareFilename)
	moduleQuery.Flags = slices.Compact(moduleQuery.Flags)

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
