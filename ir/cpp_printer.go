package ir

import (
	"fmt"
	"io"
	"strings"
)

type Position int

const (
	TypePosition = Position(iota)
	BindPosition
	ReturnPosition
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
	auto     bool
}

func (p *CppPrinter) withBindPosition(callback func()) {
	position := p.position
	p.position = BindPosition
	defer func() { p.position = position }()
	callback()
}

func (p *CppPrinter) withAutoType(value bool, callback func()) {
	auto := p.auto
	p.auto = value
	defer func() { p.auto = auto }()
	callback()
}

func (p *CppPrinter) withReturnPosition(callback func()) {
	position := p.position
	p.position = ReturnPosition
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

func (p *CppPrinter) printCast(arg IrTerm, types []IrType) {
	p.printf("static_cast<")
	p.withBindPosition(func() {
		p.printType(types[0])
		for _, typ := range types[1:] {
			p.printf(", ")
			p.printType(typ)
		}
	})
	p.printf(">")
	p.printf("(")
	p.PrintTerm(arg)
	p.printf(")")
}

func (p *CppPrinter) printCall(id IrTerm, types []IrType, arg IrTerm) {
	p.PrintTerm(id)
	if id.Is(VarTerm) && !IsOperator(id.Var.ID) && len(types) > 0 {
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
	if arg.Is(TupleTerm) {
		args := arg.Tuple.Elems
		if len(args) > 0 {
			p.PrintTerm(args[0])
			for _, t := range args[1:] {
				p.printf(", ")
				p.PrintTerm(t)
			}
		}
	} else {
		p.PrintTerm(arg)
	}
	p.printf(")")
}

func (p *CppPrinter) printAliasDecl(id string, typ IrType) {
	switch typ.Case {
	case LambdaType:
		tvars := typ.LambdaVars()
		p.printf("template <typename %s", tvars[0])
		for _, tvar := range tvars[1:] {
			p.printf(", typename %s", tvar)
		}
		p.printf("> ")
		p.printAliasDecl(id, typ.LambdaBody())

	case StructType:
		p.printf("struct %s {\n", id)
		for _, field := range typ.Fields() {
			p.printType(field.Type)
			p.printf(" %s;\n", field.ID)
		}
		p.printf("};\n")

	case VariantType:
		p.printf("using %s = ", id)
		p.printType(typ)
		p.printf(";\n")

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) printType(typ IrType) {
	if p.auto {
		p.printf("auto")
		return
	}

	switch {
	case typ.Is(AppType):
		p.printType(typ.App.Fun)
		args := typ.AppArgs()
		p.printf("<")
		p.printType(args[0])
		for _, arg := range args[1:] {
			p.printf(", ")
			p.printType(arg)
		}
		p.printf(">")

	case typ.Is(ArrayType):
		p.printf("std::array<")
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

	case typ.Is(FunType):
		c := typ.Fun

		p.printf("std::function<")
		p.printType(c.Ret)
		p.printf("(")
		p.printType(c.Arg)
		p.printf(")>")

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
		if len(tuple.Elems) > 0 {
			p.printType(tuple.Elems[0])
			for _, elem := range tuple.Elems[1:] {
				p.printf(", ")
				p.printType(elem)
			}
		}

	case typ.Is(TupleType) && p.position == BindPosition:
		tuple := typ.Tuple
		// Print rets.
		switch len(tuple.Elems) {
		case 0:
			p.printf("void")
		case 1:
			p.printType(tuple.Elems[0])
		default:
			p.printf("std::tuple<")
			p.printType(tuple.Elems[0])
			for _, elem := range tuple.Elems[1:] {
				p.printf(", ")
				p.printType(elem)
			}
			p.printf(">")
		}

	case typ.Is(VariantType):
		p.printf("std::variant<")
		if tags := typ.Tags(); len(tags) > 0 {
			p.printType(tags[0].Type)
			p.printf("/* %s */", toID(tags[0].ID))
			for _, tag := range tags[1:] {
				p.printf(", ")
				p.printType(tag.Type)
				p.printf("/* %s */", toID(tag.ID))
			}
		}
		p.printf(">")

	case typ.Is(VarType):
		p.printf("%s", typ.Var)

	default:
		panic(fmt.Errorf("printType: unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (p *CppPrinter) PrintDecl(decl IrDecl, export bool) {
	if export {
		p.printf("export ")
	}

	if decl.Is(NameDecl) {
		p.printType(NewNameType(decl.Name.ID))
		return
	}

	if decl.Is(AliasDecl) {
		p.printAliasDecl(decl.Alias.ID, decl.Alias.Type)
		return
	}

	switch typ := decl.Term.Type; typ.Case {
	case AppType, ArrayType, NameType, TupleType, VarType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.withBindPosition(func() {
				p.printType(typ)
				p.printf(" %s", id)
			})
		})

	case ForallType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			tvars := typ.ForallVars()
			p.printf("template <typename %s", tvars[0])
			for _, tvar := range tvars[1:] {
				p.printf(", typename %s", tvar)
			}
			p.printf("> ")
			p.PrintDecl(NewTermDecl(id, typ.ForallBody()), false /* export */)
		})

	case FunType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.withBindPosition(func() { p.printType(typ.Fun.Ret) })
			p.printf(" %s(", id)
			p.printType(typ.Fun.Arg)
			p.printf(");")
		})

	case StructType:
		// TODO: Handle namespacing.
		p.printf("struct %s", decl.Term.ID)

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (p *CppPrinter) PrintModuleTop(moduleName string) {
	if !strings.Contains(moduleName, "_") && moduleName != "program" {
		p.printf("export module %s;\n", moduleName)
		return
	}

	p.printf("export module %s;\n", moduleName)
	p.printf("\n")
	p.printf("import <array>;\n")
	p.printf("import <cstdlib>;\n")
	p.printf("import <functional>;\n")
	p.printf("import <iostream>;\n")
	p.printf("import <tuple>;\n")
	p.printf("import <variant>;\n")
	p.printf("import <vector>;\n")
	p.printf("\n")
}

func (p *CppPrinter) Import(module string) {
	p.printf("import %s;\n", module)
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
		p.printf(`
	// Needed because of import<vector> results in Bad file data:
	// https://stackoverflow.com/questions/70456868/vector-in-c-module-causes-useless-bad-file-data-gcc-output
	namespace std _GLIBCXX_VISIBILITY(default){}
	`)
	default:
		return fmt.Errorf("unknown section %q", id)
	}

	for _, decl := range decls {
		if isComment {
			p.printf(" * ")
		}
		p.PrintDecl(decl, false /* export */)
		p.printf("\n")
	}

	if isComment {
		p.printf("*/")
	}
	p.printf("\n")

	return nil
}

func (p *CppPrinter) PrintImpls(module string, ids []string) error {
	for _, id := range ids {
		p.printf("export import :%s;\n", id)
	}
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
		if isExport {
			p.printf("export ")
		}

		{
			// Print template type (if any).
			if typeVars := function.TypeVars; len(typeVars) > 0 {
				p.printf("template <typename %s", typeVars[0].Var)
				for _, tvar := range typeVars[1:] {
					p.printf(", typename %s", tvar.Var)
				}
				p.printf(">")
			}
		}

		{
			// Print ret type.
			p.withBindPosition(func() { p.printType(function.RetType) })
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

		p.printf(")\n")
		p.PrintTerm(function.Body)
		p.printf("\n")
	})
}

