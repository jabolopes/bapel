package ir

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/parser"
	"github.com/zyedidia/generic/stack"
)

func toID(id string) string {
	return strings.Replace(id, ".", "::", -1)
}

type Compiler struct {
	printer     *CppPrinter
	blocks      *stack.Stack[blockType]
	context     *IrContext
	typechecker *IrTypechecker
}

func (a *Compiler) printf(format string, args ...any) {
	a.printer.printf(format, args...)
}

func (a *Compiler) println() {
	a.printer.printf("\n")
}

func (a *Compiler) fun() *irFunction {
	return a.context.currentFunction()
}

func (a *Compiler) isFunctionBlock() bool {
	allowedBlocks := []blockType{functionBlock, ifThenBlock, ifElseBlock, elseBlock}

	for _, allowed := range allowedBlocks {
		if a.blocks.Peek() == allowed {
			return true
		}
	}

	return false
}

func (a *Compiler) printFunctionSignature(id string, args, rets []IrDecl) {
	if a.context.isExport(id) {
		a.printf("export ")
	}

	if strings.Contains(id, ".") {
		a.printf("namespace ")

		tokens := strings.Split(id, ".")
		tokens, id = tokens[:len(tokens)-1], tokens[len(tokens)-1]

		a.printf("%s", tokens[0])
		for _, token := range tokens[1:] {
			a.printf("::%s", token)
		}

		a.printf("{")
	}

	{
		// Print ret type.
		retTypes := make([]IrType, len(rets))
		for i := range rets {
			retTypes[i] = rets[i].Type
		}

		a.printer.bindPosition.Push(true)
		a.printer.printType(NewTupleType(retTypes))
		a.printer.bindPosition.Pop()
	}

	// Print id.
	a.printf(" %s(", id)

	// Print args.
	switch len(args) {
	case 0:
		break
	case 1:
		a.printer.printType(args[0].Type)
		a.printf(" %s", args[0].ID)
	default:
		a.printer.printType(args[0].Type)
		a.printf(" %s", args[0].ID)
		for _, arg := range args[1:] {
			a.printf(", ")
			a.printer.printType(arg.Type)
			a.printf(" %s", arg.ID)
		}
	}

	a.printf(")")
}

func (a *Compiler) endModule() error {
	if a.blocks.Pop() != moduleBlock {
		return errors.New("expected module block")
	}
	return a.context.checkModule()
}

func (a *Compiler) endImports() error {
	if a.blocks.Pop() != importsBlock {
		return errors.New("expected imports block")
	}
	a.println()
	return nil
}

func (a *Compiler) endExports() error {
	if a.blocks.Pop() != exportsBlock {
		return errors.New("expected exports block")
	}
	a.println()
	return nil
}

func (a *Compiler) endDecls() error {
	if a.blocks.Pop() != declsBlock {
		return errors.New("expected declarations block")
	}
	a.println()
	return nil
}

func (a *Compiler) endFunction() error {
	if a.blocks.Peek() != functionBlock {
		return errors.New("expected function block")
	}

	if err := a.Return(); err != nil {
		return err
	}

	a.printf("}\n")

	if strings.Contains(a.fun().id, ".") {
		// This function was defined inside a namespace, so close the namespace.
		a.printf("}\n")
	}

	a.blocks.Pop()
	a.context.leaveFunction()
	return nil
}

func (a *Compiler) endIf() error {
	if block := a.blocks.Pop(); block != ifThenBlock && block != ifElseBlock {
		return errors.New("expected if block")
	}

	a.printf("}\n")
	return nil
}

func (a *Compiler) endElse() error {
	if a.blocks.Pop() != elseBlock {
		return errors.New("expected else block")
	}

	a.printf("}\n")
	return nil
}

func (a *Compiler) Module() error {
	if a.blocks.Size() != 0 {
		return fmt.Errorf("Modules can only be defined at the toplevel")
	}

	a.printf("module;\n")
	a.printf("\n")
	a.printf("import <array>;\n")
	a.printf("import <cstdlib>;\n")
	a.printf("import <iostream>;\n")
	a.printf("import <tuple>;\n")
	a.printf("import <vector>;\n")
	a.printf("\n")
	a.printf("import c;\n")
	a.printf("\n")
	a.printf("export module bpl;\n")
	a.printf("\n")

	return nil
}

