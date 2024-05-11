package comp

import (
	"fmt"
	"io"
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/ts/stlc"
)

func newContext() (stlc.Context, error) {
	context := stlc.NewContext()

	binds := []stlc.Bind{
		stlc.NewNameBind("i8", stlc.ImportSymbol),
		stlc.NewNameBind("i16", stlc.ImportSymbol),
		stlc.NewNameBind("i32", stlc.ImportSymbol),
		stlc.NewNameBind("i64", stlc.ImportSymbol),
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
	printer *ir.CppPrinter
	context stlc.Context
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
			if c.context, err = c.context.AddBind(stlc.NewNameBind(decl.Name.ID, symbol)); err != nil {
				return err
			}
		default:
			panic(fmt.Errorf("unhandled %T %d", decl.Case, decl.Case))
		}
	}

	return c.printer.PrintModuleSection(id, decls)
}

func (c *Compiler) compileComponent(component ir.IrComponent) error {
	var err error
	if c.context, err = c.context.AddBind(stlc.NewComponentBind(component.ElemType)); err != nil {
		return err
	}

	iteratorTypeName := fmt.Sprintf("%s_iterator", component.ElemType)
	if c.context, err = c.context.AddBind(stlc.NewNameBind(iteratorTypeName, stlc.DefSymbol)); err != nil {
		return err
	}

	return c.printer.PrintComponent(component, iteratorTypeName)
}

func (c *Compiler) compileFunction(function ir.IrFunction) error {
	if err := stlc.NewInferencer(c.context).InferFunction(&function); err != nil {
		return err
	}

	typechecker := stlc.NewTypechecker(c.context)

	var err error
	if c.context, err = typechecker.TypecheckFunction(&function); err != nil {
		return err
	}

	c.printer.PrintFunction(function, function.Export)
	return nil
}

func (c *Compiler) compileImport(id string) error {
	input, err := os.Open(id)
	if err != nil {
		return err
	}
	defer input.Close()

	decls, err := query.QueryExports(input)
	if err != nil {
		return err
	}

	return c.compileSection("imports", decls)
}

func (c *Compiler) compileTerm(term ir.IrTerm) error {
	if err := stlc.NewInferencer(c.context).InferTerm(&term); err != nil {
		return err
	}

	typechecker := stlc.NewTypechecker(c.context)
	if err := typechecker.TypecheckTerm(&term); err != nil {
		return err
	}

	c.printer.PrintTerm(term)
	return nil
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
	case bplparser.ImportSource:
		return c.compileImport(*source.Import)
	case bplparser.TermSource:
		return c.compileTerm(*source.Term)
	case bplparser.TypeDefSource:
		return c.compileTypeDef(source.TypeDef.Export, source.TypeDef.Decl)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (c *Compiler) compileModule(sources []bplparser.Source) error {
	c.printer.PrintModuleTop()

	for _, source := range sources {
		if err := c.compileSource(source); err != nil {
			return err
		}
	}

	return c.context.CheckModule()
}

func (c *Compiler) compileFile(input *os.File) error {
	sources, err := bplparser.ParseFile(input)
	if err != nil {
		return err
	}

	return c.compileModule(sources)
}

func CompileFile(inputFile *os.File, output io.Writer) error {
	context, err := newContext()
	if err != nil {
		return err
	}

	compiler := &Compiler{ir.NewCppPrinter(output), context}
	return compiler.compileFile(inputFile)
}
