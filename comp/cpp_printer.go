package comp

import (
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
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

func countTypeVars(kind ir.IrKind) int {
	switch kind.Case {
	case ir.ArrowKind:
		return 1 + countTypeVars(kind.Arrow.Arg)
	default:
		return 0
	}
}

type CppPrinter struct {
	output     io.Writer
	position   Position
	autoType   bool
	moduleName string
}

func (p *CppPrinter) withBindPosition(callback func()) {
	position := p.position
	p.position = BindPosition
	defer func() { p.position = position }()
	callback()
}

func (p *CppPrinter) withAutoType(value bool, callback func()) {
	autoType := p.autoType
	p.autoType = value
	defer func() { p.autoType = autoType }()
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

func (p *CppPrinter) printCast(arg ir.IrTerm, types []ir.IrType) {
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

func (p *CppPrinter) printCall(id ir.IrTerm, types []ir.IrType, arg ir.IrTerm) {
	p.PrintTerm(id)
	if id.Is(ir.VarTerm) && !ir.IsOperator(id.Var.ID) && len(types) > 0 {
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
	if arg.Is(ir.TupleTerm) {
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

func (p *CppPrinter) printAliasDecl(id string, typ ir.IrType) {
	switch typ.Case {
	case ir.LambdaType:
		tvars := typ.LambdaVars()
		p.printf("template <typename %s", tvars[0])
		for _, tvar := range tvars[1:] {
			p.printf(", typename %s", tvar)
		}
		p.printf("> ")
		p.printAliasDecl(id, typ.LambdaBody())

	case ir.StructType:
		p.printf("struct %s {\n", id)
		for _, field := range typ.Fields() {
			p.printType(field.Type)
			p.printf(" %s;\n", field.ID)
		}
		p.printf("}\n")

	case ir.VariantType:
		p.printf("struct %s : ", id)
		p.printType(typ)
		p.printf("{}\n")

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) printType(typ ir.IrType) {
	if p.autoType {
		p.printf("auto")
		return
	}

	switch {
	case typ.Is(ir.AppType):
		p.printType(typ.App.Fun)
		args := typ.AppArgs()
		p.printf("<")
		p.printType(args[0])
		for _, arg := range args[1:] {
			p.printf(", ")
			p.printType(arg)
		}
		p.printf(">")

	case typ.Is(ir.ArrayType):
		p.printf("std::array<")
		p.printType(typ.Array.ElemType)
		p.printf(", %d>", typ.Array.Size)

	case typ.Is(ir.ForallType):
		tvars := typ.ForallVars()
		p.printf("template <typename %s", tvars[0])
		for _, tvar := range tvars[1:] {
			p.printf(", typename %s", tvar)
		}
		p.printf("> ")
		p.printType(typ.ForallBody())

	case typ.Is(ir.FunType):
		c := typ.Fun

		p.printf("std::function<")
		p.printType(c.Ret)
		p.printf("(")
		p.printType(c.Arg)
		p.printf(")>")

	case typ.Is(ir.NameType):
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

	case typ.Is(ir.TupleType) && p.position == TypePosition:
		tuple := typ.Tuple
		if len(tuple.Elems) > 0 {
			p.printType(tuple.Elems[0])
			for _, elem := range tuple.Elems[1:] {
				p.printf(", ")
				p.printType(elem)
			}
		}

	case typ.Is(ir.TupleType) && p.position == BindPosition:
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

	case typ.Is(ir.VariantType):
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

	case typ.Is(ir.VarType):
		p.printf("%s", typ.Var)

	default:
		panic(fmt.Errorf("printType: unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (p *CppPrinter) printDecl(decl ir.IrDecl, export bool) {
	if export {
		p.printf("export ")
	}

	if decl.Is(ir.NameDecl) {
		p.printInNamespace(decl.Name.ID, func(id string) {
			if args := countTypeVars(decl.Name.Kind); args > 0 {
				p.printf("template <")
				p.printf("typename t%d", 0)
				for i := 1; i < args; i++ {
					p.printf(", typename t%d", i)
				}
				p.printf("> ")
			}

			p.printf("struct %s;\n", id)
		})
		return
	}

	if decl.Is(ir.AliasDecl) {
		p.printInNamespace(decl.Alias.ID, func(id string) {
			switch typ := decl.Alias.Type; typ.Case {
			case ir.LambdaType:
				tvars := typ.LambdaVars()
				p.printf("template <typename %s", tvars[0])
				for _, tvar := range tvars[1:] {
					p.printf(", typename %s", tvar)
				}
				p.printf("> struct %s", id)

			default:
				p.printf("struct %s", id)
			}
			p.printf(";\n")
		})
		return
	}

	switch typ := decl.Term.Type; typ.Case {
	case ir.AppType, ir.ArrayType, ir.NameType, ir.TupleType, ir.VarType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.withBindPosition(func() {
				p.printType(typ)
				p.printf(" %s", id)
			})
		})

	case ir.ForallType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			tvars := typ.ForallVars()
			p.printf("template <typename %s", tvars[0])
			for _, tvar := range tvars[1:] {
				p.printf(", typename %s", tvar)
			}
			p.printf("> ")
			p.printDecl(ir.NewTermDecl(id, typ.ForallBody()), false /* export */)
		})

	case ir.FunType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.withBindPosition(func() { p.printType(typ.Fun.Ret) })
			p.printf(" %s(", id)
			p.printType(typ.Fun.Arg)
			p.printf(");")
		})

	case ir.StructType:
		// TODO: Handle namespacing.
		p.printf("struct %s", decl.Term.ID)

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (p *CppPrinter) printTypeDef(decl ir.IrDecl, export bool) {
	if export {
		p.printf("export ")
	}

	switch {
	case decl.Is(ir.NameDecl):
		p.printType(ir.NewNameType(decl.Name.ID))
		p.printf(";\n")

	case decl.Is(ir.AliasDecl):
		p.printInNamespace(decl.Alias.ID, func(id string) {
			p.printAliasDecl(id, decl.Alias.Type)
			p.printf(";")
		})

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", decl.Case, decl.Case, decl))
	}
}

func (p *CppPrinter) PrintModuleTop(moduleName string) {
	p.printf("module;\n")
	p.printf("\n")
	p.printf("#include <array>\n")
	p.printf("#include <cstdlib>\n")
	p.printf("#include <functional>\n")
	p.printf("#include <tuple>\n")
	p.printf("#include <variant>\n")
	p.printf("#include <vector>\n")
	p.printf("\n")
	p.printf("export module %s;\n", moduleName)
	p.printf("\n")
}

func (p *CppPrinter) printImports(imports ast.Imports) {
	p.printf("\n")
	for _, moduleName := range imports.IDs {
		p.printf("import %s;\n", moduleName.Value)
	}
}

func (p *CppPrinter) printImpls(impls ast.Impls) {
	p.printf("\n")
	for _, id := range impls.IDs {
		p.printf("export import :%s;\n", bplparser2.TrimExtension(id.Value))
	}
}

func (c *CppPrinter) printComponent(component ir.IrComponent) {
	iteratorTypeName := fmt.Sprintf("%s_iterator", component.ElemType)

	// TODO: Use PrintType() for types and handle namespaces correctly.
	c.printf(`template<>
struct ecs::Component<%s> : public ecs::StaticComponent<%d>::Component<%s> {
};`, component.ElemType, component.Length, component.ElemType)

	c.printf(`template<>
struct ecs::Iterator<%s> : public ecs::StaticComponent<%d>::Iterator<%s> {
};`, component.ElemType, component.Length, component.ElemType)

	c.printf(`using %s = ecs::StaticComponent<%d>::IteratorImpl<%s>;`,
		iteratorTypeName, component.Length, component.ElemType)
}

func (p *CppPrinter) printFunction(function ir.IrFunction) {
	p.printInNamespace(function.ID, func(id string) {
		if function.Export {
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

func (p *CppPrinter) PrintTerm(term ir.IrTerm) {
	if p.position == ReturnPosition || term.LastTerm {
		returning := term.LastTerm && !term.Is(ir.IndexSetTerm)
		if returning {
			p.printf("return")
		}

		if term.Is(ir.TupleTerm) && len(term.Tuple.Elems) == 0 {
			return
		}

		if returning {
			p.printf(" ")
		}
	}

	switch {
	case term.Is(ir.AppTypeTerm):
		term, types := term.AppTypes()
		p.printCast(term, types)

	case term.Is(ir.AppTermTerm):
		id, types, arg := term.AppArgs()
		p.printCall(id, types, arg)

	// TODO: This doesn't seem to be needed. Delete.
	case term.Is(ir.AssignTerm) && term.Assign.Arg.Is(ir.TupleTerm):
		p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
		p.printf(" = ")
		p.printf("std::make_tuple(")
		p.PrintTerm(term.Assign.Arg)
		p.printf(")")

	case term.Is(ir.AssignTerm):
		p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
		p.printf(" = ")
		p.PrintTerm(term.Assign.Arg)

	case term.Is(ir.BlockTerm):
		c := term.Block
		p.printf("{\n")
		for _, term := range c.Terms {
			p.PrintTerm(term)
			p.printf(";")
		}
		p.printf("}\n")

	case term.Is(ir.ConstTerm) && term.Const.Is(ir.IntLiteral):
		p.printf("%d", *term.Const.Int)

	case term.Is(ir.ConstTerm) && term.Const.Is(ir.StrLiteral):
		p.printf(`"%s"`, *term.Const.Str)

	case term.Is(ir.IfTerm):
		c := term.If

		p.printf("if (")
		p.PrintTerm(c.Condition)
		p.printf(") ")
		p.PrintTerm(c.Then)
		if c.Else != nil {
			p.printf(" else ")
			p.PrintTerm(*c.Else)
		}

	case term.Is(ir.InjectionTerm):
		c := term.Injection

		p.printType(c.VariantType)
		p.printf("{")
		p.printf("std::in_place_index<%d>, ", *c.TagIndex)
		p.PrintTerm(c.Value)
		p.printf("}")

	case term.Is(ir.StructTerm):
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

	case term.Is(ir.IndexSetTerm):
		if term.IndexSet.Obj.Type.Is(ir.TupleType) {
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

	case term.Is(ir.LambdaTerm) || term.Is(ir.TypeAbsTerm):
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

	case term.Is(ir.LetTerm):
		c := term.Let

		// There's no type (e.g., std::function) in C++20 for polymorphic
		// lambdas, so 'auto' must be used instead.
		//
		// For example:
		//   auto id = []<typename T>(T x) { return x; };
		auto := c.Value.Is(ir.TypeAbsTerm)

		p.withAutoType(auto, func() {
			p.withBindPosition(func() {
				p.printType(c.VarType)
				p.printf(" %s", c.Var)
			})
		})
		p.printf(" = ")
		p.PrintTerm(c.Value)

	case term.Is(ir.ProjectionTerm):
		c := term.Projection

		if c.Term.Type.Is(ir.ArrayType) {
			p.PrintTerm(c.Term)
			p.printf("[")
			p.PrintTerm(c.Label)
			p.printf("]")
		} else if c.Term.Type.Is(ir.StructType) {
			p.PrintTerm(c.Term)
			p.printf(".%s", *c.LabelName)
		} else {
			p.printf("std::get<%d>(", *c.Index)
			p.PrintTerm(c.Term)
			p.printf(")")
		}

	case term.Is(ir.ReturnTerm):
		c := term.Return
		p.printf("return ")
		p.withReturnPosition(func() { p.PrintTerm(c.Expr) })
		p.printf(";")

	case term.Is(ir.TupleTerm):
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

	case term.Is(ir.VarTerm):
		p.printf("%s", toID(term.Var.ID))

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (p *CppPrinter) printSource(source ast.Source) {
	switch source.Case {
	case ast.ComponentSource:
		p.printComponent(*source.Component)
	case ast.FunctionSource:
		p.printFunction(*source.Function)
	case ast.TypeDefSource:
		p.printTypeDef(source.TypeDef.Decl, source.TypeDef.Export)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (p *CppPrinter) doDecls(sources []ast.Source) {
	p.printf(`
// Needed because of import<vector> results in Bad file data:
// https://stackoverflow.com/questions/70456868/vector-in-c-module-causes-useless-bad-file-data-gcc-output
namespace std _GLIBCXX_VISIBILITY(default){}

`)

	for _, source := range sources {
		switch {
		case source.Is(ast.FunctionSource):
			p.printDecl(source.Function.Decl(), source.Function.Export)
		case source.Is(ast.TypeDefSource):
			p.printDecl(source.TypeDef.Decl, source.TypeDef.Export)
		}
	}

	p.printf("\n")
}

func (p *CppPrinter) PrintModule(module ast.Module) error {
	p.PrintModuleTop(p.moduleName)
	p.printImpls(module.Impls)
	p.printImports(module.Imports)
	p.doDecls(module.Body)
	for _, source := range module.Body {
		p.printSource(source)
	}
	return nil
}

func NewCppPrinter(output io.Writer, moduleName string) *CppPrinter {
	printer := &CppPrinter{
		output,
		TypePosition,
		false, /* autoType */
		moduleName,
	}
	return printer
}