func (p *CppPrinter) PrintTerm(term IrTerm) {
	if p.position == ReturnPosition || term.LastTerm {
		returning := term.LastTerm && !term.Is(IndexSetTerm)
		if returning {
			p.printf("return")
		}

		if term.Is(TupleTerm) && len(term.Tuple.Elems) == 0 {
			return
		}

		if returning {
			p.printf(" ")
		}
	}

	switch {
	case term.Is(AppTypeTerm):
		term, types := term.AppTypes()
		p.printCast(term, types)

	case term.Is(AppTermTerm):
		id, types, arg := term.AppArgs()
		p.printCall(id, types, arg)

	// TODO: This doesn't seem to be needed. Delete.
	case term.Is(AssignTerm) && term.Assign.Arg.Is(TupleTerm):
		p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
		p.printf(" = ")
		p.printf("std::make_tuple(")
		p.PrintTerm(term.Assign.Arg)
		p.printf(")")

	case term.Is(AssignTerm):
		p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
		p.printf(" = ")
		p.PrintTerm(term.Assign.Arg)

	case term.Is(BlockTerm):
		c := term.Block
		p.printf("{\n")
		for _, term := range c.Terms {
			p.PrintTerm(term)
			p.printf(";")
		}
		p.printf("}\n")

	case term.Is(ConstTerm):
		p.printf("%s", term.Const.Value)

	case term.Is(IfTerm):
		c := term.If

		p.printf("if (")
		p.PrintTerm(c.Condition)
		p.printf(") ")
		p.PrintTerm(c.Then)
		if c.Else != nil {
			p.printf(" else ")
			p.PrintTerm(*c.Else)
		}

	case term.Is(InjectionTerm):
		c := term.Injection

		p.printType(c.VariantType)
		p.printf("{")
		p.printf("std::in_place_index<%d>, ", *c.TagIndex)
		p.PrintTerm(c.Value)
		p.printf("}")

	case term.Is(StructTerm):
		c := term.Struct

		p.printf("{")
		if len(c.Values) > 0 {
			p.printf(".%s = ", c.Values[0].Label)
			p.PrintTerm(c.Values[0].Value)
			for _, f := range c.Values[1:] {
				p.printf(", .%s = ", f.Label)
				p.PrintTerm(f.Value)
			}
		}
		p.printf("}")

	case term.Is(IndexGetTerm):
		if term.IndexGet.Obj.Type.Is(TupleType) || term.IndexGet.Obj.Type.Is(VariantType) {
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

	case term.Is(IndexSetTerm):
		if term.IndexSet.Obj.Type.Is(TupleType) {
			p.printf("std::get<")
			p.PrintTerm(term.IndexSet.Index)
			p.printf(">(")
			p.PrintTerm(term.IndexSet.Obj)
			p.printf(") = ")
			p.PrintTerm(term.IndexSet.Value)
		} else if term.IndexSet.TagIndex != nil {
			p.PrintTerm(term.IndexSet.Obj)
			p.printf(" = ")
			p.printType(*term.IndexSet.Obj.Type)
			p.printf("{std::in_place_index<%d>, ", *term.IndexSet.TagIndex)
			p.PrintTerm(term.IndexSet.Value)
			p.printf("}")
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

	case term.Is(LambdaTerm) || term.Is(TypeAbsTerm):
		tvars, args, argTypes, body := term.ToFunction()
		p.printf("[]")

		// Print type abstraction types.
		if len(tvars) > 0 {
			p.printf("<")
			interleave(tvars, func() { p.printf(", ") }, func(_ int, tvar string) {
				p.printf("typename %s", tvar)
			})
			p.printf(">")
		}

		// Print abstraction arguments and types.
		p.printf("(")
		interleave(args, func() { p.printf(", ") }, func(i int, arg string) {
			p.printType(argTypes[i])
			p.printf(" %s", toID(arg))
		})
		p.printf(")")

		// Print abstraction body.
		p.printf("{ return ")
		p.PrintTerm(body)
		p.printf("; }")

	case term.Is(LetTerm):
		c := term.Let

		// There's no type (e.g., std::function) in C++20 for polymorphic
		// lambdas, so 'auto' must be used instead.
		//
		// For example:
		//   auto id = []<typename T>(T x) { return x; };
		auto := c.Value.Is(TypeAbsTerm)

		p.withAutoType(auto, func() {
			p.withBindPosition(func() {
				p.printType(c.VarType)
				p.printf(" %s", c.Var)
			})
		})
		p.printf(" = ")
		p.PrintTerm(c.Value)

	case term.Is(ProjectionTerm):
		c := term.Projection

		if c.Term.Type.Is(ArrayType) {
			p.PrintTerm(c.Term)
			p.printf("[")
			p.PrintTerm(c.Label)
			p.printf("]")
		} else if c.Term.Type.Is(StructType) {
			p.PrintTerm(c.Term)
			p.printf(".%s", *c.LabelName)
		} else {
			p.printf("std::get<%d>(", *c.Index)
			p.PrintTerm(c.Term)
			p.printf(")")
		}

	case term.Is(ReturnTerm):
		c := term.Return
		p.printf("return ")
		p.withReturnPosition(func() { p.PrintTerm(c.Expr) })
		p.printf(";")

	case term.Is(TupleTerm):
		if p.position == BindPosition {
			p.printf("std::tie(")
		} else {
			p.printf("std::make_tuple(")
		}

		if len(term.Tuple.Elems) > 0 {
			p.PrintTerm(term.Tuple.Elems[0])
			for _, term := range term.Tuple.Elems[1:] {
				p.printf(", ")
				p.PrintTerm(term)
			}
		}

		p.printf(")")

	case term.Is(VarTerm):
		p.printf("%s", toID(term.Var.ID))

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func NewCppPrinter(output io.Writer) *CppPrinter {
	printer := &CppPrinter{
		output,
		TypePosition,
		false, /* auto */
	}
	return printer
}
