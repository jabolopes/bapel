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
	output      io.Writer
	blocks      *stack.Stack[blockType]
	context     *IrContext
	typechecker *IrTypechecker
}

func (a *Compiler) out() io.Writer {
	return a.output
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

func (a *Compiler) printToken(token parser.Token) {
	switch token.Case {
	case parser.IDToken:
		fmt.Fprintf(a.out(), "%s", toID(token.Text))
	case parser.NumberToken:
		fmt.Fprintf(a.out(), "%s", token.Text)
	default:
		panic(fmt.Errorf("Unhandled token %d", token.Case))
	}
}

func (a *Compiler) printType(typ IrType) {
	switch typ.Case {
	case ArrayType:
		fmt.Fprintf(a.out(), "std::array<")
		a.printType(typ.ArrayType.ElemType)
		fmt.Fprintf(a.out(), ", %d>", typ.ArrayType.Size)
	case FunType:
		panic(fmt.Errorf("printType: Unimplemented function type"))
	case IntType:
		switch typ.IntType {
		case I8:
			fmt.Fprintf(a.out(), "char")
		case I16:
			fmt.Fprintf(a.out(), "int16_t")
		case I32:
			fmt.Fprintf(a.out(), "int32_t")
		case I64:
			fmt.Fprintf(a.out(), "int64_t")
		}
	case IDType:
		fmt.Fprintf(a.out(), "struct %s", toID(typ.IDType))
	default:
		panic(fmt.Errorf("printType: Unhandled case %d", typ.Case))
	}
}

func (a *Compiler) printDecl(decl IrDecl) {
	if decl.Type.Is(IntType) {
		a.printType(decl.Type)
		fmt.Fprintf(a.out(), " %s", decl.ID)
		return
	}

	typ := decl.Type.FunType

	// Print rets.
	switch len(typ.Rets) {
	case 0:
		fmt.Fprintf(a.out(), "void")
	case 1:
		a.printType(typ.Rets[0])
	default:
		fmt.Fprintf(a.out(), "std::tuple<")
		a.printType(typ.Rets[0])
		for _, ret := range typ.Rets[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(ret)
		}
		fmt.Fprintf(a.out(), ">")
	}

	// Print id.
	fmt.Fprintf(a.out(), " %s(", decl.ID)

	// Print args.
	switch len(typ.Args) {
	case 0:
		break
	case 1:
		a.printType(typ.Args[0])
	default:
		a.printType(typ.Args[0])
		for _, arg := range typ.Args[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(arg)
		}
	}

	fmt.Fprintf(a.out(), ")")
}

func (a *Compiler) printFunctionSignature(id string, args, rets []IrDecl) {
	if a.context.isExport(id) {
		fmt.Fprintf(a.out(), "export ")
	}

	if strings.Contains(id, ".") {
		fmt.Fprintf(a.out(), "namespace ")

		tokens := strings.Split(id, ".")
		tokens, id = tokens[:len(tokens)-1], tokens[len(tokens)-1]

		fmt.Fprintf(a.out(), "%s", tokens[0])
		for _, token := range tokens[1:] {
			fmt.Fprintf(a.out(), "::%s", token)
		}

		fmt.Fprintf(a.out(), "{")
	}

	// Print rets.
	switch len(rets) {
	case 0:
		fmt.Fprintf(a.out(), "void")
	case 1:
		a.printType(rets[0].Type)
	default:
		fmt.Fprintf(a.out(), "std::tuple<")
		a.printType(rets[0].Type)
		for _, ret := range rets[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(ret.Type)
		}
		fmt.Fprintf(a.out(), ">")
	}

	// Print id.
	fmt.Fprintf(a.out(), " %s(", id)

	// Print args.
	switch len(args) {
	case 0:
		break
	case 1:
		a.printType(args[0].Type)
		fmt.Fprintf(a.out(), " %s", args[0].ID)
	default:
		a.printType(args[0].Type)
		fmt.Fprintf(a.out(), " %s", args[0].ID)
		for _, arg := range args[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(arg.Type)
			fmt.Fprintf(a.out(), " %s", arg.ID)
		}
	}

	fmt.Fprintf(a.out(), ")")
}

func (a *Compiler) printCall(id string, args []parser.Token, rets []string) {
	switch len(rets) {
	case 0:
		break
	case 1:
		fmt.Fprintf(a.out(), "%s = ", rets[0])
	default:
		fmt.Fprintf(a.out(), "std::tie(%s", rets[0])
		for _, ret := range rets[1:] {
			fmt.Fprintf(a.out(), ", %s", ret)
		}
		fmt.Fprintf(a.out(), ") = ")
	}

	fmt.Fprintf(a.out(), "%s(", toID(id))

	switch len(args) {
	case 0:
		break
	case 1:
		fmt.Fprintf(a.out(), "%s", args[0].Text)
	default:
		fmt.Fprintf(a.out(), "%s", args[0].Text)
		for _, arg := range args[1:] {
			fmt.Fprintf(a.out(), ", %s", arg.Text)
		}
	}

	fmt.Fprintf(a.out(), ")")
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
	fmt.Fprintln(a.out())
	return nil
}

func (a *Compiler) endExports() error {
	if a.blocks.Pop() != exportsBlock {
		return errors.New("expected exports block")
	}
	fmt.Fprintln(a.out())
	return nil
}

func (a *Compiler) endDecls() error {
	if a.blocks.Pop() != declsBlock {
		return errors.New("expected declarations block")
	}
	fmt.Fprintln(a.out())
	return nil
}

func (a *Compiler) endFunction() error {
	if a.blocks.Peek() != functionBlock {
		return errors.New("expected function block")
	}

	if err := a.Return(); err != nil {
		return err
	}

	fmt.Fprintf(a.out(), "}\n")

	if strings.Contains(a.fun().id, ".") {
		// This function was defined inside a namespace, so close the namespace.
		fmt.Fprintf(a.out(), "}\n")
	}

	a.blocks.Pop()
	a.context.leaveFunction()
	return nil
}

func (a *Compiler) endIf() error {
	if block := a.blocks.Pop(); block != ifThenBlock && block != ifElseBlock {
		return errors.New("expected if block")
	}

	fmt.Fprintf(a.out(), "}\n")
	return nil
}

func (a *Compiler) endElse() error {
	if a.blocks.Pop() != elseBlock {
		return errors.New("expected else block")
	}

	fmt.Fprintf(a.out(), "}\n")
	return nil
}

func (a *Compiler) Module() error {
	if a.blocks.Size() != 0 {
		return fmt.Errorf("Modules can only be defined at the toplevel")
	}

	fmt.Fprintf(a.out(), "module;\n")
	fmt.Fprintf(a.out(), "\n")
	fmt.Fprintf(a.out(), "import <array>;\n")
	fmt.Fprintf(a.out(), "import <cstdlib>;\n")
	fmt.Fprintf(a.out(), "import <iostream>;\n")
	fmt.Fprintf(a.out(), "import <tuple>;\n")
	fmt.Fprintf(a.out(), "import <vector>;\n")
	fmt.Fprintf(a.out(), "\n")
	fmt.Fprintf(a.out(), "import c;\n")
	fmt.Fprintf(a.out(), "\n")
	fmt.Fprintf(a.out(), "export module bpl;\n")
	fmt.Fprintf(a.out(), "\n")

	return nil
}

func (a *Compiler) Section(section string) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only start a '%s' block within a module block", section)
	}

	switch section {
	case "imports":
		a.blocks.Push(importsBlock)
		fmt.Fprintf(a.out(), "// IMPORTS\n")
	case "decls":
		a.blocks.Push(declsBlock)
		fmt.Fprintf(a.out(), "// HEADER\n")
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
		a.printDecl(decl)
		fmt.Fprintf(a.out(), ";\n")
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
	fmt.Fprintf(a.out(), " {\n")

	for _, ret := range rets {
		a.printType(ret.Type)
		fmt.Fprintf(a.out(), " %s;\n", ret.ID)
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

	fmt.Fprintf(a.out(), "struct %s {\n", id)
	for _, field := range typ.Fields {
		a.printType(field.Type)
		fmt.Fprintf(a.out(), " %s;\n", field.Name)
	}
	fmt.Fprintf(a.out(), "};\n")

	return nil
}

func (a *Compiler) Entity(id string) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if _, ok := a.context.lookupSymbol(id, FindAny); !ok {
		return fmt.Errorf("entity %q must have a previously defined type (e.g., struct)", id)
	}

	fmt.Fprintf(a.out(), "ecs::StaticPool<int, struct %s, 1024> %s_entity;\n", id, id)
	return nil
}

