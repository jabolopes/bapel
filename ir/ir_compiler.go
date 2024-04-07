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

func functionDecl(id string, typeVars []string, args, rets []IrDecl) IrDecl {
	argTypes := make([]IrType, len(args))
	for i := range args {
		argTypes[i] = args[i].Type()
	}

	retTypes := make([]IrType, len(rets))
	for i := range rets {
		retTypes[i] = rets[i].Type()
	}

	typ := NewForallType(typeVars, NewFunctionType(NewTupleType(argTypes), NewTupleType(retTypes)))
	return NewTermDecl(id, typ)
}

type Compiler struct {
	printer     *CppPrinter
	blocks      *stack.Stack[block]
	context     *IrContext
	inferencer  *IrInferencer
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
	a.context.removeTillMarker(block.function.id)
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
	bind, ok := a.context.lookupBind(id, FindAny)
	decl := bind.Decl
	if !ok || decl.Case != TermDecl {
		return false
	}

	if decl.Type().Is(FunType) {
		return true
	}

	if decl.Type().Is(ForallType) {
		return decl.Type().Forall.Type.Is(FunType)
	}

	return false
}

func (a *Compiler) Module() error {
	if a.blocks.Size() != 0 {
		return fmt.Errorf("modules can only be defined at the toplevel")
	}

	a.printf("export module bpl;\n")
	a.printf("\n")
	a.printf("import <array>;\n")
	a.printf("import <cstdlib>;\n")
	a.printf("import <iostream>;\n")
	a.printf("import <tuple>;\n")
	a.printf("import <vector>;\n")
	a.printf("\n")
	a.printf("import c;\n")
	a.printf("\n")

	return nil
}

func (a *Compiler) Section(id string, decls []IrDecl) error {
	if a.blocks.Peek().typ != moduleBlock {
		return fmt.Errorf("can only start a '%s' block within a module block", id)
	}

	switch id {
	case "imports":
		a.blocks.Push(newBlock(importsBlock))
		a.printf("// IMPORTS\n")
	case "decls":
		a.blocks.Push(newBlock(declsBlock))
		a.printf("// HEADER\n")
	case "exports":
		a.blocks.Push(newBlock(exportsBlock))
	default:
		return fmt.Errorf("unknown section %q", id)
	}

	for _, decl := range decls {
		if err := a.Declare(decl); err != nil {
			return err
		}
	}

	return a.End()
}

func (a *Compiler) Declare(decl IrDecl) error {
	block := a.blocks.Peek().typ
	switch {
	case block == importsBlock:
		if err := a.context.AddBind(NewDeclBind(ImportSymbol, decl)); err != nil {
			return err
		}

	case block == exportsBlock:
		if err := a.context.AddBind(NewDeclBind(ExportSymbol, decl)); err != nil {
			return err
		}

	case block == declsBlock:
		if err := a.context.AddBind(NewDeclBind(DeclSymbol, decl)); err != nil {
			return err
		}
		a.printer.printDecl(decl)
		a.printf(";\n")

	case decl.Case == TypeDecl:
		if block != moduleBlock {
			return fmt.Errorf("types can only be defined in a module block")
		}
		if err := a.context.AddBind(NewDeclBind(DefSymbol, decl)); err != nil {
			return err
		}
		a.printer.PrintDef(decl)

	case decl.Case == TermDecl:
		if !a.isFunctionBlock() {
			return fmt.Errorf("terms can only be defined inside a function block")
		}
		if err := a.context.AddBind(NewDeclBind(DefSymbol, decl)); err != nil {
			return err
		}
		a.printer.PrintDef(decl)

	default:
		return fmt.Errorf("declaration / definition %s is not allowed in %s", decl, block)
	}

	return nil
}

func (a *Compiler) Function(id string, typeVars []string, args, rets []IrDecl) error {
	if a.blocks.Peek().typ != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if err := a.context.AddBind(NewDeclBind(DefSymbol, functionDecl(id, typeVars, args, rets))); err != nil {
		return err
	}

	retIDs := make([]string, len(rets))
	for i := range rets {
		retIDs[i] = rets[i].Term.ID
	}
	a.blocks.Push(newFunctionBlock(id, retIDs))
	a.context.enterFunction(id, typeVars, args, rets)

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
		// Print template type (if any).
		if len(typeVars) > 0 {
			a.printf("template <typename %s", typeVars[0])
			for _, tvar := range typeVars[1:] {
				a.printf(", typename %s", tvar)
			}
			a.printf(">")
		}
	}

	{
		// Print ret type.
		retTypes := make([]IrType, len(rets))
		for i := range rets {
			retTypes[i] = rets[i].Type()
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
		a.printer.printType(args[0].Type())
		a.printf(" %s", args[0].Term.ID)
	default:
		a.printer.printType(args[0].Type())
		a.printf(" %s", args[0].Term.ID)
		for _, arg := range args[1:] {
			a.printf(", ")
			a.printer.printType(arg.Type())
			a.printf(" %s", arg.Term.ID)
		}
	}

	a.printf(") {\n")

	for _, ret := range rets {
		a.printer.printType(ret.Type())
		a.printf(" %s;\n", ret.Term.ID)
	}

	return nil
}

func (a *Compiler) Entity(entity IrEntity) error {
	if a.blocks.Peek().typ != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if _, ok := a.context.lookupBind(entity.ID, FindAny); !ok {
		return fmt.Errorf("entity %q must have a previously defined type (e.g., struct)", entity.ID)
	}

	a.printf("ecs::StaticComponent<%s, %d> Component_%s{};\n", entity.ID, entity.Length, entity.ID)
	return nil
}

func (a *Compiler) Term(term IrTerm) error {
	if !a.isFunctionBlock() {
		return errors.New("terms can only occur within a function block")
	}

	if err := a.inferencer.Infer(&term); err != nil {
		return err
	}

	if err := a.typechecker.TypecheckTerm(&term); err != nil {
		return err
	}

	a.printer.PrintTerm(term)

	if term.Case == IfTerm {
		if term.If.Then {
			a.blocks.Push(newBlock(ifThenBlock))
		} else {
			a.blocks.Push(newBlock(ifElseBlock))
		}
	}
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

func NewCompiler(output io.Writer) *Compiler {
	context := NewIrContext()
	context.AddBind(NewDeclBind(ImportSymbol, NewTypeDecl(NewNameType("i8"))))
	context.AddBind(NewDeclBind(ImportSymbol, NewTypeDecl(NewNameType("i16"))))
	context.AddBind(NewDeclBind(ImportSymbol, NewTypeDecl(NewNameType("i32"))))
	context.AddBind(NewDeclBind(ImportSymbol, NewTypeDecl(NewNameType("i64"))))
	context.AddBind(NewDeclBind(ImportSymbol, NewTypeDecl(NewNameType("Number"))))
	context.AddBind(NewDeclBind(ImportSymbol, NewTypeDecl(NewNumberType())))
	context.AddBind(NewDeclBind(ImportSymbol,
		NewTermDecl("+",
			NewForallType(
				[]string{"a"},
				NewFunctionType(NewTupleType([]IrType{NewVarType("a"), NewVarType("a")}), NewVarType("a"))))))
	context.AddBind(NewDeclBind(ImportSymbol,
		NewTermDecl("-",
			NewForallType(
				[]string{"a"},
				NewFunctionType(NewTupleType([]IrType{NewVarType("a"), NewVarType("a")}), NewVarType("a"))))))

	compiler := &Compiler{
		NewCppPrinter(output),
		stack.New[block](), /* blocks */
		context,
		NewInferencer(context),
		NewIrTypechecker(context),
	}
	return compiler
}
