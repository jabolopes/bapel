package comp

import (
	"cmp"
	"fmt"
	"path"
	"slices"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/query"
)

type importedModule struct {
	importErr error
	// Whether this node is still being visited. This is used to detect
	// cyclic imports.
	visiting bool
}

type Resolver struct {
	querier         query.Querier
	sourceFile      ast.SourceFile
	unit            *ir.IrUnit
	importedModules map[string]importedModule
}

func (r *Resolver) resolveImport(moduleID ir.ModuleID) (retErr error) {
	if importedModule, ok := r.importedModules[moduleID.Name]; ok {
		if importedModule.visiting {
			// TODO: Include the cycle in the error message, i.e., the path
			// between the modules that forms the cycle.
			retErr = fmt.Errorf("import cycle with module %q", moduleID)
		} else {
			retErr = importedModule.importErr
		}
		return
	}

	r.importedModules[moduleID.Name] = importedModule{nil, true /* visiting */}
	defer func() {
		r.importedModules[moduleID.Name] = importedModule{retErr, false /* visiting */}
	}()

	moduleQuery, err := r.querier.QueryModuleExports(moduleID)
	if err != nil {
		return err
	}

	{
		for _, moduleID := range moduleQuery.Imports {
			if err := r.resolveImport(moduleID); err != nil {
				return err
			}
		}
	}

	r.unit.Imports = append(r.unit.Imports, ir.NewImport(moduleID))
	r.unit.ImportDecls = append(r.unit.ImportDecls, moduleQuery.Decls...)
	return nil
}

func (r *Resolver) resolveImports(imports ast.Imports) error {
	for _, moduleID := range imports.IDs {
		if err := r.resolveImport(moduleID); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveImpl(implFilename ir.Filename) error {
	sourceFileQuery, err := query.QuerySourceFile(implFilename.Value)
	if err != nil {
		return err
	}

	for _, moduleID := range sourceFileQuery.Imports {
		if err := r.resolveImport(moduleID); err != nil {
			return err
		}
	}

	r.unit.Impls = append(r.unit.Impls, ir.NewImpl(implFilename))
	r.unit.ImplDecls = append(r.unit.ImplDecls, sourceFileQuery.Decls...)
	return nil
}

func (r *Resolver) resolveImpls(relativeImplFilenames []ir.Filename) error {
	baseFilename := r.sourceFile.Header.Filename

	for _, relativeImplFilename := range relativeImplFilenames {
		implFilename := r.querier.ImplSourceFilename(baseFilename, relativeImplFilename)

		if err := r.resolveImpl(implFilename); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveImplSourceFileImpls() error {
	moduleQuery, err := r.querier.QueryModule(r.sourceFile.Header.ModuleID)
	if err != nil {
		return err
	}

	basename := path.Base(r.sourceFile.Header.Filename.Value)

	index := slices.IndexFunc(moduleQuery.Impls, func(filename ir.Filename) bool {
		return filename.Value == basename
	})
	if index == -1 {
		return fmt.Errorf("implementation file %q belongs to module %q but it's not part of the base source file `impls` section", r.sourceFile.Header.Filename, r.sourceFile.Header.ModuleID)
	}

	aboveImpls := moduleQuery.Impls[0:index]
	return r.resolveImpls(aboveImpls)
}

func (r *Resolver) resolveDecls() error {
	for _, source := range r.sourceFile.Body {
		if !source.Is(ast.DeclSource) {
			continue
		}

		c := source.Decl

		if c.Decl.Is(ir.AliasDecl) {
			// Add typename.
			decl := ir.NewNameDecl(c.Decl.ID(), c.Decl.Alias.Kind, c.Decl.Export)
			decl.Pos = c.Decl.Pos

			r.unit.Decls = append(r.unit.Decls, decl)
		}

		r.unit.Decls = append(r.unit.Decls, c.Decl)
	}

	return nil
}

func (r *Resolver) resolveFunctions() error {
	for _, source := range r.sourceFile.Body {
		if !source.Is(ast.FunctionSource) {
			continue
		}

		c := source.Function

		decl := c.Decl()

		r.unit.Decls = append(r.unit.Decls, decl)
		r.unit.Functions = append(r.unit.Functions, *c)
	}

	return nil
}

func (r *Resolver) resolve() error {
	if err := r.resolveImports(r.sourceFile.Imports); err != nil {
		return err
	}

	switch {
	case r.sourceFile.Header.Is(ast.BaseSourceFile):
		if err := r.resolveImpls(r.sourceFile.Impls.Filenames); err != nil {
			return err
		}
	case r.sourceFile.Header.Is(ast.ImplSourceFile):
		if err := r.resolveImplSourceFileImpls(); err != nil {
			return err
		}
	}

	if err := r.resolveDecls(); err != nil {
		return err
	}

	if err := r.resolveFunctions(); err != nil {
		return err
	}

	r.unit.Imports = ir.CleanImports(r.unit.Imports)

	typesBeforeTerms := func(x, y ir.IrDecl) int {
		if c := cmp.Compare(x.Case, y.Case); c != 0 {
			// Sort name decl before alias decl before term decl.
			return -c
		}

		switch {
		case x.Is(ir.TermDecl):
			return cmp.Compare(x.Term.ID, y.Term.ID)
		case x.Is(ir.AliasDecl):
			return cmp.Compare(x.Alias.ID, y.Alias.ID)
		case x.Is(ir.NameDecl):
			return cmp.Compare(x.Name.ID, y.Name.ID)
		}

		return 0
	}

	// TODO: Implement topological sorting.
	slices.SortFunc(r.unit.ImportDecls, typesBeforeTerms)
	slices.SortFunc(r.unit.ImplDecls, typesBeforeTerms)
	slices.SortFunc(r.unit.Decls, typesBeforeTerms)

	return nil
}

func ResolveSourceFile(querier query.Querier, sourceFile ast.SourceFile) (ir.IrUnit, error) {
	var c ir.IrUnitCase
	switch sourceFile.Header.Case {
	case ast.BaseSourceFile:
		c = ir.BaseUnit
	case ast.ImplSourceFile:
		c = ir.ImplUnit
	}

	unit := &ir.IrUnit{
		Case:     c,
		ModuleID: sourceFile.Header.ModuleID.Name,
		Filename: sourceFile.Header.Filename.Value,
	}

	r := &Resolver{querier, sourceFile, unit, map[string]importedModule{}}
	if err := r.resolve(); err != nil {
		return ir.IrUnit{}, err
	}

	return *unit, nil
}