func (a *Compiler) DefineLocal(decl IrDecl) error {
	if !a.isFunctionBlock() {
		return fmt.Errorf("can only define local variables inside a function")
	}

	if err := a.context.addLocal(decl); err != nil {
		return err
	}

	a.printType(decl.Type)
	fmt.Fprintf(a.out(), " %s;\n", decl.ID)
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
		fmt.Fprintf(a.out(), "%s = %s[%s];\n", rets[0], id.Text, index.Text)
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
		fmt.Fprintf(a.out(), "%s[%s] = %s;\n", id.Text, index.Text, value.Text)
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

		fmt.Fprintf(a.out(), "%s = %s;\n", toID(rets[0]), toID(args[0].Text))
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

		if err := a.typechecker.CheckAssign(NewCallTerm(id.Text, argTerms), NewTupleTerm(retTerms)); err != nil {
			return err
		}

		a.printCall(id.Text, args, rets)
		fmt.Fprintf(a.out(), ";\n")
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

		if err := a.typechecker.CheckAssign(NewTupleTerm(argTerms), NewTupleTerm(retTerms)); err != nil {
			return err
		}

		fmt.Fprintf(a.out(), "%s = %s;\n", rets[0], args[0].Text)
		return nil

	case 2:
		// x <- <unaryOp> y
		// x <- <unaryOp> 123

		// TODO: Check if argument is immediate or variable, and validate
		// variables are defined.
		//
		// TODO: Validate operation is defined.
		fmt.Fprintf(a.out(), "%s = %s %s;\n", rets[0], args[0].Text, args[1].Text)
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
		fmt.Fprintf(a.out(), "%s = %s %s %s;\n", rets[0], args[0].Text, args[1].Text, args[2].Text)
		return nil
	default:
		return fmt.Errorf("expected 1, 2 or 3 arguments; got %q", args)
	}
}

