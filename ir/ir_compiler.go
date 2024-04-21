package ir

import (
	"fmt"
	"io"
	"strings"
)

func toID(id string) string {
	if strings.Contains(id, ".") {
		return "::" + strings.Replace(id, ".", "::", -1)
	}
	return id
}

type Compiler struct {
	printer *CppPrinter
	context *IrContext
}

func (a *Compiler) printf(format string, args ...any) {
	a.printer.printf(format, args...)
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

func (a *Compiler) EndModule() error {
	return a.context.checkModule()
}

func (a *Compiler) Section(id string, decls []IrDecl) error {
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
	if err := a.context.AddBind(NewDeclBind(DefSymbol, function.Decl())); err != nil {
		return err
	}

	a.context.enterFunction(function.ID, function.TypeVars, function.Args, function.Rets)

	if a.context.isExport(function.ID) {
		a.printf("export ")
	}

	var retErr error
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

		if retErr = a.Term(function.Body); retErr != nil {
			return
		}

		a.printReturn(function.ID, function.Rets)

		a.context.removeTillMarker(function.ID)
	})

	return retErr
}

func (a *Compiler) Component(component IrComponent) error {
	if _, ok := a.context.lookupBind(component.ID, FindAny); ok {
		return fmt.Errorf("name %q is already defined", component.ID)
	}

	typ := NewComponentType(component.ID, component.ElemType)

	if err := isTypeWellFormed(*a.context, typ); err != nil {
		return fmt.Errorf("component %s has an ill-formed type %s: %v", component.ID, typ, err)
	}

	if err := a.context.AddBind(NewDeclBind(DefSymbol, NewTypeDecl(typ))); err != nil {
		return err
	}

	getter := fmt.Sprintf("%s_get", component.ID)
	getterType := NewFunctionType(NewNameType("i64"), NewTupleType([]IrType{component.ElemType, NewNameType("i8")}))
	if err := a.context.AddBind(NewDeclBind(DefSymbol, NewTermDecl(getter, getterType))); err != nil {
		return err
	}

	setter := fmt.Sprintf("%s_set", component.ID)
	setterType := NewFunctionType(NewTupleType([]IrType{NewNameType("i64"), component.ElemType}), NewTupleType(nil))
	if err := a.context.AddBind(NewDeclBind(DefSymbol, NewTermDecl(setter, setterType))); err != nil {
		return err
	}

	// TODO: Use cpp_printer to print all arguments correctly.
	a.printf("ecs::StaticComponent<%s, %d> %s{};\n",
		component.ElemType, component.Length, component.ID)
	a.printf("std::pair<%s, bool> %s(int64_t entityId) { return ecs::get<%s>(&%s, entityId); }",
		component.ElemType, getter, component.ElemType, component.ID)
	a.printf("void %s(int64_t entityId, %s value) { ecs::set<%s>(&%s, entityId, value); }",
		setter, component.ElemType, component.ElemType, component.ID)

	return nil
}

func (a *Compiler) Term(term IrTerm) error {
	typechecker := NewIrTypechecker(a.context)
	if err := typechecker.TypecheckTerm(&term); err != nil {
		return err
	}

	a.printer.PrintTerm(term)
	return nil
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

	return &Compiler{NewCppPrinter(output), context}
}
