package comp

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/ts/stlc"
)

func newContext() (stlc.Context, error) {
	context := stlc.NewContext()

	binds := []stlc.Bind{
		stlc.NewConstBind("i8", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i16", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i32", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i64", ir.NewTypeKind(), stlc.ImportSymbol),
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
	printer    *ir.CppPrinter
	context    stlc.Context
	moduleName string
	// Whether the file being compiled is an implementation source file
	// instead of a module file.
	isImplFile bool
	// If a module contains C++ files, we can no longer check the module for
	// declared but undefined symbols, since we can't yet inspect the C++ module.
	disableCheckModule bool
}

func (c *Compiler) compileSection(id string, decls []ir.IrDecl) error {
	var symbol stlc.Symbol
	switch id {
	case "imports":
		symbol = stlc.ImportSymbol
	case "exports":
		symbol = stlc.ExportSymbol
	case "decls":
		symbol = stlc.DeclSymbol
	default:
		return fmt.Errorf("unknown section %q", id)
	}

	for _, decl := range decls {
		var err error
		switch decl.Case {
		case ir.TermDecl:
			if c.context, err = c.context.AddBind(stlc.NewTermBind(decl.Term.ID, decl.Term.Type, symbol)); err != nil {
				return err
			}
		case ir.AliasDecl:
			if c.context, err = c.context.AddBind(stlc.NewAliasBind(decl.Alias.ID, decl.Alias.Type, symbol)); err != nil {
				return err
			}
		case ir.NameDecl:
			if c.context, err = c.context.AddBind(stlc.NewConstBind(decl.Name.ID, decl.Name.Kind, symbol)); err != nil {
				return err
			}
		default:
			panic(fmt.Errorf("unhandled %T %d", decl.Case, decl.Case))
		}
	}

	c.printer.PrintModuleSection(id, decls)
	return nil
}

func (c *Compiler) compileComponent(component ir.IrComponent) error {
	var err error
	if c.context, err = c.context.AddBind(stlc.NewComponentBind(component.ElemType)); err != nil {
		return err
	}

	iteratorTypeName := fmt.Sprintf("%s_iterator", component.ElemType)
	if c.context, err = c.context.AddBind(stlc.NewConstBind(iteratorTypeName, ir.NewTypeKind(), stlc.DefSymbol)); err != nil {
		return err
	}

	return c.printer.PrintComponent(component, iteratorTypeName)
}

func (c *Compiler) compileFunction(function ir.IrFunction) error {
	typechecker := stlc.NewTypechecker(c.context)

	if err := typechecker.InferFunction(&function); err != nil {
		return err
	}

	var err error
	if c.context, err = typechecker.TypecheckFunction(&function); err != nil {
		return err
	}

	c.printer.PrintFunction(function, function.Export)
	return nil
}

func (c *Compiler) compileImport(importModuleName string) error {
	if ext := path.Ext(importModuleName); len(ext) > 0 {
		return fmt.Errorf("module %q imports %q which should be a module name but instead it looks like a file with the extension %q",
			c.moduleName, importModuleName, ext)
	}

	importFile := fmt.Sprintf("%s.bpl", importModuleName)
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

func (c *Compiler) compileImports(modules []string) error {
	for _, moduleName := range modules {
		if err := c.compileImport(moduleName); err != nil {
			return err
		}
	}
	c.printer.PrintImportsSection(modules)
	return nil
}

func (c *Compiler) compileImpls(filenames []string) error {
	impls := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		if path.Ext(filename) == ".cpp" {
			c.disableCheckModule = true
		}

		impls = append(impls, TrimExtension(filename))
	}

	return c.printer.PrintImpls(c.moduleName, impls)
}

func (c *Compiler) compileTypeDef(export bool, decl ir.IrDecl) error {
	var err error
	if c.context, err = c.context.AddBind(stlc.NewAliasBind(decl.Alias.ID, decl.Alias.Type, stlc.DefSymbol)); err != nil {
		return err
	}

	c.printer.PrintDecl(decl, export)
	return nil
}

func (c *Compiler) compileSource(source bplparser.Source) error {
	switch source.Case {
	case bplparser.SectionSource:
		return c.compileSection(source.Section.ID, source.Section.Decls)
	case bplparser.ComponentSource:
		return c.compileComponent(*source.Component)
	case bplparser.FunctionSource:
		return c.compileFunction(*source.Function)
	case bplparser.ImportsSource:
		return c.compileImports(source.Imports.IDs)
	case bplparser.ImplsSource:
		return c.compileImpls(source.Impls.IDs)
	case bplparser.TypeDefSource:
		return c.compileTypeDef(source.TypeDef.Export, source.TypeDef.Decl)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (c *Compiler) compileModule(sources []bplparser.Source) error {
	c.printer.PrintModuleTop(c.moduleName)

	for _, source := range sources {
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

func (c *Compiler) compileFile(filename string, input io.Reader) error {
	sources, err := bplparser2.ParseFile(filename, input)
	if err != nil {
		return err
	}

	return c.compileModule(sources)
}

func CompileModuleFile(inputFilename string, input io.Reader, output io.Writer) error {
	context, err := newContext()
	if err != nil {
		return err
	}

	moduleName := TrimExtension(path.Base(inputFilename))

	compiler := &Compiler{
		ir.NewCppPrinter(output),
		context,
		moduleName,
		false, /* isImplFile */
		false, /* disableCheckModule */
	}
	return compiler.compileFile(inputFilename, input)
}

func CompileImplFile(inputFilename, moduleName string, input io.Reader, output io.Writer) error {
	context, err := newContext()
	if err != nil {
		return err
	}

	implModuleName := fmt.Sprintf("%s:%s", moduleName, TrimExtension(path.Base(inputFilename)))

	compiler := &Compiler{
		ir.NewCppPrinter(output),
		context,
		implModuleName,
		true,  /* isImplFile */
		false, /* disableCheckModule */
	}
	return compiler.compileFile(inputFilename, input)
}
