package ir

import (
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/parser"
)

type Position int

const (
	TypePosition = Position(iota)
	BindPosition
)

func toID(id string) string {
	if strings.Contains(id, ".") {
		return "::" + strings.Replace(id, ".", "::", -1)
	}
	return id
}

type CppPrinter struct {
	output   io.Writer
	position Position
}

func (p *CppPrinter) withBindPosition(callback func()) {
	position := p.position
	p.position = BindPosition
	defer func() { p.position = position }()
	callback()
}

func (p *CppPrinter) printInNamespace(id string, callback func(string)) {
	if !strings.Contains(id, ".") {
		callback(id)
		return
	}

	p.printf("namespace ")

	tokens := strings.Split(id, ".")
	tokens, id = tokens[:len(tokens)-1], tokens[len(tokens)-1]

	p.printf("%s", tokens[0])
	for _, token := range tokens[1:] {
		p.printf("::%s", token)
	}

	p.printf(" { ")
	callback(id)
	p.printf(" }")
}

func (p *CppPrinter) printf(format string, args ...any) {
	fmt.Fprintf(p.output, format, args...)
}

func (p *CppPrinter) printCall(id IrTerm, types []IrType, arg IrTerm) {
	p.PrintTerm(id)
	if id.Is(TokenTerm) && !IsOperator(id.Token.Text) && len(types) > 0 {
		p.printf("<")
		p.withBindPosition(func() {
			p.printType(types[0])
			for _, typ := range types[1:] {
				p.printf(", ")
				p.printType(typ)
			}
		})
		p.printf(">")
	}
	p.printf("(")
	p.PrintTerm(arg)
	p.printf(")")
}

func (a *CppPrinter) printReturn(id string, retIDs []string) {
	a.printf("return")

	switch len(retIDs) {
	case 0:
		break
	case 1:
		a.printf(" %s", retIDs[0])
	default:
		a.printf(" {%s", retIDs[0])
		for _, ret := range retIDs[1:] {
			a.printf(", %s", ret)
		}
		a.printf("}")
	}

	a.printf(";\n")
}

func (p *CppPrinter) printToken(token parser.Token) {
	switch token.Case {
	case parser.IDToken:
		p.printf("%s", toID(token.Text))
	case parser.NumberToken:
		p.printf("%s", token.Text)
	default:
		panic(fmt.Errorf("unhandled %d %d", token.Case, token.Case))
	}
}

func (p *CppPrinter) printAliasDecl(id string, value IrType) {
	switch value.Case {
	case StructType:
		p.printf("struct %s {\n", id)
		for _, field := range value.Fields() {
			p.printType(field.Type)
			p.printf(" %s;\n", field.ID)
		}
		p.printf("};\n")

	default:
		panic(fmt.Errorf("unhandled %T %d", value.Case, value.Case))
	}
}

