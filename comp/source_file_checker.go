package comp

import (
	"fmt"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/ts/stlc"
)

func newContext() (stlc.Context, error) {
	context := stlc.NewContext()

	binds := []stlc.Bind{
		// Fundamental types.
		stlc.NewConstBind("bool", ir.NewTypeKind()),
		stlc.NewConstBind("i8", ir.NewTypeKind()),
		stlc.NewConstBind("i16", ir.NewTypeKind()),
		stlc.NewConstBind("i32", ir.NewTypeKind()),
		stlc.NewConstBind("i64", ir.NewTypeKind()),
		stlc.NewConstBind("void", ir.NewTypeKind()),
		// Fundamental terms.
		stlc.NewTermDeclBind("true", ir.NewNameType("bool")),
		stlc.NewTermDeclBind("false", ir.NewNameType("bool")),
		// Operators.
		stlc.NewTermDeclBind("||", ir.OperatorType("||")),
		stlc.NewTermDeclBind("&&", ir.OperatorType("&&")),
		stlc.NewTermDeclBind("!=", ir.OperatorType("!=")),
		stlc.NewTermDeclBind("==", ir.OperatorType("==")),
		stlc.NewTermDeclBind(">", ir.OperatorType(">")),
		stlc.NewTermDeclBind(">=", ir.OperatorType(">=")),
		stlc.NewTermDeclBind("<", ir.OperatorType("<")),
		stlc.NewTermDeclBind("<=", ir.OperatorType("<=")),
		stlc.NewTermDeclBind("+", ir.OperatorType("+")),
		stlc.NewTermDeclBind("-", ir.OperatorType("-")),
		stlc.NewTermDeclBind("*", ir.OperatorType("*")),
		stlc.NewTermDeclBind("/", ir.OperatorType("/")),
		stlc.NewTermDeclBind("!", ir.OperatorType("!")),
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

type sourceFileChecker struct {
	context stlc.Context
	// If a module contains C++ files, we can no longer check the module for
	// declared but undefined symbols, since we can't yet inspect the C++ module.
	disableCheckSourceFile bool
	// Term symbols to track which symbols are declared / defined. Declared but
	// undefined terms are not allowed.
	symbols map[string]symbol
}

func (c *sourceFileChecker) addSymbol(decl ir.IrDecl) error {
	var err error
	c.context, err = c.context.AddSymbol(decl)
	return err
}

func (c *sourceFileChecker) checkFunction(function *ir.IrFunction) error {
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

func (c *sourceFileChecker) checkSource(source *ast.Source) error {
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

		return c.addSymbol(decl)

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

		return c.checkFunction(source.Function)

	case ast.ImportSource:
		return c.addSymbol(source.Import.Decl)

	case ast.ImplSource:
		return c.addSymbol(source.Impl.Decl)

	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (c *sourceFileChecker) checkSourceFile(sourceFile *ast.SourceFile) error {
	i := 0
	for ; i < len(sourceFile.Body); i++ {
		source := &sourceFile.Body[i]

		if !source.Is(ast.ImportSource) {
			break
		}

		if err := c.checkSource(source); err != nil {
			return err
		}
	}

	var err error
	c.context, err = c.context.EnterScope()
	if err != nil {
		return err
	}

	for ; i < len(sourceFile.Body); i++ {
		if err := c.checkSource(&sourceFile.Body[i]); err != nil {
			return err
		}
	}

	if !c.disableCheckSourceFile {
		for _, symbol := range c.symbols {
			if symbol.declared && !symbol.defined {
				return fmt.Errorf("%v: symbol %q is declared but it is not defined in that source file",
					symbol.decl.Pos, symbol.decl.ID())
			}
		}
	}

	return nil
}

func CheckSourceFile(querier query.Querier, inputFilename string) (ast.SourceFile, error) {
	sourceFile, err := parse.ParseSourceFile(inputFilename)
	if err != nil {
		return ast.SourceFile{}, err
	}

	if err := ResolveSourceFile(querier, &sourceFile); err != nil {
		return ast.SourceFile{}, err
	}

	context, err := newContext()
	if err != nil {
		return ast.SourceFile{}, err
	}

	checker := &sourceFileChecker{
		context,
		false, /* disableCheckSourceFile */
		map[string]symbol{},
	}

	if err := checker.checkSourceFile(&sourceFile); err != nil {
		return ast.SourceFile{}, err
	}

	return sourceFile, nil
}
