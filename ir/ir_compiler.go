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
	inferencer  *IrInferencer
	typechecker *IrTypechecker
}

func (a *Compiler) printf(format string, args ...any) {
	a.printer.printf(format, args...)
}

func (a *Compiler) endModule() error {
	if a.blocks.Pop().typ != moduleBlock {
		return errors.New("expected module block")
	}
	return a.context.checkModule()
}

func (a *Compiler) printReturn(id string, rets []IrDecl) {
	retIDs := make([]string, len(rets))
	for i := range rets {
		retIDs[i] = rets[i].Term.ID
	}

	a.printf("return ")

	switch len(retIDs) {
	case 0:
		break
	case 1:
		a.printf("%s", retIDs[0])
	default:
		a.printf("{%s", retIDs[0])
		for _, ret := range retIDs[1:] {
			a.printf(", %s", ret)
		}
		a.printf("}")
	}

	a.printf(";\n}\n")
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

	var symbol IrSymbol
	isComment := false
	switch id {
	case "imports":
		symbol = ImportSymbol
		isComment = true
		a.printf("/*\n * IMPORTS\n *\n")
	case "exports":
		a.printf("/*\n * EXPORTS\n *\n")
		isComment = true
		symbol = ExportSymbol
	case "decls":
		symbol = DeclSymbol
		a.printf("/*\n * HEADER\n */\n")
	default:
		return fmt.Errorf("unknown section %q", id)
	}

	for _, decl := range decls {
		if err := a.context.AddBind(NewDeclBind(symbol, decl)); err != nil {
			return err
		}

		if isComment {
			a.printf(" * ")
		}
		a.printer.printDecl(decl)
		a.printf("\n")
	}

	if isComment {
		a.printf("*/\n")
	}
	a.printf("\n")

	return nil
}

func (a *Compiler) TypeDefinition(typ IrType) error {
	if err := a.context.AddBind(NewDeclBind(DefSymbol, NewTypeDecl(typ))); err != nil {
		return err
	}
	// TODO: Should this be PrintType?
	a.printer.PrintDef(NewTypeDecl(typ))
	return nil
}

func (a *Compiler) Function(function IrFunction) error {
	if a.blocks.Peek().typ != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if err := a.context.AddBind(NewDeclBind(DefSymbol, function.Decl())); err != nil {
		return err
	}

	a.context.enterFunction(function.ID, function.TypeVars, function.Args, function.Rets)

	if a.context.isExport(function.ID) {
		a.printf("export ")
	}

	a.printer.printInNamespace(function.ID, func(id string) {
		{
			// Print template type (if any).
			if typeVars := function.TypeVars; len(typeVars) > 0 {
				a.printf("template <typename %s", typeVars[0])
				for _, tvar := range typeVars[1:] {
					a.printf(", typename %s", tvar)
				}
				a.printf(">")
			}
		}

		{
			// Print ret type.
			retTypes := make([]IrType, len(function.Rets))
			for i := range function.Rets {
				retTypes[i] = function.Rets[i].Type()
			}

			a.printer.withBindPosition(func() { a.printer.printType(NewTupleType(retTypes)) })
		}

		// Print id.
		a.printf(" %s(", id)

		// Print args.
		switch args := function.Args; len(args) {
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

		for _, ret := range function.Rets {
			a.printer.printType(ret.Type())
			a.printf(" %s;\n", ret.Term.ID)
		}

		// TODO: Handler error.
		a.Term(function.Body)

		a.printReturn(function.ID, function.Rets)

		a.context.removeTillMarker(function.ID)
	})

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
	if err := a.inferencer.Infer(&term); err != nil {
		return err
	}

	if err := a.typechecker.TypecheckTerm(&term); err != nil {
		return err
	}

	a.printer.PrintTerm(term)
	return nil
}

func (a *Compiler) End() error {
	switch block := a.blocks.Peek().typ; block {
	case moduleBlock:
		return a.endModule()
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