func (a *Compiler) Section(section string) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only start a '%s' block within a module block", section)
	}

	switch section {
	case "imports":
		a.blocks.Push(importsBlock)
		a.printf("// IMPORTS\n")
	case "decls":
		a.blocks.Push(declsBlock)
		a.printf("// HEADER\n")
	case "exports":
		a.blocks.Push(exportsBlock)
	default:
		return fmt.Errorf("unknown section %q", section)
	}

	return nil
}

func (a *Compiler) Declare(decl IrDecl) error {
	if block := a.blocks.Peek(); block != importsBlock && block != exportsBlock && block != declsBlock {
		return fmt.Errorf("declarations can occur only within an 'imports', an 'exports', or a 'decls' block.")
	}

	if _, ok := a.context.lookupSymbol(decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared or defined in this module", decl.ID)
	}

	switch a.blocks.Peek() {
	case importsBlock:
		if err := a.context.addImport(decl); err != nil {
			return err
		}
	case exportsBlock:
		if err := a.context.addExport(decl); err != nil {
			return err
		}
	case declsBlock:
		if err := a.context.addDecl(decl); err != nil {
			return err
		}
		a.printer.printDecl(decl)
		a.printf(";\n")
	}

	return nil
}

func (a *Compiler) Function(id string, args, rets []IrDecl) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	function := NewFunction(id, args, rets)
	if err := a.context.addFunction(function.decl()); err != nil {
		return err
	}

	a.blocks.Push(functionBlock)
	a.context.enterFunction(id, args, rets)

	a.printFunctionSignature(id, args, rets)
	a.printf(" {\n")

	for _, ret := range rets {
		a.printer.printType(ret.Type)
		a.printf(" %s;\n", ret.ID)
	}

	return nil
}

func (a *Compiler) Struct(id string, typ IrStructType) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if err := a.context.addStruct(NewTypeDecl(id, NewStructType(typ))); err != nil {
		return err
	}

	a.printf("struct %s {\n", id)
	for _, field := range typ.Fields {
		a.printer.printType(field.Type)
		a.printf(" %s;\n", field.Name)
	}
	a.printf("};\n")

	return nil
}

func (a *Compiler) Entity(id string) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if _, ok := a.context.lookupSymbol(id, FindAny); !ok {
		return fmt.Errorf("entity %q must have a previously defined type (e.g., struct)", id)
	}

	a.printf("ecs::StaticPool<int, struct %s, 1024> %s_entity;\n", id, id)
	return nil
}

func (a *Compiler) DefineLocal(decl IrDecl) error {
	if !a.isFunctionBlock() {
		return fmt.Errorf("can only define local variables inside a function")
	}

	if err := a.context.addLocal(decl); err != nil {
		return err
	}

	a.printer.printType(decl.Type)
	a.printf(" %s;\n", decl.ID)
	return nil
}

func (a *Compiler) Statement(args []IrTerm, ret IrTerm) error {
	isFunction := false
	var id string
	if len(args) > 0 && args[0].Case == TokenTerm && args[0].Token.Case == parser.IDToken {
		symbol, ok := a.context.lookupSymbol(args[0].Token.Text, FindAny)
		if ok && symbol.Decl.Type.Is(FunType) {
			isFunction = true
			id = args[0].Token.Text
			args = args[1:]
		}
	}

	var statement IrTerm
	if isFunction {
		// Call / assign call.
		//
		// funID [arg1 ...]
		// ret1 [ret2 ...] <- funID [arg1 ...]
		//
		// Examples:
		//   f
		//   f a b c
		//   x <- f
		//   x y <- f a b c
		statement = NewStatementTerm(NewAssignTerm(NewCallTerm(id, args), ret))
	} else {
		// x <- y
		// x <- 123
		statement = NewStatementTerm(NewAssignTerm(NewTupleTerm(args), ret))
	}

	if err := a.typechecker.CheckTerm(NewTupleType(nil), statement); err != nil {
		return err
	}

	a.printer.PrintTerm(statement)
	return nil
}