func (p *CppPrinter) printType(typ IrType) {
	switch {
	case typ.Case == ArrayType:
		fmt.Fprintf(p.out(), "std::array<")
		p.printType(typ.Array.ElemType)
		p.printf(", %d>", typ.Array.Size)

	case typ.Is(ForallType):
		tvars := typ.ForallVars()
		p.printf("template <typename %s", tvars[0])
		for _, tvar := range tvars[1:] {
			p.printf(", typename %s", tvar)
		}
		p.printf("> ")
		p.printType(typ.ForallBody())

	case typ.Is(NameType):
		switch typ.Name {
		case "i8":
			p.printf("int8_t")
		case "i16":
			p.printf("int16_t")
		case "i32":
			p.printf("int32_t")
		case "i64":
			p.printf("int64_t")
		default:
			p.printf("%s", toID(typ.Name))
		}

	case typ.Is(TupleType) && p.position == TypePosition:
		tuple := typ.Tuple
		if len(tuple) > 0 {
			p.printType(tuple[0])
			for _, elem := range tuple[1:] {
				p.printf(", ")
				p.printType(elem)
			}
		}

	case typ.Is(TupleType) && p.position == BindPosition:
		tuple := typ.Tuple
		// Print rets.
		switch len(tuple) {
		case 0:
			p.printf("void")
		case 1:
			p.printType(tuple[0])
		default:
			p.printf("std::tuple<")
			p.printType(tuple[0])
			for _, elem := range tuple[1:] {
				p.printf(", ")
				p.printType(elem)
			}
			p.printf(">")
		}

	case typ.Is(VarType):
		p.printf("%s", typ.Var)

	default:
		panic(fmt.Errorf("printType: unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) PrintDecl(decl IrDecl) {
	if decl.Is(NameDecl) {
		p.printType(NewNameType(decl.Name.ID))
		return
	}

	if decl.Is(AliasDecl) {
		p.printAliasDecl(decl.Alias.ID, decl.Alias.Type)
		return
	}

	switch typ := decl.Term.Type; typ.Case {
	case ForallType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			tvars := typ.ForallVars()
			p.printf("template <typename %s", tvars[0])
			for _, tvar := range tvars[1:] {
				p.printf(", typename %s", tvar)
			}
			p.printf("> ")
			p.PrintDecl(NewTermDecl(id, typ.ForallBody()))
		})

	case FunType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.withBindPosition(func() { p.printType(typ.Fun.Ret) })
			p.printf(" %s(", id)
			p.printType(typ.Fun.Arg)
			p.printf(");")
		})

	case NameType:
		p.printType(typ)
		p.printf(" %s", decl.Term.ID)

	case StructType:
		// TODO: Handle namespacing.
		p.printf("struct %s", decl.Term.ID)

	case TupleType:
		c := typ.Tuple
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.printf("std::tuple<")
			p.printType(c[0])
			for _, typ := range c[1:] {
				p.printf(", ")
				p.printType(typ)
			}
			p.printf("> %s", id)
		})

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) PrintModuleTop() {
	p.printf("export module bpl;\n")
	p.printf("\n")
	p.printf("import <array>;\n")
	p.printf("import <cstdlib>;\n")
	p.printf("import <iostream>;\n")
	p.printf("import <tuple>;\n")
	p.printf("import <vector>;\n")
	p.printf("\n")
	p.printf("import c;\n")
	p.printf("\n")
}

func (p *CppPrinter) PrintModuleSection(id string, decls []IrDecl) error {
	isComment := false
	switch id {
	case "imports":
		isComment = true
		p.printf("/*\n * IMPORTS\n *\n")
	case "exports":
		p.printf("/*\n * EXPORTS\n *\n")
		isComment = true
	case "decls":
		p.printf("/*\n * HEADER\n */\n")
	default:
		return fmt.Errorf("unknown section %q", id)
	}

	for _, decl := range decls {
		if isComment {
			p.printf(" * ")
		}
		p.PrintDecl(decl)
		p.printf("\n")
	}

	if isComment {
		p.printf("*/\n")
	}
	p.printf("\n")

	return nil
}

func (c *CppPrinter) PrintComponent(component IrComponent, iteratorTypeName string) error {
	// TODO: Use PrintType() for types and handle namespaces correctly.
	c.printf(`template<>
struct ecs::Component<%s> : public ecs::StaticComponent<%d>::Component<%s> {
};`, component.ElemType, component.Length, component.ElemType)

	c.printf(`template<>
struct ecs::Iterator<%s> : public ecs::StaticComponent<%d>::Iterator<%s> {
};`, component.ElemType, component.Length, component.ElemType)

	c.printf(`using %s = ecs::StaticComponent<%d>::IteratorImpl<%s>;`,
		iteratorTypeName, component.Length, component.ElemType)

	return nil
}

