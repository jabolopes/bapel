package comp

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
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

func toPartitionName(filename ir.Filename) string {
	return parse.TrimExtension(path.Base(filename.Value))
}

func toModulePartitionName(unit ir.IrUnit) string {
	switch unit.Case {
	case ir.BaseUnit:
		return unit.ModuleID.Name
	case ir.ImplUnit:
		return fmt.Sprintf("%s:%s", unit.ModuleID, toPartitionName(unit.Filename))
	default:
		panic(fmt.Errorf("unhandled %T %d", unit.Case, unit.Case))
	}
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
	output   io.Writer
	position Position
	autoType bool
	idgen    int
	err      error
}

func (p *CppPrinter) genID() string {
	id := fmt.Sprintf("__v_%d", p.idgen)
	p.idgen++
	return id
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
	if arg.Is(ir.ConstTerm) {
		p.printf("static_cast<")
		p.withBindPosition(func() {
			ir.Interleave(types, func() { p.printf(", ") }, func(_ int, typ ir.IrType) {
				p.printType(typ)
			})
		})
		p.printf(">(")
		p.PrintTerm(arg)
		p.printf(")")
		return
	}

	p.PrintTerm(arg)
	p.printf("<")
	p.withBindPosition(func() {
		ir.Interleave(types, func() { p.printf(", ") }, func(_ int, typ ir.IrType) {
			p.printType(typ)
		})
	})
	p.printf(">")
}

func (p *CppPrinter) printCall(id ir.IrTerm, types []ir.IrType, arg ir.IrTerm) {
	if id.Is(ir.VarTerm) && ir.IsOperator(id.Var.ID) {
		if arg.Is(ir.TupleTerm) {
			p.printf("(")
			p.PrintTerm(arg.Tuple.Elems[0])
			p.printf(")")

			p.PrintTerm(id)

			p.printf("(")
			p.PrintTerm(ir.NewTupleTerm(arg.Tuple.Elems[1:]))
			p.printf(")")
		} else {
			p.PrintTerm(id)
			p.printf(" ")
			p.PrintTerm(arg)
		}

		return
	}

	p.PrintTerm(id)

	if id.Is(ir.VarTerm) && !ir.IsOperator(id.Var.ID) && len(types) > 0 {
		p.printf("<")
		p.withBindPosition(func() {
			ir.Interleave(types, func() { p.printf(", ") }, func(_ int, typ ir.IrType) {
				p.printType(typ)
			})
		})
		p.printf(">")
	}

	p.printf("(")
	if arg.Is(ir.TupleTerm) {
		ir.Interleave(arg.Tuple.Elems, func() { p.printf(", ") }, func(_ int, t ir.IrTerm) {
			p.PrintTerm(t)
		})
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
		p.withBindPosition(func() { p.printType(c.Ret) })
		p.printf("(")
		p.withBindPosition(func() { p.printType(c.Arg) })
		p.printf(")>")

	case typ.Is(ir.NameType):
		switch typ.Name {
		case "bool":
			p.printf("bool")
		case "i8":
			p.printf("int8_t")
		case "i16":
			p.printf("int16_t")
		case "i32":
			p.printf("int32_t")
		case "i64":
			p.printf("int64_t")
		case "f32":
			p.printf("float")
		case "f64":
			p.printf("double")
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
			p.printf("std::monostate")
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
			p.withBindPosition(func() { p.printType(tags[0].Type) })
			p.printf("/* %s */", toID(tags[0].ID))
			for _, tag := range tags[1:] {
				p.printf(", ")
				p.withBindPosition(func() { p.printType(tag.Type) })
				p.printf("/* %s */", toID(tag.ID))
			}
		}
		p.printf(">")

	case typ.Is(ir.VarType):
		p.printf("%s", typ.Var)

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (p *CppPrinter) printDecl(decl ir.IrDecl) {
	if decl.Export {
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
			p.printDecl(ir.NewTermDecl(id, typ.ForallBody(), false /* export */))
		})

	case ir.FunType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.withBindPosition(func() { p.printType(typ.Fun.Ret) })
			p.printf(" %s(", id)
			p.printType(typ.Fun.Arg)
			p.printf(");")
		})

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (p *CppPrinter) printTypeDef(decl ir.IrDecl) {
	if decl.Export {
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

func (p *CppPrinter) printModuleTop(moduleName string) {
	p.printf("module;\n")
	p.printf("\n")
	p.printf("#include <array>\n")
	p.printf("#include <cstdlib>\n")
	p.printf("#include <functional>\n")
	p.printf("#include <optional>\n")
	p.printf("#include <string>\n")
	p.printf("#include <tuple>\n")
	p.printf("#include <variant>\n")
	p.printf("#include <vector>\n")
	p.printf("\n")
	p.printf("export module %s;\n", moduleName)
	p.printf("\n")
}

func (p *CppPrinter) printImports(imports []ir.IrImport) {
	p.printf("\n")
	for _, imp := range imports {
		p.printf("import %s;\n", imp.ModuleID)
	}
	p.printf("\n")
}

func (p *CppPrinter) printImpls(impls []ir.IrImpl) {
	p.printf("\n")
	for _, impl := range impls {
		p.printf("export import :%s;\n", toPartitionName(impl.RelativeFilename))
	}
	p.printf("\n")
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

func (p *CppPrinter) printProjectionTerm(term ir.IrTerm) error {
	if !term.Is(ir.ProjectionTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ProjectionTerm, ir.ProjectionTerm))
	}

	c := term.Projection

	if c.Term.Type.Is(ir.StructType) {
		_, field, err := c.Term.Type.FieldByLabel(c.Label)
		if err != nil {
			return err
		}

		p.PrintTerm(c.Term)
		p.printf(".%s", field.ID)
		return nil
	}

	var index int
	var err error

	if c.Term.Type.Is(ir.TupleType) {
		index, _, err = c.Term.Type.ElemByLabel(c.Label)
	} else if c.Term.Type.Is(ir.VariantType) {
		index, _, err = c.Term.Type.TagByLabel(c.Label)
	}

	if err != nil {
		return err
	}

	p.printf("std::get<%d>(", index)
	p.PrintTerm(c.Term)
	p.printf(")")

	return nil
}

func (p *CppPrinter) printSetTerm(term ir.IrTerm) error {
	if !term.Is(ir.SetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.SetTerm, ir.SetTerm))
	}

	c := term.Set

	switch {
	case term.Type.Is(ir.StructType):
		structID := p.genID()

		p.printf("([%s = ", structID)
		p.PrintTerm(c.Term)
		p.printf("]() mutable {\n")
		for _, lv := range c.Values {
			_, field, err := term.Type.FieldByLabel(lv.Label)
			if err != nil {
				return err
			}

			p.printf("%s.%s = ", structID, field.ID)
			p.PrintTerm(lv.Value)
			p.printf(";\n")
		}
		p.printf("return %s;\n", structID)
		p.printf("})()")

	case term.Type.Is(ir.TupleType):
		tupleID := p.genID()

		p.printf("([%s = ", tupleID)
		p.PrintTerm(c.Term)
		p.printf("]() mutable {\n")
		for _, lv := range c.Values {
			p.printf("std::get<%s>(%s) = ", lv.Label, tupleID)
			p.PrintTerm(lv.Value)
			p.printf(";\n")
		}
		p.printf("return %s;\n", tupleID)
		p.printf("})()")

	default:
		panic(fmt.Errorf("unhandled type %s", term.Type))
	}

	return nil
}