func (a *Compiler) Return() error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("can only be used within a function block")
	}

	fmt.Fprintf(a.out(), "return ")

	switch rets := a.fun().rets; len(rets) {
	case 0:
		break
	case 1:
		fmt.Fprintf(a.out(), "%s", rets[0].ID)
	default:
		fmt.Fprintf(a.out(), "{%s", rets[0].ID)
		for _, ret := range rets[1:] {
			fmt.Fprintf(a.out(), ", %s", ret.ID)
		}
		fmt.Fprintf(a.out(), "}")
	}

	fmt.Fprintf(a.out(), ";\n")
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

		if err := a.typechecker.CheckIf(NewCallTerm(id.Text, argTerms)); err != nil {
			return err
		}

		fmt.Fprintf(a.out(), "if (")
		if !then {
			fmt.Fprintf(a.out(), "!")
		}
		a.printCall(id.Text, args, nil /* rets */)
	} else {
		argTerms := make([]IrTerm, len(args))
		for i := range args {
			argTerms[i] = NewTokenTerm(args[i])
		}

		if err := a.typechecker.CheckIf(NewTupleTerm(argTerms)); err != nil {
			return err
		}

		fmt.Fprintf(a.out(), "if (")
		if !then {
			fmt.Fprintf(a.out(), "!")
		}
		if len(args) > 0 {
			a.printToken(args[0])
			for _, arg := range args[1:] {
				fmt.Fprintf(a.out(), " ")
				a.printToken(arg)
			}
		}
	}

	fmt.Fprintf(a.out(), ") {\n")

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
	fmt.Fprintf(a.out(), "} else {\n")
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
	fmt.Fprintf(a.out(), "std::cout << %d << std::endl;\n", value)
	return nil
}

func (a *Compiler) PrintVar(sign Sign, id string) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'print var' can only be used in a function block")
	}

	if _, err := a.context.getDecl(id, FindAny); err != nil {
		return err
	}

	fmt.Fprintf(a.out(), "std::cout << %s << std::endl;\n", id)
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
		output,
		stack.New[blockType](), /* blocks */
		context,
		NewIrTypechecker(context),
	}
	return compiler
}