func (p *CppPrinter) PrintFunction(function IrFunction, isExport bool) {
	p.printInNamespace(function.ID, func(id string) {
		{
			// Print template type (if any).
			if typeVars := function.TypeVars; len(typeVars) > 0 {
				p.printf("template <typename %s", typeVars[0])
				for _, tvar := range typeVars[1:] {
					p.printf(", typename %s", tvar)
				}
				p.printf(">")
			}
		}

		{
			// Print ret type.
			retTypes := make([]IrType, len(function.Rets))
			for i := range function.Rets {
				retTypes[i] = function.Rets[i].Term.Type
			}

			p.withBindPosition(func() { p.printType(NewTupleType(retTypes)) })
		}

		// Print id.
		p.printf(" %s(", id)

		// Print args.
		switch args := function.Args; len(args) {
		case 0:
			break
		case 1:
			p.withBindPosition(func() { p.printType(args[0].Term.Type) })
			p.printf(" %s", args[0].Term.ID)
		default:
			p.withBindPosition(func() { p.printType(args[0].Term.Type) })
			p.printf(" %s", args[0].Term.ID)
			for _, arg := range args[1:] {
				p.printf(", ")
				p.withBindPosition(func() { p.printType(arg.Term.Type) })
				p.printf(" %s", arg.Term.ID)
			}
		}

		p.printf(") {\n")

		for _, ret := range function.Rets {
			p.withBindPosition(func() { p.printType(ret.Term.Type) })
			p.printf(" %s;\n", ret.Term.ID)
		}

		p.PrintTerm(function.Body)

		{
			retIDs := make([]string, len(function.Rets))
			for i, decl := range function.Rets {
				retIDs[i] = decl.Term.ID
			}
			p.printReturn(function.ID, retIDs)
		}
		p.printf("}\n")
	})
}

func (p *CppPrinter) PrintTerm(term IrTerm) {
	switch term.Case {
	case AppTermTerm:
		id, types, arg := term.AppArgs()
		p.printCall(id, types, arg)

	case AssignTerm:
		p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
		p.printf(" = ")
		p.PrintTerm(term.Assign.Arg)

	case BlockTerm:
		c := term.Block
		p.printf("{\n")
		for _, term := range c.Terms {
			p.PrintTerm(term)
			p.printf(";")
		}
		p.printf("}\n")

	case IfTerm:
		c := term.If

		p.printf("if (")
		if c.Negate {
			p.printf("!")
		}
		p.PrintTerm(c.Condition)
		p.printf(") ")
		p.PrintTerm(c.Then)
		if c.Else != nil {
			p.printf(" else ")
			p.PrintTerm(*c.Else)
		}

	case IndexGetTerm:
		if term.IndexGet.Obj.Type.Is(TupleType) {
			p.printf("std::get<")
			p.PrintTerm(term.IndexGet.Index)
			p.printf(">(")
			p.PrintTerm(term.IndexGet.Obj)
			p.printf(")")
		} else if len(term.IndexGet.Field) == 0 {
			p.PrintTerm(term.IndexGet.Obj)
			p.printf("[")
			p.PrintTerm(term.IndexGet.Index)
			p.printf("]")
		} else {
			p.PrintTerm(term.IndexGet.Obj)
			p.printf(".%s", term.IndexGet.Field)
		}

	case IndexSetTerm:
		if term.IndexSet.Obj.Type.Is(TupleType) {
			p.printf("std::get<")
			p.PrintTerm(term.IndexSet.Index)
			p.printf(">(")
			p.PrintTerm(term.IndexSet.Obj)
			p.printf(") = ")
			p.PrintTerm(term.IndexSet.Value)
		} else if len(term.IndexSet.Field) == 0 {
			p.PrintTerm(term.IndexSet.Obj)
			p.printf("[")
			p.PrintTerm(term.IndexSet.Index)
			p.printf("] = ")
			p.PrintTerm(term.IndexSet.Value)
		} else {
			p.PrintTerm(term.IndexSet.Obj)
			p.printf(".%s = ", term.IndexSet.Field)
			p.PrintTerm(term.IndexSet.Value)
		}

	case LetTerm:
		c := term.Let
		p.PrintDecl(c.Decl)

	case TokenTerm:
		p.printToken(*term.Token)

	case TupleTerm:
		if p.position == BindPosition {
			p.printf("std::tie(")
		}

		if len(term.Tuple) > 0 {
			p.PrintTerm(term.Tuple[0])
			for _, term := range term.Tuple[1:] {
				p.printf(", ")
				p.PrintTerm(term)
			}
		}

		if p.position == BindPosition {
			p.printf(")")
		}

	default:
		panic(fmt.Errorf("unhandled IrTerm %d", term.Case))
	}
}

func NewCppPrinter(output io.Writer) *CppPrinter {
	printer := &CppPrinter{
		output,
		TypePosition,
	}
	return printer
}