func (p *CppPrinter) PrintTerm(term ir.IrTerm) {
	if p.position == ReturnPosition || term.LastTerm {
		returning := term.LastTerm
		if returning {
			p.printf("return")
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

	case term.Is(ir.SetTerm):
		if err := p.printSetTerm(term); err != nil {
			p.err = errors.Join(p.err, err)
		}

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

	case term.Is(ir.LambdaTerm) || term.Is(ir.TypeAbsTerm):
		tvars, args, argTypes, body := term.ToFunction()
		p.printf("[]")

		// Print type abstraction types.
		if len(tvars) > 0 {
			p.printf("<")
			ir.Interleave(tvars, func() { p.printf(", ") }, func(_ int, tvar string) {
				p.printf("typename %s", tvar)
			})
			p.printf(">")
		}

		// Print abstraction arguments and types.
		p.printf("(")
		ir.Interleave(args, func() { p.printf(", ") }, func(i int, arg string) {
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
		if err := p.printProjectionTerm(term); err != nil {
			p.err = errors.Join(p.err, err)
		}

	case term.Is(ir.MatchTerm):
		c := term.Match

		variantID := p.genID()

		p.printf("([&] {")
		p.printf("auto %s = ", variantID)
		p.PrintTerm(c.Term)
		p.printf(";\n")
		p.printf("switch (%s.index()) {\n", variantID)
		for _, arm := range c.Arms {
			p.printf("case %d: {", *arm.Index)
			p.printf("auto &%s = std::get<%d>(%s);\n", arm.Arg, *arm.Index, variantID)
			p.printf("return ")
			p.PrintTerm(arm.Body)
			p.printf(";\n")
			p.printf("}\n")
		}
		p.printf("} })")

	case term.Is(ir.ReturnTerm):
		c := term.Return

		p.printf("return ")
		p.withReturnPosition(func() { p.PrintTerm(c.Expr) })
		p.printf(";")

	case term.Is(ir.TupleTerm):
		if p.position == BindPosition {
			p.printf("std::tie(")
		} else if len(term.Tuple.Elems) == 0 {
			p.printf("std::monostate(")
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

func (p *CppPrinter) printDecls(decls []ir.IrDecl) {
	for _, decl := range decls {
		switch {
		case decl.Is(ir.AliasDecl):
			p.printTypeDef(decl)
		default:
			p.printDecl(decl)
		}
	}
}

func (p *CppPrinter) printUnit(unit ir.IrUnit) error {
	p.printModuleTop(toModulePartitionName(unit))
	p.printImpls(unit.Impls)
	p.printImports(unit.Imports)
	p.printDecls(unit.Decls)
	for _, function := range unit.Functions {
		p.printFunction(function)
	}

	return nil
}

func newCppPrinter(output io.Writer) *CppPrinter {
	return &CppPrinter{
		output,
		TypePosition,
		false, /* autoType */
		0,     /* idgen */
		nil,   /* err */
	}
}

func printUnitToCpp(unit ir.IrUnit, output io.Writer) error {
	printer := newCppPrinter(output)
	printErr := printer.printUnit(unit)
	return errors.Join(printErr, printer.err)
}
