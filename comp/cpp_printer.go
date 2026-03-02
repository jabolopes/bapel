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
)

func toID(id string) string {
	if strings.Contains(id, ir.NamespaceSeparator) {
		return "::" + strings.Replace(id, ir.NamespaceSeparator, "::", -1)
	}
	return id
}

func toModuleName(moduleID ir.ModuleID) string {
	return strings.Replace(moduleID.Name, ir.ModuleIDSeparator, ".", -1)
}

func toPartitionName(filename ir.Filename) string {
	return parse.TrimExtension(path.Base(filename.Value))
}

func toModulePartitionName(unit ir.IrUnit) string {
	switch unit.Case {
	case ir.BaseUnit:
		return toModuleName(unit.ModuleID)
	case ir.ImplUnit:
		return fmt.Sprintf("%s:%s", toModuleName(unit.ModuleID), toPartitionName(unit.Filename))
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

func isCppStatement(term ir.IrTerm) bool {
	switch term.Case {
	case ir.AssignTerm, ir.BlockTerm, ir.IfTerm, ir.LetTerm, ir.MatchTerm, ir.ReturnTerm:
		return true
	}

	if term.Is(ir.AppTermTerm) {
		id, _, _ := term.AppArgs()
		if id.Is(ir.VarTerm) && (id.Var.ID == "ifthen" || id.Var.ID == "ifelse") {
			return true
		}
	}

	return false
}

type CppPrinter struct {
	output         io.Writer
	idgen          int
	position       Position
	autoType       bool
	lastTerm       bool
	varDestination string
	// Stores anonymous types (structs) keyed by the hash of the type.
	//
	// This is necessary because C++ does not allow anonymous structs,
	// e.g., as function arguments and return types, among others. This
	// gives a name based on the hash of the type and replaces anonymous
	// structs types with a typename with that hashed name.
	anonymousTypes map[string]anonymousType
	err            error
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

func (p *CppPrinter) withLastTerm(value bool, callback func()) {
	lastTerm := p.lastTerm
	p.lastTerm = value
	defer func() { p.lastTerm = lastTerm }()
	callback()
}

func (p *CppPrinter) withVarDestination(value string, callback func()) {
	orig := p.varDestination
	p.varDestination = value
	defer func() { p.varDestination = orig }()
	p.withLastTerm(true, callback)
}

func (p *CppPrinter) handleLastTerm(callback func()) {
	if !p.lastTerm {
		callback()
		return
	}

	if len(p.varDestination) == 0 {
		p.printf("return ")
	} else {
		p.printf("%s = ", p.varDestination)
	}
	p.withLastTerm(false, callback)
}

func (p *CppPrinter) printInNamespace(id string, callback func(string)) {
	if !strings.Contains(id, ir.NamespaceSeparator) {
		callback(id)
		return
	}

	p.printf("namespace ")

	tokens := strings.Split(id, ir.NamespaceSeparator)
	tokens, id = tokens[:len(tokens)-1], tokens[len(tokens)-1]

	ir.Interleave(tokens, func() { p.printf("::") }, func(_ int, token string) {
		p.printf("%s", token)
	})

	p.printf(" { ")
	callback(id)
	p.printf(" }")
}

func (p *CppPrinter) printf(format string, args ...any) {
	fmt.Fprintf(p.output, format, args...)
}

func (p *CppPrinter) printAppTypeTerm(term ir.IrTerm) {
	arg, types := term.AppTypes()

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

func (p *CppPrinter) printAppTermTerm(term ir.IrTerm) {
	id, types, arg := term.AppArgs()

	if id.Is(ir.VarTerm) && id.Var.ID == "ifthen" {
		condition, then := arg.Tuple.Elems[0], arg.Tuple.Elems[1]

		p.printf("if (")
		p.withLastTerm(false, func() { p.PrintTerm(condition) })
		p.printf(") ")
		p.PrintTerm(then)
		return
	}

	if id.Is(ir.VarTerm) && id.Var.ID == "ifelse" {
		condition, then, elseTerm := arg.Tuple.Elems[0], arg.Tuple.Elems[1], arg.Tuple.Elems[2]

		p.printf("if (")
		p.withLastTerm(false, func() { p.PrintTerm(condition) })
		p.printf(") ")
		p.PrintTerm(then)
		p.printf(" else ")
		p.PrintTerm(elseTerm)
		return
	}

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

func (p *CppPrinter) printAssignTerm(term ir.IrTerm) {
	p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
	p.printf(" = ")
	p.PrintTerm(term.Assign.Arg)
}

func (p *CppPrinter) printBlockTerm(term ir.IrTerm) {
	c := term.Block

	p.printf("{\n")

	for _, term := range c.Terms[:len(c.Terms)-1] {
		p.withLastTerm(false, func() {
			p.PrintTerm(term)
			p.printf(";")
		})
	}

	p.PrintTerm(c.Terms[len(c.Terms)-1])
	p.printf(";")

	p.printf("}\n")
}

func (p *CppPrinter) printConstTerm(term ir.IrTerm) {
	switch {
	case term.Const.Is(ir.IntLiteral):
		p.printf("%d", *term.Const.Int)
	case term.Const.Is(ir.FloatLiteral):
		p.printf("%d.%d", term.Const.Float.Integer, term.Const.Float.Decimal)
	case term.Const.Is(ir.StrLiteral):
		p.printf(`"%s"`, *term.Const.Str)
	}
}

func (p *CppPrinter) printIfTerm(term ir.IrTerm) {
	c := term.If

	p.printf("if (")
	p.withLastTerm(false, func() {
		p.PrintTerm(c.Condition)
	})
	p.printf(") ")
	p.PrintTerm(c.Then)
	if c.Else != nil {
		p.printf(" else ")
		p.PrintTerm(*c.Else)
	}
}

func (p *CppPrinter) printInjectionTerm(term ir.IrTerm) {
	c := term.Injection

	p.printType(c.VariantType)
	p.printf("{")
	if c.TagIndex != nil {
		p.printf("std::in_place_index<%d>, ", *c.TagIndex)
	}
	p.PrintTerm(c.Value)
	p.printf("}")
}

func (p *CppPrinter) printLambdaTerm(term ir.IrTerm) {
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
	p.withLastTerm(true, func() { p.PrintTerm(body) })
}

func (p *CppPrinter) printAliasDecl(id string, typ ir.IrType) {
	switch typ.Case {
	case ir.LambdaType:
		tvars := typ.LambdaVars()
		p.printf("template <")
		ir.Interleave(tvars, func() { p.printf(", ") }, func(_ int, tvar string) {
			p.printf("typename %s", tvar)
		})
		p.printf("> ")
		p.printAliasDecl(id, typ.LambdaBody())

	case ir.StructType:
		p.printf("struct %s {\n", id)
		for _, field := range typ.Fields() {
			p.printType(field.Type)
			p.printf(" %s;\n", field.ID)
		}
		p.printf("}")

	case ir.TupleType:
		p.printf("struct %s : ", id)
		p.withBindPosition(func() { p.printType(typ) })
		p.printf("{ %s(const ", id)
		p.withBindPosition(func() { p.printType(typ) })
		p.printf("& arg) : ")
		p.withBindPosition(func() { p.printType(typ) })
		p.printf("(arg) {}}")

	case ir.VariantType:
		p.printf("struct %s : ", id)
		p.printType(typ)
		p.printf("{}")

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
		ir.Interleave(args, func() { p.printf(", ") }, func(_ int, arg ir.IrType) {
			p.printType(arg)
		})
		p.printf(">")

	case typ.Is(ir.ArrayType):
		p.printf("std::array<")
		p.printType(typ.Array.ElemType)
		p.printf(", %d>", typ.Array.Size)

	case typ.Is(ir.ExistVarType):
		p.printf("%s", typ)

	case typ.Is(ir.ForallType):
		tvars := typ.ForallVars()
		p.printf("template <")
		ir.Interleave(tvars, func() { p.printf(", ") }, func(_ int, tvar string) {
			p.printf("typename %s", tvar)
		})
		p.printf("> ")
		p.printType(typ.ForallBody())

	case typ.Is(ir.FunType):
		c := typ.Fun

		p.printf("std::function<")
		p.withBindPosition(func() { p.printType(c.Ret) })
		p.printf("(")
		p.printType(c.Arg)
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

	case typ.Is(ir.StructType):
		if anonymousType, ok := p.anonymousTypes[p.hashType(typ)]; ok {
			p.printType(anonymousType.nameType)
			return
		}

		p.printf("struct {")
		ir.Interleave(typ.Fields(), func() { p.printf(" ") }, func(_ int, field ir.StructField) {
			p.printType(field.Type)
			p.printf(" %s;", field.ID)
		})
		p.printf("}")

	case typ.Is(ir.TupleType) && p.position == TypePosition:
		ir.Interleave(typ.Tuple.Elems, func() { p.printf(", ") }, func(_ int, elem ir.IrType) {
			p.printType(elem)
		})

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
			ir.Interleave(tuple.Elems, func() { p.printf(", ") }, func(_ int, elem ir.IrType) {
				p.printType(elem)
			})
			p.printf(">")
		}

	case typ.Is(ir.VariantType):
		p.printf("std::variant<")
		ir.Interleave(typ.Tags(), func() { p.printf(", ") }, func(_ int, tag ir.VariantTag) {
			p.withBindPosition(func() { p.printType(tag.Type) })
			p.printf("/* %s */", toID(tag.ID))
		})
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
				p.printf("template <")
				ir.Interleave(tvars, func() { p.printf(", ") }, func(_ int, tvar string) {
					p.printf("typename %s", tvar)
				})
				p.printf("> struct %s", id)

			default:
				p.printf("struct %s", id)
			}
			p.printf(";\n")
		})
		return
	}

	switch typ := decl.Term.Type; typ.Case {
	case ir.AppType, ir.ArrayType, ir.ExistVarType, ir.NameType, ir.TupleType, ir.VarType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			p.withBindPosition(func() {
				p.printType(typ)
				p.printf(" %s", id)
			})
		})

	case ir.ForallType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			tvars := typ.ForallVars()
			p.printf("template <")
			ir.Interleave(tvars, func() { p.printf(", ") }, func(_ int, tvar string) {
				p.printf("typename %s", tvar)
			})
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
	p.printf("#include <cmath>\n")
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
		p.printf("import %s;\n", toModuleName(imp.ModuleID))
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
				p.printf("template <")
				ir.Interleave(typeVars, func() { p.printf(", ") }, func(_ int, varkind ir.VarKind) {
					p.printf("typename %s", varkind.Var)
				})
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
		ir.Interleave(function.Args, func() { p.printf(", ") }, func(_ int, arg ir.FunctionArg) {
			p.withBindPosition(func() { p.printType(arg.Type) })
			p.printf(" %s", arg.ID)
		})

		p.printf(")\n")
		p.withLastTerm(true, func() { p.PrintTerm(function.Body) })
		p.printf("\n")
	})
}

