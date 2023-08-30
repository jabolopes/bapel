package ir

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/zyedidia/generic/stack"
)

func toID(id string) string {
	return strings.Replace(id, ".", "::", -1)
}

type Compiler struct {
	printer     *CppPrinter
	blocks      *stack.Stack[block]
	context     *IrContext
	typechecker *IrTypechecker
}

func (a *Compiler) printf(format string, args ...any) {
	a.printer.printf(format, args...)
}

func (a *Compiler) println() {
	a.printer.printf("\n")
}

func (a *Compiler) isFunctionBlock() bool {
	allowedBlocks := []blockType{functionBlock, ifThenBlock, ifElseBlock, elseBlock}

	for _, allowed := range allowedBlocks {
		if a.blocks.Peek().typ == allowed {
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

		a.printer.withBindPosition(func() { a.printer.printType(NewTupleType(retTypes)) })
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
	if a.blocks.Pop().typ != moduleBlock {
		return errors.New("expected module block")
	}
	return a.context.checkModule()
}

func (a *Compiler) endImports() error {
	if a.blocks.Pop().typ != importsBlock {
		return errors.New("expected imports block")
	}
	a.println()
	return nil
}

func (a *Compiler) endExports() error {
	if a.blocks.Pop().typ != exportsBlock {
		return errors.New("expected exports block")
	}
	a.println()
	return nil
}

func (a *Compiler) endDecls() error {
	if a.blocks.Pop().typ != declsBlock {
		return errors.New("expected declarations block")
	}
	a.println()
	return nil
}

func (a *Compiler) endFunction() error {
	block := a.blocks.Peek()
	if block.typ != functionBlock {
		return errors.New("expected function block")
	}

	if err := a.Return(); err != nil {
		return err
	}

	a.printf("}\n")

	if strings.Contains(block.function.id, ".") {
		// This function was defined inside a namespace, so close the namespace.
		a.printf("}\n")
	}

	a.blocks.Pop()
	a.context.leaveFunction(block.function.id)
	return nil
}

func (a *Compiler) endIf() error {
	if block := a.blocks.Pop().typ; block != ifThenBlock && block != ifElseBlock {
		return errors.New("expected if block")
	}

	a.printf("}\n")
	return nil
}

func (a *Compiler) endElse() error {
	if a.blocks.Pop().typ != elseBlock {
		return errors.New("expected else block")
	}

	a.printf("}\n")
	return nil
}

func (a *Compiler) IsFunction(id string) bool {
	symbol, ok := a.context.lookupSymbol(id, FindAny)
	return ok && symbol.Decl.Type.Is(FunType)
}

func (a *Compiler) Module() error {
	if a.blocks.Size() != 0 {
		return fmt.Errorf("modules can only be defined at the toplevel")
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
	if a.blocks.Peek().typ != moduleBlock {
		return fmt.Errorf("can only start a '%s' block within a module block", section)
	}

	switch section {
	case "imports":
		a.blocks.Push(newBlock(importsBlock))
		a.printf("// IMPORTS\n")
	case "decls":
		a.blocks.Push(newBlock(declsBlock))
		a.printf("// HEADER\n")
	case "exports":
		a.blocks.Push(newBlock(exportsBlock))
	default:
		return fmt.Errorf("unknown section %q", section)
	}

	return nil
}

func (a *Compiler) Declare(decl IrDecl) error {
	if block := a.blocks.Peek().typ; block != importsBlock && block != exportsBlock && block != declsBlock {
		return fmt.Errorf("declarations can occur only within an 'imports', an 'exports', or a 'decls' block")
	}

	if _, ok := a.context.lookupSymbol(decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared or defined in this module", decl.ID)
	}

	switch a.blocks.Peek().typ {
	case importsBlock:
		if err := a.context.addDeclaration(NewSymbol(ImportSymbol, decl)); err != nil {
			return err
		}
	case exportsBlock:
		if err := a.context.addDeclaration(NewSymbol(ExportSymbol, decl)); err != nil {
			return err
		}
	case declsBlock:
		if err := a.context.addDeclaration(NewSymbol(DeclSymbol, decl)); err != nil {
			return err
		}
		a.printer.printDecl(decl)
		a.printf(";\n")
	}

	return nil
}

func functionDecl(id string, args, rets []IrDecl) IrDecl {
	argTypes := make([]IrType, len(args))
	for i := range args {
		argTypes[i] = args[i].Type
	}

	retTypes := make([]IrType, len(rets))
	for i := range rets {
		retTypes[i] = rets[i].Type
	}

	return NewTermDecl(id, NewFunctionType(IrFunctionType{argTypes, retTypes}))
}

func (a *Compiler) Function(id string, args, rets []IrDecl) error {
	if a.blocks.Peek().typ != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if err := a.context.addDefinition(functionDecl(id, args, rets)); err != nil {
		return err
	}

	retIDs := make([]string, len(rets))
	for i := range rets {
		retIDs[i] = rets[i].ID
	}
	a.blocks.Push(newFunctionBlock(id, retIDs))
	a.context.enterFunction(id, args, rets)

	a.printFunctionSignature(id, args, rets)
	a.printf(" {\n")

	for _, ret := range rets {
		a.printer.printType(ret.Type)
		a.printf(" %s;\n", ret.ID)
	}

	return nil
}

func (a *Compiler) Define(decl IrDecl) error {
	switch decl.Case {
	case TypeDecl:
		if a.blocks.Peek().typ != moduleBlock {
			return fmt.Errorf("types can only be defined in a module block")
		}

	case TermDecl:
		if !a.isFunctionBlock() {
			return fmt.Errorf("terms can only be defined inside a function block")
		}

	default:
		panic(fmt.Errorf("unhandled decl case %d", decl.Case))
	}

	if err := a.context.addDefinition(decl); err != nil {
		return err
	}

	a.printer.PrintDef(decl)
	return nil
}

func (a *Compiler) Entity(id string) error {
	if a.blocks.Peek().typ != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if _, ok := a.context.lookupSymbol(id, FindAny); !ok {
		return fmt.Errorf("entity %q must have a previously defined type (e.g., struct)", id)
	}

	a.printf("ecs::StaticComponent<%s, 1024> Component_%s{};\n", id, id)
	return nil
}

func (a *Compiler) Statement(statement IrTerm) error {
	if err := a.typechecker.TypecheckTerm(statement); err != nil {
		return err
	}

	a.printer.PrintTerm(statement)
	return nil
}

func (a *Compiler) Return() error {
	if a.blocks.Peek().typ != functionBlock {
		return fmt.Errorf("can only be used within a function block")
	}

	a.printf("return ")

	switch rets := a.blocks.Peek().function.retIDs; len(rets) {
	case 0:
		break
	case 1:
		a.printf("%s", rets[0])
	default:
		a.printf("{%s", rets[0])
		for _, ret := range rets[1:] {
			a.printf(", %s", ret)
		}
		a.printf("}")
	}

	a.printf(";\n")
	return nil
}

func (a *Compiler) If(ifTerm IrTerm) error {
	if ifTerm.Case != IfTerm {
		panic(fmt.Errorf("expected IfTerm; got %d", ifTerm.Case))
	}

	if !a.isFunctionBlock() {
		return errors.New("'if' can only be used in a function block")
	}

	if err := a.typechecker.TypecheckTerm(ifTerm); err != nil {
		return err
	}

	a.printer.PrintTerm(ifTerm)
	if ifTerm.If.Then {
		a.blocks.Push(newBlock(ifThenBlock))
	} else {
		a.blocks.Push(newBlock(ifElseBlock))
	}

	return nil
}

func (a *Compiler) Else() error {
	if a.blocks.Pop().typ != ifThenBlock {
		return errors.New("expected if block")
	}

	// After the opcode, put a placeholder offset to be rewritten by
	// 'endElse'.
	a.printf("} else {\n")
	a.blocks.Push(newBlock(elseBlock))
	return nil
}

func (a *Compiler) End() error {
	switch block := a.blocks.Peek().typ; block {
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
		return fmt.Errorf("unexpected block type %d", block)
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
		stack.New[block](), /* blocks */
		context,
		NewIrTypechecker(context),
	}
	return compiler
}
