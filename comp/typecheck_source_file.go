package comp

import (
	"fmt"

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

type TypecheckOptions struct {
	// Whether to skip context initialization with the default symbols.
	SkipDefaultContext bool
	// Whether to skip function typechecking. Type inference remains
	// enabled either way.
	SkipTermTypechecker bool
	// If a module contains C++ files, we can no longer check the module for
	// declared but undefined symbols, since we can't yet inspect the C++ module.
	SkipUndefinedTermChecks bool
}

type sourceFileChecker struct {
	options TypecheckOptions
	context stlc.Context
	// Term symbols to track which symbols are declared / defined. Declared but
	// undefined terms are not allowed.
	symbols map[string]symbol
}

func (c *sourceFileChecker) addSymbol(decl ir.IrDecl) error {
	var err error
	c.context, err = c.context.AddSymbol(decl)
	return err
}

func (c *sourceFileChecker) addDecl(decl ir.IrDecl) error {
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
}

func (c *sourceFileChecker) checkFunction(function *ir.IrFunction) error {
	typechecker := stlc.NewTypechecker(c.context)

	if c.options.SkipTermTypechecker {
		var err error
		if c.context, err = typechecker.InferFunction(function); err != nil {
			return err
		}
	} else {
		if _, err := typechecker.InferFunction(function); err != nil {
			return err
		}

		var err error
		if c.context, err = typechecker.TypecheckFunction(function); err != nil {
			return err
		}
	}

	return nil
}

func (c *sourceFileChecker) addFunction(function *ir.IrFunction) error {
	decl := function.Decl()

	symbol, ok := c.symbols[decl.ID()]
	if !ok {
		symbol.decl = decl
	}

	if symbol.defined {
		return fmt.Errorf("symbol %q already defined in %v", decl.ID(), decl.Pos)
	}

	symbol.defined = true
	c.symbols[decl.ID()] = symbol

	return c.checkFunction(function)
}

func (c *sourceFileChecker) checkUnit(unit *ir.IrUnit) error {
	for _, decl := range unit.ImportDecls {
		if err := c.addSymbol(decl); err != nil {
			return err
		}
	}

	var err error
	c.context, err = c.context.EnterScope()
	if err != nil {
		return err
	}

	for _, decl := range unit.ImplDecls {
		if err := c.addSymbol(decl); err != nil {
			return err
		}
	}

	for _, decl := range unit.Decls {
		if err := c.addDecl(decl); err != nil {
			return err
		}
	}

	for i := range unit.Functions {
		if err := c.addFunction(&unit.Functions[i]); err != nil {
			return err
		}
	}

	if !c.options.SkipUndefinedTermChecks {
		for _, symbol := range c.symbols {
			if symbol.declared && !symbol.defined {
				return fmt.Errorf("%v: symbol %q is declared but it is not defined in that source file",
					symbol.decl.Pos, symbol.decl.ID())
			}
		}
	}

	return nil
}

func typecheckUnit(options TypecheckOptions, unit *ir.IrUnit) error {
	var context stlc.Context
	if options.SkipDefaultContext {
		context = stlc.NewContext()
	} else {
		var err error
		context, err = newContext()
		if err != nil {
			return err
		}
	}

	checker := &sourceFileChecker{
		options,
		context,
		map[string]symbol{},
	}

	if err := checker.checkUnit(unit); err != nil {
		return fmt.Errorf("failed to typecheck %q:\n  %v", unit.Filename, err)
	}

	return nil
}

func TypecheckSourceFile(querier query.Querier, options TypecheckOptions, inputFilename string) (ir.IrUnit, error) {
	sourceFile, err := parse.ParseSourceFile(inputFilename)
	if err != nil {
		return ir.IrUnit{}, err
	}

	if !sourceFile.Valid() {
		return ir.IrUnit{}, sourceFile.Error()
	}

	unit, err := ResolveSourceFile(querier, sourceFile)
	if err != nil {
		return ir.IrUnit{}, err
	}

	if err := typecheckUnit(options, &unit); err != nil {
		return ir.IrUnit{}, err
	}

	return unit, nil
}