func (p *CppPrinter) printLetTerm(term ir.IrTerm) {
	if !term.Is(ir.LetTerm) {
		panic(fmt.Errorf("expected %T", ir.LetTerm))
	}

	c := term.Let

	// There's no type (e.g., std::function) in C++20 for polymorphic
	// lambdas, so 'auto' must be used instead.
	//
	// For example:
	//   auto id = []<typename T>(T x) { return x; };
	auto := c.Value.Is(ir.TypeAbsTerm)

	p.withAutoType(auto, func() {
		p.withBindPosition(func() {
			p.printType(*c.VarType)
			p.printf(" %s", c.Var)
		})
	})

	if isCppStatement(c.Value) {
		p.withVarDestination(c.Var, func() {
			p.printf(";\n")
			p.PrintTerm(c.Value)
		})
	} else {
		p.printf(" = ")
		p.PrintTerm(c.Value)
	}
}

func (p *CppPrinter) printMatchTerm(term ir.IrTerm) {
	if !term.Is(ir.MatchTerm) {
		panic(fmt.Errorf("expected %T", ir.MatchTerm))
	}

	c := term.Match

	variantID := p.genID()

	p.printf("{")
	p.printf("auto %s = ", variantID)
	p.withLastTerm(false, func() {
		p.PrintTerm(c.Term)
	})
	p.printf(";\n")
	p.printf("switch (%s.index()) {\n", variantID)
	for _, arm := range c.Arms {
		p.printf("case %d: {", *arm.Index)
		p.printf("auto &%s = std::get<%d>(%s);\n", arm.Arg, *arm.Index, variantID)
		p.PrintTerm(arm.Body)
		p.printf(";\n")
		p.printf("}\n")
	}
	p.printf("} }")
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

func (p *CppPrinter) printReturnTerm(term ir.IrTerm) {
	if !term.Is(ir.ReturnTerm) {
		panic(fmt.Errorf("expected %T", ir.ReturnTerm))
	}

	c := term.Return

	p.printf("return ")
	p.PrintTerm(c.Expr)
	p.printf(";")
}

func (p *CppPrinter) printTupleTerm(term ir.IrTerm) {
	if !term.Is(ir.TupleTerm) {
		panic(fmt.Errorf("expected %T", ir.TupleTerm))
	}

	if p.position == BindPosition {
		p.printf("std::tie(")
	} else if len(term.Tuple.Elems) == 0 {
		p.printf("std::monostate(")
	} else {
		p.printf("std::make_tuple(")
	}

	ir.Interleave(term.Tuple.Elems, func() { p.printf(", ") }, func(_ int, elem ir.IrTerm) {
		p.PrintTerm(elem)
	})

	p.printf(")")
}

func (p *CppPrinter) printSetTerm(term ir.IrTerm) error {
	if !term.Is(ir.SetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.SetTerm, ir.SetTerm))
	}

	c := term.Set

	switch {
	case term.Type.Is(ir.StructType):
		structID := p.genID()

		p.printf("([&, %s = ", structID)
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

func (p *CppPrinter) printStructTerm(term ir.IrTerm) {
	if !term.Is(ir.StructTerm) {
		panic(fmt.Errorf("expected %T %d", ir.StructTerm, ir.StructTerm))
	}

	c := term.Struct

	p.printf("{")
	ir.Interleave(c.Values, func() { p.printf(", ") }, func(_ int, field ir.LabelValue) {
		p.printf(".%s = ", field.Label)
		p.PrintTerm(field.Value)
	})
	p.printf("}")
}

func (p *CppPrinter) PrintTerm(term ir.IrTerm) {
	switch {
	case term.Is(ir.AppTypeTerm):
		p.handleLastTerm(func() { p.printAppTypeTerm(term) })

	case term.Is(ir.AppTermTerm):
		if isCppStatement(term) {
			p.printAppTermTerm(term)
		} else {
			p.handleLastTerm(func() { p.printAppTermTerm(term) })
		}

	case term.Is(ir.AssignTerm):
		p.handleLastTerm(func() { p.printAssignTerm(term) })

	case term.Is(ir.BlockTerm):
		p.printBlockTerm(term)

	case term.Is(ir.ConstTerm):
		p.handleLastTerm(func() { p.printConstTerm(term) })

	case term.Is(ir.IfTerm):
		p.printIfTerm(term)

	case term.Is(ir.InjectionTerm):
		p.handleLastTerm(func() { p.printInjectionTerm(term) })

	case term.Is(ir.LambdaTerm) || term.Is(ir.TypeAbsTerm):
		p.handleLastTerm(func() {
			p.withLastTerm(false, func() { p.printLambdaTerm(term) })
		})

	case term.Is(ir.LetTerm):
		p.withLastTerm(false, func() { p.printLetTerm(term) })

	case term.Is(ir.MatchTerm):
		p.printMatchTerm(term)

	case term.Is(ir.ProjectionTerm):
		p.handleLastTerm(func() {
			if err := p.printProjectionTerm(term); err != nil {
				p.err = errors.Join(p.err, err)
			}
		})

	case term.Is(ir.ReturnTerm):
		p.withLastTerm(false, func() { p.printReturnTerm(term) })

	case term.Is(ir.TupleTerm):
		p.handleLastTerm(func() { p.printTupleTerm(term) })

	case term.Is(ir.SetTerm):
		p.handleLastTerm(func() {
			if err := p.printSetTerm(term); err != nil {
				p.err = errors.Join(p.err, err)
			}
		})

	case term.Is(ir.StructTerm):
		p.handleLastTerm(func() { p.printStructTerm(term) })

	case term.Is(ir.VarTerm):
		p.handleLastTerm(func() { p.printf("%s", toID(term.Var.ID)) })

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

	if err := p.recordAnonymousTypesFromUnit(&unit); err != nil {
		return err
	}

	p.printDecls(unit.Decls)
	for _, function := range unit.Functions {
		p.printFunction(function)
	}

	return nil
}

func newCppPrinter(output io.Writer) *CppPrinter {
	return &CppPrinter{
		output,
		0, /* idgen */
		TypePosition,
		false,                      /* autoType */
		false,                      /* lastTerm */
		"",                         /* varDestination */
		map[string]anonymousType{}, /* anonymousTypes */
		nil,                        /* err */
	}
}

func printUnitToCpp(unit ir.IrUnit, output io.Writer) error {
	printer := newCppPrinter(output)
	printErr := printer.printUnit(unit)
	return errors.Join(printErr, printer.err)
}
