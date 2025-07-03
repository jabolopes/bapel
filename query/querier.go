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

func (q Querier) ModuleBaseFilename(moduleID ast.ModuleID) ast.Filename {
	return q.finder.moduleBaseFilename(moduleID)
}

func (q Querier) ModuleImplFilename(baseFilename ast.Filename, relativeImplFilename ast.Filename) ast.Filename {
	return q.finder.moduleImplFilename(baseFilename, relativeImplFilename)
}

func (q Querier) QueryModuleDecls(moduleID ast.ModuleID) ([]ir.IrDecl, error) {
	baseFilename := q.finder.moduleBaseFilename(moduleID)

	module, decls, err := queryDeclsBplFile(baseFilename.Value)
	if err != nil {
		return nil, err
	}

	for _, relativeImplFilename := range module.Impls.Filenames {
		implFilename := q.finder.moduleImplFilename(baseFilename, relativeImplFilename)

		implDecls, err := QueryFileDecls(implFilename.Value)
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
// The module body is not populated in the ast.Module because this only returns
// module metadata.
func (q Querier) QueryModuleMetadata(moduleID ast.ModuleID) (ast.Module, error) {
	baseFilename := q.finder.moduleBaseFilename(moduleID)

	module, err := parseModuleNoBody(baseFilename.Value)
	if err != nil {
		return ast.Module{}, err
	}

	for _, relativeImplFilename := range module.Impls.Filenames {
		if !strings.HasSuffix(relativeImplFilename.Value, ".bpl") {
			continue
		}

		implFilename := q.finder.moduleImplFilename(baseFilename, relativeImplFilename)

		implModule, err := parseModuleNoBody(implFilename.Value)
		if err != nil {
			return ast.Module{}, err
		}

		module.Imports.IDs = append(module.Imports.IDs, implModule.Imports.IDs...)
		module.Flags.Filenames = append(module.Flags.Filenames, implModule.Flags.Filenames...)
		module.Validation.Join(implModule.Validation)
	}

	slices.SortFunc(module.Imports.IDs, ast.CompareModuleID)
	module.Imports.IDs = slices.CompactFunc(module.Imports.IDs, func(id1, id2 ast.ModuleID) bool {
		return ast.CompareModuleID(id1, id2) == 0
	})

	slices.SortFunc(module.Flags.Filenames, ast.CompareFilename)
	module.Flags.Filenames = slices.Compact(module.Flags.Filenames)

	return module, nil
}

func New() (Querier, error) {
	finder, err := newModuleFinder()
	if err != nil {
		return Querier{}, err
	}

	return Querier{finder}, nil
}
