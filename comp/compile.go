package comp

import (
	"fmt"
	"io"
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/stlc"
)

func newContext() (stlc.Context, error) {
	context := stlc.NewContext()

	binds := []stlc.Bind{
		stlc.NewDeclBind(stlc.ImportSymbol, ir.NewTypeDecl(ir.NewNameType("i8"))),
		stlc.NewDeclBind(stlc.ImportSymbol, ir.NewTypeDecl(ir.NewNameType("i16"))),
		stlc.NewDeclBind(stlc.ImportSymbol, ir.NewTypeDecl(ir.NewNameType("i32"))),
		stlc.NewDeclBind(stlc.ImportSymbol, ir.NewTypeDecl(ir.NewNameType("i64"))),
		stlc.NewDeclBind(stlc.ImportSymbol,
			ir.NewTermDecl("+",
				ir.NewForallType(
					[]string{"a"},
					ir.NewFunctionType(ir.NewTupleType([]ir.IrType{ir.NewVarType("a"), ir.NewVarType("a")}), ir.NewVarType("a"))))),
		stlc.NewDeclBind(stlc.ImportSymbol,
			ir.NewTermDecl("-",
				ir.NewForallType(
					[]string{"a"},
					ir.NewFunctionType(ir.NewTupleType([]ir.IrType{ir.NewVarType("a"), ir.NewVarType("a")}), ir.NewVarType("a"))))),
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
		if c.context, err = c.context.AddBind(stlc.NewDeclBind(symbol, decl)); err != nil {
			return err
		}
	}

	return c.printer.PrintModuleSection(id, decls)
}

func (c *Compiler) compileComponent(component ir.IrComponent) error {
	typ := ir.NewComponentType(component.ElemType)

	var err error
	if c.context, err = c.context.AddBind(stlc.NewDeclBind(stlc.DefSymbol, ir.NewTypeDecl(typ))); err != nil {
		return err
	}

	getterName := fmt.Sprintf("%s_get", component.ElemType)
	getterType := ir.NewFunctionType(ir.NewNameType("i64"), ir.NewTupleType([]ir.IrType{component.ElemType, ir.NewNameType("i8")}))

	if c.context, err = c.context.AddBind(stlc.NewDeclBind(stlc.DefSymbol, ir.NewTermDecl(getterName, getterType))); err != nil {
		return err
	}

	setterName := fmt.Sprintf("%s_set", component.ElemType)
	setterType := ir.NewFunctionType(ir.NewTupleType([]ir.IrType{ir.NewNameType("i64"), component.ElemType}), ir.NewTupleType(nil))
	if c.context, err = c.context.AddBind(stlc.NewDeclBind(stlc.DefSymbol, ir.NewTermDecl(setterName, setterType))); err != nil {
		return err
	}

	return c.printer.PrintComponent(component, getterName, setterName)
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

	c.printer.PrintFunction(function, c.context.IsExport(function.ID))
	return nil
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

func (c *Compiler) compileTypeDef(typ ir.IrType) error {
	decl := ir.NewTypeDecl(typ)

	var err error
	if c.context, err = c.context.AddBind(stlc.NewDeclBind(stlc.DefSymbol, decl)); err != nil {
		return err
	}

	// TODO: Should this be PrintType?
	c.printer.PrintDef(decl)
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
	case bplparser.TermSource:
		return c.compileTerm(*source.Term)
	case bplparser.TypeDefSource:
		return c.compileTypeDef(source.TypeDef.Type)
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
