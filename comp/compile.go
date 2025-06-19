package comp

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/ts/stlc"
)

func newContext() (stlc.Context, error) {
	context := stlc.NewContext()

	binds := []stlc.Bind{
		stlc.NewConstBind("bool", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i8", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i16", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i32", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i64", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("void", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewTermBind("!",
			ir.Forall(
				"a", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("a"))),
			stlc.ImportSymbol),
		stlc.NewTermBind("+",
			ir.Forall(
				"a", ir.NewTypeKind(), ir.NewFunctionType(ir.NewTupleType([]ir.IrType{ir.NewVarType("a"), ir.NewVarType("a")}), ir.NewVarType("a"))),
			stlc.ImportSymbol),
		stlc.NewTermBind("-",
			ir.Forall(
				"a", ir.NewTypeKind(), ir.NewFunctionType(ir.NewTupleType([]ir.IrType{ir.NewVarType("a"), ir.NewVarType("a")}), ir.NewVarType("a"))),
			stlc.ImportSymbol),
	}

	for _, bind := range binds {
		var err error
		if context, err = context.AddBind(bind); err != nil {
			return context, err
		}
	}

	return context, nil
}

type Compiler struct {
	context stlc.Context
	module  ast.Module
	// If a module contains C++ files, we can no longer check the module for
	// declared but undefined symbols, since we can't yet inspect the C++ module.
	disableCheckModule bool
}

func (c *Compiler) compileSection(id string, decls []ir.IrDecl) error {
	var symbol stlc.Symbol
	switch id {
	case "imports":
		symbol = stlc.ImportSymbol
	case "impls":
		symbol = stlc.ImplSymbol
	case "exports":
		symbol = stlc.ExportSymbol
	case "decls":
		symbol = stlc.DeclSymbol
	default:
		return fmt.Errorf("unknown section %q", id)
	}

	for _, decl := range decls {
		var err error
		c.context, err = c.context.AddDecl(decl, symbol)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileComponent(component ir.IrComponent) error {
	var err error
	if c.context, err = c.context.AddBind(stlc.NewComponentBind(component.ElemType)); err != nil {
		return err
	}

	iteratorTypeName := fmt.Sprintf("%s_iterator", component.ElemType)
	c.context, err = c.context.AddBind(stlc.NewConstBind(iteratorTypeName, ir.NewTypeKind(), stlc.DefSymbol))
	return err
}

func (c *Compiler) compileFunction(function ir.IrFunction) error {
	typechecker := stlc.NewTypechecker(c.context)

	if _, err := typechecker.InferFunction(&function); err != nil {
		return err
	}

	var err error
	if c.context, err = typechecker.TypecheckFunction(&function); err != nil {
		return err
	}

	return nil
}

func (c *Compiler) compileImport(importModuleName ast.ID) error {
	if ext := path.Ext(importModuleName.Value); len(ext) > 0 {
		return fmt.Errorf("%s\n  module %q imports %q which should be a module name but instead it looks like a file with the extension %q",
			importModuleName.Pos, c.module.Header.Name, importModuleName.Value, ext)
	}

	importFile := fmt.Sprintf("%s.bpl", importModuleName.Value)
	input, err := os.Open(importFile)
	if err != nil {
		return err
	}
	defer input.Close()

	decls, err := query.QueryExports(importFile, input)
	if err != nil {
		return err
	}

	return c.compileSection("imports", decls)
}

func (c *Compiler) compileImports(imports ast.Imports) error {
	for _, moduleName := range imports.IDs {
		if err := c.compileImport(moduleName); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileImpls(filenames []ast.ID) error {
	var decls []ir.IrDecl
	for _, id := range filenames {
		filename := id.Value
		if path.Ext(filename) != ".bpl" {
			c.disableCheckModule = true
			continue
		}

		input, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer input.Close()

		implDecls, err := query.QueryDecls(filename, input)
		if err != nil {
			return err
		}

		decls = append(decls, implDecls...)
	}

	return c.compileSection("impls", decls)
}

func (c *Compiler) addAliasBind(decl ir.IrDecl) error {
	var err error
	c.context, err = c.context.AddAliasBind(decl)
	return err
}

func (c *Compiler) compileSource(source ast.Source) error {
	switch source.Case {
	case ast.ComponentSource:
		return c.compileComponent(*source.Component)
	case ast.FunctionSource:
		return c.compileFunction(*source.Function)
	case ast.TypeDefSource:
		return c.addAliasBind(source.TypeDef.Decl)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (c *Compiler) doDecls(sources []ast.Source) error {
	// In principle, we should be sorting the symbols using a
	// topological sorting of a graph constructed between the symbols
	// they define and their free variables.
	//
	// This topological sorting should include defined in impls files
	// also, since the module is defined by the module file + the impl
	// files, and in this language the order of types and terms within a
	// module should not matter.
	//
	// Until that is implemented, this uses a simpler sorting which is
	// to sort types before terms, and expect that the program sorted
	// types in usage order, which might not be true.
	//
	// TODO: Implement proper sorting of symbols and avoid the need for
	// forward declarations, except for mutually recursive terms and
	// maybe mutually recursive types.
	var typeDecls []ir.IrDecl
	var termDecls []ir.IrDecl

	for _, source := range sources {
		switch {
		case source.Is(ast.FunctionSource):
			termDecls = append(termDecls, source.Function.Decl())
		case source.Is(ast.TypeDefSource):
			typeDecls = append(typeDecls, source.TypeDef.Decl)
		}
	}

	return c.compileSection("decls", append(typeDecls, termDecls...))
}

func (c *Compiler) compileModule(module ast.Module) error {
	if err := c.compileImports(module.Imports); err != nil {
		return err
	}
	if err := c.compileSection("exports", module.Exports.Decls); err != nil {
		return err
	}
	if err := c.compileImpls(module.Impls.IDs); err != nil {
		return err
	}
	if err := c.doDecls(module.Body); err != nil {
		return err
	}

	for _, source := range module.Body {
		if err := c.compileSource(source); err != nil {
			return err
		}
	}

	if !c.disableCheckModule {
		if err := c.context.CheckModule(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) compileFile(filename string, input io.Reader) (ast.Module, error) {
	module, err := bplparser2.ParseFile(filename, input)
	if err != nil {
		return ast.Module{}, err
	}

	c.module = module
	if err := c.compileModule(module); err != nil {
		return ast.Module{}, err
	}

	return module, nil
}

func CompileModule(inputFilename string, input io.Reader, output io.Writer) error {
	context, err := newContext()
	if err != nil {
		return err
	}

	compiler := &Compiler{
		context,
		ast.Module{},
		false, /* disableCheckModule */
	}

	module, err := compiler.compileFile(inputFilename, input)
	if err != nil {
		return err
	}

	var moduleName string
	switch module.Header.Case {
	case ast.BaseModule:
		moduleName = module.Header.Name
	case ast.ImplModule:
		moduleName = fmt.Sprintf("%s:%s", module.Header.BaseModuleName.Value, module.Header.Name)
	}

	printer := NewCppPrinter(output, moduleName)
	printer.PrintModule(module)

	return nil
}
