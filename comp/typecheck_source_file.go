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
		stlc.NewConstBind("f32", ir.NewTypeKind()),
		stlc.NewConstBind("f64", ir.NewTypeKind()),
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
		// 'if' type.
		stlc.NewTermDeclBind("ifthen",
			ir.ForallVars([]ir.VarKind{ir.VarKind{"a", ir.NewTypeKind()}},
				ir.NewFunctionType(
					ir.NewTupleType([]ir.IrType{ir.NewNameType("bool"), ir.Tvar("a")}),
					ir.Tvar("a")))),
		stlc.NewTermDeclBind("ifelse",
			ir.ForallVars([]ir.VarKind{ir.VarKind{"a", ir.NewTypeKind()}},
				ir.NewFunctionType(
					ir.NewTupleType([]ir.IrType{ir.NewNameType("bool"), ir.Tvar("a"), ir.Tvar("a")}),
					ir.Tvar("a")))),
		stlc.NewTermDeclBind("core::for",
			ir.ForallVars([]ir.VarKind{ir.VarKind{"a", ir.NewTypeKind()}},
				ir.NewFunctionType(
					ir.NewTupleType([]ir.IrType{
						ir.NewNameType("bool"),
						ir.NewFunctionType(ir.NewTupleType(nil), ir.Tvar("a")),
					}),
					ir.NewTupleType(nil)))),
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

	for i := range unit.TraitImpls {
		if err := c.addTraitImpl(&unit.TraitImpls[i]); err != nil {
			return err
		}
	}

	for i := range unit.Functions {
		if err := c.addFunction(&unit.Functions[i]); err != nil {
			return err
		}
	}

	for i := range unit.TraitImpls {
		if err := c.checkTraitImpl(&unit.TraitImpls[i]); err != nil {
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

func (c *sourceFileChecker) addTraitImpl(impl *ir.IrTraitImpl) error {
	var err error
	typechecker := stlc.NewTypechecker(c.context)
	reducedType := typechecker.ReduceType(impl.TypeName)
	c.context, err = c.context.AddBind(stlc.NewTraitImplBind(impl.TraitName, reducedType))
	return err
}

func (c *sourceFileChecker) checkTraitImpl(impl *ir.IrTraitImpl) error {
	// 1. Find the trait in the context.
	bind, err := c.context.GetTraitBind(impl.TraitName)
	if err != nil {
		return err
	}
	trait := bind.Trait

	// 2. Verify all trait methods are implemented.
	implementedMethods := make(map[string]ir.IrFunction)
	for _, m := range impl.Methods {
		implementedMethods[m.ID] = m
	}

	// Bind 'Self' to the target type in the context used for type checking the methods.
	methodContext, err := c.context.AddBind(stlc.NewAliasBind("Self", impl.TypeName))
	if err != nil {
		return err
	}

	for _, traitMethod := range trait.Methods {
		implMethod, ok := implementedMethods[traitMethod.ID]
		if !ok {
			return fmt.Errorf("method %q of trait %q is not implemented for %s", traitMethod.ID, impl.TraitName, impl.TypeName)
		}

		// 3. Verify signature matches.
		var expectedArgs []ir.IrType
		for _, arg := range traitMethod.Args {
			t := ir.SubstituteType(arg.Type, ir.NewNameType("Self"), impl.TypeName)
			expectedArgs = append(expectedArgs, t)
		}
		expectedRet := ir.SubstituteType(traitMethod.RetType, ir.NewNameType("Self"), impl.TypeName)
		expectedType := ir.NewFunctionType(ir.NewTupleType(expectedArgs), expectedRet)

		actualArgs := make([]ir.IrType, len(implMethod.Args))
		for j := range implMethod.Args {
			actualArgs[j] = implMethod.Args[j].Type
		}
		actualType := ir.NewFunctionType(ir.NewTupleType(actualArgs), implMethod.RetType)

		typechecker := stlc.NewTypechecker(methodContext)
		if err := typechecker.Subtype(expectedType, actualType); err != nil {
			return fmt.Errorf("method %q has type %s that does not match trait signature %s:\n  %v",
				implMethod.ID, actualType, expectedType, err)
		}
		if err := typechecker.Subtype(actualType, expectedType); err != nil {
			return fmt.Errorf("method %q has type %s that does not match trait signature %s:\n  %v",
				implMethod.ID, actualType, expectedType, err)
		}

		// 4. Type check the impl method body.
		if _, err := typechecker.InferFunction(&implMethod); err != nil {
			return err
		}
		if _, err := typechecker.TypecheckFunction(&implMethod); err != nil {
			return err
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

	unit, err := ResolveSourceFile(querier, sourceFile)
	if err != nil {
		return ir.IrUnit{}, err
	}

	if err := typecheckUnit(options, &unit); err != nil {
		return ir.IrUnit{}, err
	}

	return unit, nil
}
