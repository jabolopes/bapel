package comp

import (
	"fmt"
	"path"
	"slices"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/query"
)

type Resolver struct {
	module *ast.Module
}

func (r *Resolver) resolveImport(moduleID ast.ID) ([]ast.Source, error) {
	if ext := path.Ext(moduleID.Value); len(ext) > 0 {
		return nil, fmt.Errorf("%s\n  module ID %q looks like a file with the extension %q; it should be a module ID (without extension)",
			moduleID.Pos, moduleID.Value, ext)
	}

	decls, err := query.QueryModuleExports(moduleID.Value)
	if err != nil {
		return nil, err
	}

	sources := make([]ast.Source, 0, len(decls))
	for _, decl := range decls {
		sources = append(sources, ast.NewImportSource(moduleID, decl))
	}

	return sources, nil
}

func (r *Resolver) resolveImports(imports ast.Imports) ([]ast.Source, error) {
	var allSources []ast.Source
	for _, moduleID := range imports.IDs {
		sources, err := r.resolveImport(moduleID)
		if err != nil {
			return nil, err
		}
		allSources = append(allSources, sources...)
	}
	return allSources, nil
}

func (r *Resolver) resolveImpl(filename ast.ID) ([]ast.Source, error) {
	decls, err := query.QueryFileDecls(filename.Value)
	if err != nil {
		return nil, err
	}

	sources := make([]ast.Source, 0, len(decls))
	for _, decl := range decls {
		var source ast.Source
		if decl.Export {
			source = ast.NewDeclSource(decl)
		} else {
			source = ast.NewImplSource(filename.Value, decl)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *Resolver) resolveImpls(filenames []ast.ID) ([]ast.Source, error) {
	var allSources []ast.Source
	for _, filename := range filenames {
		sources, err := r.resolveImpl(filename)
		if err != nil {
			return nil, err
		}

		allSources = append(allSources, sources...)
	}

	return allSources, nil
}

func (r *Resolver) resolveDecls() ([]ast.Source, error) {
	var sources []ast.Source
	for _, source := range r.module.Body {
		if !source.Is(ast.DeclSource) {
			continue
		}

		c := source.Decl

		if !c.Decl.Is(ir.AliasDecl) {
			continue
		}

		decl := ir.NewNameDecl(c.Decl.ID(), c.Decl.Alias.Kind, c.Decl.Export)
		decl.Pos = c.Decl.Pos

		sources = append(sources, ast.NewDeclSource(decl))
	}

	return sources, nil
}

func (r *Resolver) resolveFunctions() ([]ast.Source, error) {
	var sources []ast.Source
	for _, source := range r.module.Body {
		if !source.Is(ast.FunctionSource) {
			continue
		}

		c := source.Function

		decl := c.Decl()

		sources = append(sources, ast.NewDeclSource(decl))
	}

	return sources, nil
}

func (r *Resolver) resolve() error {
	importSources, err := r.resolveImports(r.module.Imports)
	if err != nil {
		return err
	}

	implSources, err := r.resolveImpls(r.module.Impls.IDs)
	if err != nil {
		return err
	}

	declSources, err := r.resolveDecls()
	if err != nil {
		return err
	}

	moreDeclSources, err := r.resolveFunctions()
	if err != nil {
		return err
	}

	typesBeforeTerms := func(x, y ast.Source) int {
		if x.Case != y.Case {
			return 0
		}

		var declX ir.IrDecl
		var declY ir.IrDecl

		switch {
		case x.Is(ast.DeclSource):
			declX = x.Decl.Decl
		case x.Is(ast.ImportSource):
			declX = x.Import.Decl
		case x.Is(ast.ImplSource):
			declX = x.Impl.Decl
		}

		switch {
		case y.Is(ast.DeclSource):
			declY = y.Decl.Decl
		case y.Is(ast.ImportSource):
			declY = y.Import.Decl
		case y.Is(ast.ImplSource):
			declY = x.Impl.Decl
		}

		typeX := declX.Is(ir.NameDecl) || declX.Is(ir.AliasDecl)
		typeY := declY.Is(ir.NameDecl) || declY.Is(ir.AliasDecl)
		if typeX == typeY {
			return 0
		}

		if typeX {
			return -1
		}

		return 1
	}

	// TODO: Implement topological sorting.
	slices.SortFunc(importSources, typesBeforeTerms)
	slices.SortFunc(implSources, typesBeforeTerms)
	slices.SortFunc(declSources, typesBeforeTerms)

	r.module.Body =
		append(importSources,
			append(implSources,
				append(declSources,
					append(moreDeclSources,
						r.module.Body...)...)...)...)

	return nil
}

func ResolveModule(module *ast.Module) error {
	r := &Resolver{module}
	return r.resolve()
}
