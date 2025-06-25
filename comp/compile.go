package comp

import (
	"fmt"
	"io"
	"log"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
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

func (c *Compiler) addSymbol(decl ir.IrDecl, symbol stlc.Symbol) error {
	var err error
	c.context, err = c.context.AddSymbol(decl, symbol)
	return err
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

func (c *Compiler) compileFunction(function *ir.IrFunction) error {
	typechecker := stlc.NewTypechecker(c.context)

	if _, err := typechecker.InferFunction(function); err != nil {
		return err
	}

	var err error
	if c.context, err = typechecker.TypecheckFunction(function); err != nil {
		return err
	}

	return nil
}

func (c *Compiler) compileSource(source *ast.Source) error {
	switch source.Case {
	case ast.ComponentSource:
		return c.compileComponent(*source.Component)
	case ast.DeclSource:
		symbol := stlc.DeclSymbol
		if source.Decl.Decl.Export {
			symbol = stlc.ExportSymbol
		}
		return c.addSymbol(source.Decl.Decl, symbol)
	case ast.FunctionSource:
		return c.compileFunction(source.Function)
	case ast.ImportSource:
		return c.addSymbol(source.Import.Decl, stlc.ImportSymbol)
	case ast.ImplSource:
		return c.addSymbol(source.Impl.Decl, stlc.ImplSymbol)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (c *Compiler) compileModule() error {
	for i := range c.module.Body {
		if err := c.compileSource(&c.module.Body[i]); err != nil {
			return err
		}
	}

	// TODO: Finish.
	//
	// if !c.disableCheckModule {
	// 	if err := c.context.CheckModule(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (c *Compiler) compileFile(filename string, input io.Reader) (ast.Module, error) {
	module, err := bplparser2.ParseFile(filename, input)
	if err != nil {
		return ast.Module{}, err
	}

	if err := ResolveModule(&module); err != nil {
		return ast.Module{}, err
	}

	log.Printf("HERE %+s %q", module, module.Header.Name)

	c.module = module
	if err := c.compileModule(); err != nil {
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
