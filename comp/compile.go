package comp

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/stlc"
)

func newContext() (stlc.Context, error) {
	context := stlc.NewContext()

	binds := []stlc.Bind{
		// Fundamental types.
		stlc.NewConstBind("bool", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i8", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i16", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i32", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("i64", ir.NewTypeKind(), stlc.ImportSymbol),
		stlc.NewConstBind("void", ir.NewTypeKind(), stlc.ImportSymbol),
		// Fundamental terms.
		stlc.NewTermBind("true", ir.NewNameType("bool"), stlc.ImportSymbol),
		stlc.NewTermBind("false", ir.NewNameType("bool"), stlc.ImportSymbol),
		// Operators.
		stlc.NewTermBind("!=", ir.OperatorType("!="), stlc.ImportSymbol),
		stlc.NewTermBind("==", ir.OperatorType("=="), stlc.ImportSymbol),
		stlc.NewTermBind(">", ir.OperatorType(">"), stlc.ImportSymbol),
		stlc.NewTermBind(">=", ir.OperatorType(">="), stlc.ImportSymbol),
		stlc.NewTermBind("<", ir.OperatorType("<"), stlc.ImportSymbol),
		stlc.NewTermBind("<=", ir.OperatorType("<="), stlc.ImportSymbol),
		stlc.NewTermBind("+", ir.OperatorType("+"), stlc.ImportSymbol),
		stlc.NewTermBind("-", ir.OperatorType("-"), stlc.ImportSymbol),
		stlc.NewTermBind("*", ir.OperatorType("*"), stlc.ImportSymbol),
		stlc.NewTermBind("/", ir.OperatorType("/"), stlc.ImportSymbol),
		stlc.NewTermBind("!", ir.OperatorType("!"), stlc.ImportSymbol),
	}

	for _, bind := range binds {
		var err error
		if context, err = context.AddBind(bind); err != nil {
			return context, err
		}
	}

	return context, nil
}

type symbol struct {
	decl     ir.IrDecl
	declared bool
	defined  bool
}

type Compiler struct {
	context stlc.Context
	// If a module contains C++ files, we can no longer check the module for
	// declared but undefined symbols, since we can't yet inspect the C++ module.
	disableCheckModule bool
	// Term symbols to track which symbols are declared / defined. Declared but
	// undefined terms are not allowed.
	symbols map[string]symbol
}

func (c *Compiler) addSymbol(decl ir.IrDecl, symbol stlc.Symbol) error {
	var err error
	c.context, err = c.context.AddSymbol(decl, symbol)
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
	case ast.DeclSource:
		decl := source.Decl.Decl

		if decl.Is(ir.TermDecl) {
			symbol, ok := c.symbols[decl.ID()]
			if !ok {
				symbol.decl = decl
			}

			if symbol.declared {
				return fmt.Errorf("symbol %q already declared in %v", decl.ID(), decl.Pos)
			}

			symbol.declared = true
			c.symbols[decl.ID()] = symbol
		}

		return c.addSymbol(decl, stlc.DeclSymbol)
	case ast.FunctionSource:
		decl := source.Function.Decl()

		symbol, ok := c.symbols[decl.ID()]
		if !ok {
			symbol.decl = decl
		}

		if symbol.defined {
			return fmt.Errorf("symbol %q already defined in %v", decl.ID(), decl.Pos)
		}

		symbol.defined = true
		c.symbols[decl.ID()] = symbol

		return c.compileFunction(source.Function)
	case ast.ImportSource:
		return c.addSymbol(source.Import.Decl, stlc.ImportSymbol)
	case ast.ImplSource:
		return c.addSymbol(source.Impl.Decl, stlc.ImportSymbol)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (c *Compiler) compileModule(module *ast.Module) error {
	for i := range module.Body {
		if err := c.compileSource(&module.Body[i]); err != nil {
			return err
		}
	}

	if !c.disableCheckModule {
		for _, symbol := range c.symbols {
			if symbol.declared && !symbol.defined {
				return fmt.Errorf("%v: symbol %q is declared but it is not defined in that module",
					symbol.decl.Pos, symbol.decl.ID())
			}
		}
	}

	return nil
}

func compileModule(inputFile, outputFile *os.File) error {
	module, err := bplparser2.ParseFile(inputFile.Name(), inputFile)
	if err != nil {
		return err
	}

	if err := ResolveModule(&module); err != nil {
		return err
	}

	context, err := newContext()
	if err != nil {
		return err
	}

	compiler := &Compiler{
		context,
		false, /* disableCheckModule */
		map[string]symbol{},
	}

	if err := compiler.compileModule(&module); err != nil {
		return err
	}

	return printModuleToCpp(module, outputFile)
}

func CompileBPLToCCM(inputFilename, outputFilename string) error {
	glog.V(1).Infof("Compiling %q to %q...", inputFilename, outputFilename)

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if err := compileModule(inputFile, outputFile); err != nil {
		return err
	}

	return outputFile.Close()
}