func (a *Compiler) Assign(args []parser.Token, rets []string) error {
	if !a.isFunctionBlock() {
		return errors.New("assignment / function call can only be used in a function block")
	}

	if len(args) == 0 {
		return fmt.Errorf("expected at least 1 argument; got %q", args)
	}

	switch args[0].Text {
	case "array.get":
		// ret <- array.get array index
		//
		// Examples:
		//   x <- array.get myarray 10
		args = args[1:]

		if len(rets) != 1 {
			return fmt.Errorf("expected exactly 1 return variable; got %q", rets)
		}

		id, args, err := parser.ShiftID(args)
		if err != nil {
			return err
		}

		if id.Case != parser.IDToken {
			return fmt.Errorf("expected identifier as first token; got %v", id)
		}

		index, args, err := parser.Shift(args, fmt.Errorf("expected number as second token; got %v", args))
		if err != nil {
			return err
		}

		// TODO: Check types.
		a.printf("%s = %s[%s];\n", rets[0], id.Text, index.Text)
		return nil

	case "array.set":
		// array.set array index value
		//
		// Examples:
		//   array.set myarray 10 myvalue
		args = args[1:]

		if len(rets) != 0 {
			return fmt.Errorf("expected no return variables; got %q", rets)
		}

		id, args, err := parser.ShiftID(args)
		if err != nil {
			return err
		}

		if id.Case != parser.IDToken {
			return fmt.Errorf("expected identifier as first token; got %v", id)
		}

		index, args, err := parser.Shift(args, fmt.Errorf("expected number as second token; got %v", args))
		if err != nil {
			return err
		}

		value, args, err := parser.Shift(args, fmt.Errorf("expected value as third argument; got %v", args))
		if err != nil {
			return err
		}

		// TODO: Check types.
		a.printf("%s[%s] = %s;\n", id.Text, index.Text, value.Text)
		return nil

	case "widen":
		// x <- widen y
		args = args[1:]

		if len(rets) != 1 {
			return fmt.Errorf("expected exactly 1 return variable; got %q", rets)
		}

		if len(args) != 1 {
			return fmt.Errorf("expected exactly 1 argument variable; got %q", args)
		}

		if err := a.typechecker.CheckWiden(args[0], rets[0]); err != nil {
			return err
		}

		a.printf("%s = %s;\n", toID(rets[0]), toID(args[0].Text))
		return nil
	}

	// Call / assign call.
	//
	// funID [arg1 ...]
	// ret1 [ret2 ...] <- funID [arg1 ...]
	//
	// Examples:
	//   f
	//   f a b c
	//   x <- f
	//   x y <- f a b c
	if symbol, ok := a.context.lookupSymbol(args[0].Text, FindAny); ok && symbol.Decl.Type.Is(FunType) {
		id, args, err := parser.ShiftID(args)
		if err != nil {
			return err
		}
		if id.Case != parser.IDToken {
			return fmt.Errorf("expected identifier as first token; got %v", id)
		}

		argTerms := make([]IrTerm, len(args))
		for i := range args {
			argTerms[i] = NewTokenTerm(args[i])
		}

		retTokens, err := parser.ParseTokens(rets)
		if err != nil {
			return err
		}

		retTerms := make([]IrTerm, len(retTokens))
		for i := range retTokens {
			retTerms[i] = NewTokenTerm(retTokens[i])
		}

		statement := NewStatementTerm(NewAssignTerm(NewCallTerm(id.Text, argTerms), NewTupleTerm(retTerms)))
		if err := a.typechecker.CheckTerm(NewTupleType(nil), statement); err != nil {
			return err
		}

		a.printer.PrintTerm(statement)
		return nil
	}

	if len(rets) != 1 {
		return fmt.Errorf("expected exactly 1 return variable; got %q", rets)
	}

	switch len(args) {
	case 1:
		// x <- y
		// x <- 123

		argTerms := make([]IrTerm, len(args))
		for i := range args {
			argTerms[i] = NewTokenTerm(args[i])
		}

		retTokens, err := parser.ParseTokens(rets)
		if err != nil {
			return err
		}

		retTerms := make([]IrTerm, len(retTokens))
		for i := range retTokens {
			retTerms[i] = NewTokenTerm(retTokens[i])
		}

		statement := NewStatementTerm(NewAssignTerm(NewTupleTerm(argTerms), NewTupleTerm(retTerms)))
		if err := a.typechecker.CheckTerm(NewTupleType(nil), statement); err != nil {
			return err
		}

		a.printer.PrintTerm(statement)
		return nil

	case 2:
		// x <- <unaryOp> y
		// x <- <unaryOp> 123

		// TODO: Check if argument is immediate or variable, and validate
		// variables are defined.
		//
		// TODO: Validate operation is defined.
		a.printf("%s = %s %s;\n", rets[0], args[0].Text, args[1].Text)
		return nil

	case 3:
		// x <- y   <binaryOp> z
		// x <- 123 <binaryOp> 456
		// x <- y   <binaryOp> 123
		// x <- 123 <binaryOp> y

		// TODO: Check if argument is immediate or variable, and validate
		// variables are defined.
		//
		// TODO: Validate operation is defined.
		a.printf("%s = %s %s %s;\n", rets[0], args[0].Text, args[1].Text, args[2].Text)
		return nil
	default:
		return fmt.Errorf("expected 1, 2 or 3 arguments; got %q", args)
	}
}

func (a *Compiler) Return() error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("can only be used within a function block")
	}

	a.printf("return ")

	switch rets := a.fun().rets; len(rets) {
	case 0:
		break
	case 1:
		a.printf("%s", rets[0].ID)
	default:
		a.printf("{%s", rets[0].ID)
		for _, ret := range rets[1:] {
			a.printf(", %s", ret.ID)
		}
		a.printf("}")
	}

	a.printf(";\n")
	return nil
}

func (a *Compiler) If(then bool, args []parser.Token) error {
	if !a.isFunctionBlock() {
		return errors.New("'if' can only be used in a function block")
	}

	isFunction := false
	if len(args) > 0 {
		symbol, ok := a.context.lookupSymbol(args[0].Text, FindAny)
		if ok && symbol.Decl.Type.Is(FunType) {
			isFunction = true
		}
	}

	var condition IrTerm
	if isFunction {
		id, args, err := parser.ShiftID(args)
		if err != nil {
			return err
		}
		if id.Case != parser.IDToken {
			return fmt.Errorf("expected identifier as first token; got %v", id)
		}

		argTerms := make([]IrTerm, len(args))
		for i := range args {
			argTerms[i] = NewTokenTerm(args[i])
		}

		condition = NewCallTerm(id.Text, argTerms)
	} else {
		argTerms := make([]IrTerm, len(args))
		for i := range args {
			argTerms[i] = NewTokenTerm(args[i])
		}

		condition = NewTupleTerm(argTerms)
	}

	if err := a.typechecker.CheckTerm(NewTupleType(nil), NewIfTerm(then, condition)); err != nil {
		return err
	}

	a.printf("if (")
	if !then {
		a.printf("!")
	}
	a.printer.PrintTerm(condition)
	a.printf(") {\n")

	if then {
		a.blocks.Push(ifThenBlock)
	} else {
		a.blocks.Push(ifElseBlock)
	}

	return nil
}

func (a *Compiler) Else() error {
	if a.blocks.Pop() != ifThenBlock {
		return errors.New("expected if block")
	}

	// After the opcode, put a placeholder offset to be rewritten by
	// 'endElse'.
	a.printf("} else {\n")
	a.blocks.Push(elseBlock)
	return nil
}

func (a *Compiler) End() error {
	switch block := a.blocks.Peek(); block {
	case moduleBlock:
		return a.endModule()
	case importsBlock:
		return a.endImports()
	case exportsBlock:
		return a.endExports()
	case declsBlock:
		return a.endDecls()
	case functionBlock:
		return a.endFunction()
	case ifThenBlock, ifElseBlock:
		return a.endIf()
	case elseBlock:
		return a.endElse()
	default:
		return fmt.Errorf("Unexpected block type %d", block)
	}
}

func (a *Compiler) PrintImmediate(typ IrIntType, sign Sign, value uint64) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'print immediate' can only be used in a function block")
	}

	// TODO: Handle signed and unsigned.
	a.printf("std::cout << %d << std::endl;\n", value)
	return nil
}

func (a *Compiler) PrintVar(sign Sign, id string) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'print var' can only be used in a function block")
	}

	if _, err := a.context.getDecl(id, FindAny); err != nil {
		return err
	}

	a.printf("std::cout << %s << std::endl;\n", id)
	return nil
}

func (a *Compiler) PrintStack(typ IrIntType, sign Sign) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'print stack' can only be used in a function block")
	}

	return errors.New("PrintStack is not implemented")
}

func NewCompiler(output io.Writer) *Compiler {
	context := NewIrContext()
	compiler := &Compiler{
		NewCppPrinter(output),
		stack.New[blockType](), /* blocks */
		context,
		NewIrTypechecker(context),
	}
	return compiler
}
