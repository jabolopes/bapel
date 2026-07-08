package comp

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type Position int

const (
	TypePosition = Position(iota)
	BindPosition
)

type PrinterMode int

const (
	ModePublicHeader PrinterMode = iota
	ModePrivateHeader
	ModeSource
)

func toHeaderPath(moduleID ir.ModuleID) string {
	return strings.Replace(moduleID.Name, ir.ModuleIDSeparator, "/", -1) + ".h"
}

func (p *CppPrinter) isTypeDecl(decl ir.IrDecl) bool {
	return decl.Is(ir.NameDecl) || decl.Is(ir.AliasDecl)
}

func (p *CppPrinter) toID(id string) string {
	if strings.Contains(id, ir.NamespaceSeparator) {
		parts := strings.Split(id, ir.NamespaceSeparator)
		if len(parts) >= 2 {
			typePrefix := strings.Join(parts[:len(parts)-1], ir.NamespaceSeparator)
			if decl, ok := p.findDecl(typePrefix); ok && p.isTypeDecl(decl) {
				methodName := parts[len(parts)-1]
				return "::" + strings.Replace(inherentCppName(typePrefix), ir.NamespaceSeparator, "::", -1) + "::" + methodName
			}
		}
		return "::" + strings.Replace(id, ir.NamespaceSeparator, "::", -1)
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

func countAliasTypeVars(typ ir.IrType) int {
	if typ.Is(ir.LambdaType) {
		return len(typ.LambdaVars())
	}
	return 0
}

func isCppStatement(term ir.IrTerm) bool {
	switch term.Case {
	case ir.AssignTerm, ir.BlockTerm, ir.LetTerm, ir.MatchTerm, ir.ReturnTerm:
		return true
	}

	if term.Is(ir.AppTermTerm) {
		id, _, _ := term.AppArgs()
		if id.Is(ir.VarTerm) && (id.Var.ID == "ifthen" || id.Var.ID == "ifelse" || id.Var.ID == "core::for") {
			return true
		}
	}

	return false
}

type CppPrinter struct {
	output         io.Writer
	Mode           PrinterMode
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
	unit           *ir.IrUnit
}

func (p *CppPrinter) genID() string {
	id := fmt.Sprintf("__v_%d", p.idgen)
	p.idgen++
	return id
}

func (p *CppPrinter) findDecl(id string) (ir.IrDecl, bool) {
	for _, decl := range p.unit.Decls {
		if decl.ID() == id {
			return decl, true
		}
	}
	for _, decl := range p.unit.ImportDecls {
		if decl.ID() == id {
			return decl, true
		}
	}
	for _, decl := range p.unit.ImplDecls {
		if decl.ID() == id {
			return decl, true
		}
	}
	return ir.IrDecl{}, false
}

func (p *CppPrinter) findTraitDecl(id string) (ir.IrDecl, bool) {
	for _, decl := range p.unit.Decls {
		if decl.ID() == id && decl.Is(ir.TraitDecl) {
			return decl, true
		}
	}
	for _, decl := range p.unit.ImportDecls {
		if decl.ID() == id && decl.Is(ir.TraitDecl) {
			return decl, true
		}
	}
	for _, decl := range p.unit.ImplDecls {
		if decl.ID() == id && decl.Is(ir.TraitDecl) {
			return decl, true
		}
	}
	return ir.IrDecl{}, false
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

	if id.Is(ir.VarTerm) && id.Var.ID == "core::for" {
		cond := arg.Tuple.Elems[0]
		bodyLambda := arg.Tuple.Elems[1]
		_, _, body := bodyLambda.ToFunction()

		p.printf("while (")
		p.withLastTerm(false, func() { p.PrintTerm(cond) })
		p.printf(") ")
		p.withLastTerm(false, func() { p.PrintTerm(body) })
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

	// Trait method handling
	if id.Is(ir.VarTerm) && strings.Contains(id.Var.ID, ir.NamespaceSeparator) {
		parts := strings.Split(id.Var.ID, ir.NamespaceSeparator)
		if len(parts) >= 2 {
			traitName := strings.Join(parts[:len(parts)-1], ir.NamespaceSeparator)
			methodName := parts[len(parts)-1]
			if _, ok := p.findTraitDecl(traitName); ok {
				p.printf("::%s<", strings.Replace(traitCppName(traitName), ir.NamespaceSeparator, "::", -1))
				p.withBindPosition(func() {
					ir.Interleave(types, func() { p.printf(", ") }, func(_ int, typ ir.IrType) {
						p.printType(typ)
					})
				})
				p.printf(">::%s(", methodName)
				if arg.Is(ir.TupleTerm) {
					ir.Interleave(arg.Tuple.Elems, func() { p.printf(", ") }, func(_ int, t ir.IrTerm) {
						p.PrintTerm(t)
					})
				} else {
					p.PrintTerm(arg)
				}
				p.printf(")")
				return
			}
		}
	}

	// Generic inherent method handling
	if id.Is(ir.VarTerm) && !ir.IsOperator(id.Var.ID) && strings.Contains(id.Var.ID, ir.NamespaceSeparator) {
		parts := strings.Split(id.Var.ID, ir.NamespaceSeparator)
		typePrefix := strings.Join(parts[:len(parts)-1], ir.NamespaceSeparator)
		methodName := parts[len(parts)-1]

		if decl, ok := p.findDecl(typePrefix); ok && (decl.Is(ir.NameDecl) || decl.Is(ir.AliasDecl)) {
			var arity int
			if decl.Is(ir.NameDecl) {
				arity = countTypeVars(decl.Name.Kind)
			} else {
				arity = countAliasTypeVars(decl.Alias.Type)
			}

			if len(types) >= arity {
				typeArgs := types[:arity]
				methodArgs := types[arity:]

				mappedType := typePrefix
				if p.isTypeDecl(decl) {
					mappedType = inherentCppName(typePrefix)
				}

				p.printf("::%s", strings.Replace(mappedType, ir.NamespaceSeparator, "::", -1))
				if arity > 0 {
					p.printf("<")
					p.withBindPosition(func() {
						ir.Interleave(typeArgs, func() { p.printf(", ") }, func(_ int, typ ir.IrType) {
							p.printType(typ)
						})
					})
					p.printf(">")
				}
				p.printf("::%s", methodName)

				if len(methodArgs) > 0 {
					p.printf("<")
					p.withBindPosition(func() {
						ir.Interleave(methodArgs, func() { p.printf(", ") }, func(_ int, typ ir.IrType) {
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
				return
			}
		}
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

	lastTerm := c.Terms[len(c.Terms)-1]
	p.PrintTerm(lastTerm)
	p.printf(";}\n")
}

func (p *CppPrinter) printConstTerm(term ir.IrTerm) {
	switch {
	case term.Const.Is(ir.IntLiteral):
		p.printf("%d", *term.Const.Int)
	case term.Const.Is(ir.FloatLiteral):
		p.printf("%d.%d", term.Const.Float.Integer, term.Const.Float.Decimal)
	case term.Const.Is(ir.RuneLiteral):
		p.printf(`'%s'`, *term.Const.Rune)
	case term.Const.Is(ir.StrLiteral):
		p.printf(`"%s"`, *term.Const.Str)
	}
}

func (p *CppPrinter) printInjectionTerm(term ir.IrTerm) {
	c := term.Injection

	p.printType(c.VariantType)
	p.printf("(")
	if c.TagIndex != nil {
		p.printf("std::in_place_index<%d>, ", *c.TagIndex)
	}
	p.PrintTerm(c.Value)
	p.printf(")")
}

func (p *CppPrinter) printLambdaTerm(term ir.IrTerm) {
	tvars, args, body := term.ToFunction()
	p.printf("[&]")

	// Print type abstraction types.
	if len(tvars) > 0 {
		p.printf("<")
		ir.Interleave(tvars, func() { p.printf(", ") }, func(_ int, tvar ir.TypeParam) {
			p.printf("typename %s", tvar.Var)
		})
		p.printf(">")
	}

	// Print abstraction arguments and types.
	arglessLambda := len(args) == 1 && args[0].Type.Is(ir.TupleType) && len(args[0].Type.Elems()) == 0

	p.printf("(")
	if !arglessLambda {
		ir.Interleave(args, func() { p.printf(", ") }, func(_ int, arg ir.FunctionArg) {
			p.printType(arg.Type)
			p.printf(" %s", p.toID(arg.ID))
		})
	}
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

	case ir.NameType:
		p.printf("using %s =", id)
		p.printType(typ)

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
		p.printf(" { using ")
		p.printType(typ)
		p.printf("::variant; }")

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
			name := typ.Name
			if !strings.Contains(name, ir.NamespaceSeparator) {
				if _, ok := p.findDecl(name); ok {
					p.printf("::%s", p.toID(name))
					break
				}
			}
			p.printf("%s", p.toID(name))
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
			p.printf("/* %s */", p.toID(tag.ID))
		})
		p.printf(">")

	case typ.Is(ir.VarType):
		p.printf("%s", typ.Var)

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (p *CppPrinter) printDecl(decl ir.IrDecl) {

	if p.Mode == ModeSource {
		return
	}
	if p.Mode == ModePublicHeader && !decl.Export {
		return
	}
	if p.Mode == ModePrivateHeader && decl.Export {
		return
	}

	if decl.Is(ir.TraitDecl) {
		p.printTraitDecl(decl)
		return
	}

	if decl.Is(ir.NameDecl) {
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
			tvars := typ.ForallTypeParams()
			p.printTemplateParams(tvars, false)
			p.printDecl(ir.NewTermDecl(id, typ.ForallBody(), decl.Export))
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

	if p.Mode == ModeSource {
		return
	}
	if p.Mode == ModePublicHeader && !decl.Export {
		return
	}
	if p.Mode == ModePrivateHeader && decl.Export {
		return
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


func (p *CppPrinter) printModuleTop(unit ir.IrUnit) {
	switch p.Mode {
	case ModePublicHeader:
		p.printf("#pragma once\n")
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

	case ModePrivateHeader:
		p.printf("#pragma once\n")
		p.printf("\n")
		headerPath := toHeaderPath(unit.ModuleID)
		p.printf("#include \"%s\"\n", headerPath)
		p.printf("\n")

	case ModeSource:
		p.printf("\n")
		privateHeaderPath := strings.TrimSuffix(toHeaderPath(unit.ModuleID), ".h") + "_private.h"
		p.printf("#include \"%s\"\n", privateHeaderPath)
		p.printf("\n")
	}
}

func (p *CppPrinter) printImports(imports []ir.IrImport) {
	if p.Mode != ModePublicHeader {
		return
	}
	p.printf("\n")
	for _, imp := range imports {
		headerPath := toHeaderPath(imp.ModuleID)
		p.printf("#include \"%s\"\n", headerPath)
	}
	p.printf("\n")
}

func (p *CppPrinter) printImpls(impls []ir.IrImpl) {
	if p.Mode != ModePublicHeader {
		return
	}
	p.printf("\n")
	for _, impl := range impls {
		if path.Ext(impl.RelativeFilename.Value) == ".h" {
			p.printf("#include \"%s\"\n", impl.RelativeFilename.Value)
		}
	}
	p.printf("\n")
}


func (p *CppPrinter) printFunctionSignature(function ir.IrFunction) {
	p.printInNamespace(function.ID, func(id string) {
		p.printTemplateParams(function.TypeParams, false)
		p.withBindPosition(func() { p.printType(function.RetType) })
		p.printf(" %s(", id)
		ir.Interleave(function.Args, func() { p.printf(", ") }, func(_ int, arg ir.FunctionArg) {
			p.printType(arg.Type)
			p.printf(" %s", arg.ID)
		})
		p.printf(");\n")
	})
}

func (p *CppPrinter) printFunctionFull(function ir.IrFunction) {
	p.printInNamespace(function.ID, func(id string) {
		p.printTemplateParams(function.TypeParams, true)
		p.withBindPosition(func() { p.printType(function.RetType) })
		p.printf(" %s(", id)
		ir.Interleave(function.Args, func() { p.printf(", ") }, func(_ int, arg ir.FunctionArg) {
			p.printType(arg.Type)
			p.printf(" %s", arg.ID)
		})
		p.printf(")\n")
		p.withLastTerm(true, func() { p.PrintTerm(function.Body) })
		p.printf("\n")
	})
}

func (p *CppPrinter) printFunction(function ir.IrFunction) {
	isTemplate := len(function.TypeParams) > 0
	isPub := function.Export

	switch p.Mode {
	case ModePublicHeader:
		if !isPub {
			return
		}
		if isTemplate {
			p.printFunctionFull(function)
		} else {
			p.printFunctionSignature(function)
		}
	case ModePrivateHeader:
		if isPub {
			return
		}
		if isTemplate {
			p.printFunctionFull(function)
		} else {
			p.printFunctionSignature(function)
		}
	case ModeSource:
		if isTemplate {
			return
		}
		p.printFunctionFull(function)
	}
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

	objType := term.Type
	if c.ReducedType != nil {
		objType = c.ReducedType
	}

	if objType.Is(ir.StructType) {
		_, field, err := objType.FieldByLabel(c.Label)
		if err != nil {
			return err
		}

		p.PrintTerm(c.Term)
		p.printf(".%s", field.ID)
		return nil
	}

	var index int
	var err error

	if objType.Is(ir.TupleType) {
		index, _, err = objType.ElemByLabel(c.Label)
	} else if objType.Is(ir.VariantType) {
		index, _, err = objType.TagByLabel(c.Label)
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

	objType := term.Type
	if c.ReducedType != nil {
		objType = c.ReducedType
	}

	switch {
	case objType.Is(ir.StructType):
		structID := p.genID()

		p.printf("([&, %s = ", structID)
		p.PrintTerm(c.Term)
		p.printf("]() mutable {\n")
		for _, lv := range c.Values {
			_, field, err := objType.FieldByLabel(lv.Label)
			if err != nil {
				return err
			}

			p.printf("%s.%s = ", structID, field.ID)
			p.PrintTerm(lv.Value)
			p.printf(";\n")
		}
		p.printf("return %s;\n", structID)
		p.printf("})()")

	case objType.Is(ir.TupleType):
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
		panic(fmt.Errorf("unhandled type %s", *objType))
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
		p.handleLastTerm(func() { p.printf("%s", p.toID(term.Var.ID)) })

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
	p.printModuleTop(unit)
	p.printImpls(unit.Impls)
	p.printImports(unit.Imports)

	if err := p.recordAnonymousTypesFromUnit(&unit); err != nil {
		return err
	}

	p.printDecls(unit.Decls)
	if p.Mode == ModePrivateHeader {
		p.printDecls(unit.ImplDecls)
	}
	for _, function := range unit.Functions {
		p.printFunction(function)
	}
	for _, impl := range unit.TraitImpls {
		p.printTraitImpl(impl)
	}

	return nil
}

func newCppPrinter(output io.Writer, mode PrinterMode, unit *ir.IrUnit) *CppPrinter {
	return &CppPrinter{
		output,
		mode,
		0, /* idgen */
		TypePosition,
		false,                      /* autoType */
		false,                      /* lastTerm */
		"",                         /* varDestination */
		map[string]anonymousType{}, /* anonymousTypes */
		nil,                        /* err */
		unit,
	}
}

func PrintUnitToCpp(unit ir.IrUnit, mode PrinterMode, output io.Writer) error {
	printer := newCppPrinter(output, mode, &unit)
	printErr := printer.printUnit(unit)
	return errors.Join(printErr, printer.err)
}

func (p *CppPrinter) printTraitDecl(decl ir.IrDecl) {
	p.printInNamespace(traitCppName(decl.Trait.ID), func(id string) {
		p.printf("template <typename Self")
		for _, tp := range decl.Trait.TypeParams {
			p.printf(", typename %s", tp.Var)
		}
		p.printf(">\n")
		p.printf("struct %s;\n", id)
	})
}

func inherentCppName(bapelName string) string {
	parts := strings.Split(bapelName, ir.NamespaceSeparator)
	runes := []string{"inherents"}
	parts = append(parts[:len(parts)-1], append(runes, parts[len(parts)-1])...)
	return strings.Join(parts, ir.NamespaceSeparator)
}

func (p *CppPrinter) printTraitImpl(impl ir.IrTraitImpl) {
	var exported bool
	if name := baseTypeName(impl.TypeName); name != "" {
		if decl, ok := p.findDecl(name); ok {
			exported = decl.Export
		}
	}

	if p.Mode == ModeSource {
		return
	}
	if p.Mode == ModePublicHeader && !exported {
		return
	}
	if p.Mode == ModePrivateHeader && exported {
		return
	}

	if impl.Case == ir.InherentImpl {
		baseName := baseTypeName(impl.TypeName)
		p.printInNamespace(inherentCppName(baseName), func(id string) {
			if len(impl.TypeParams) > 0 {
				p.printf("template <")
				ir.Interleave(impl.TypeParams, func() { p.printf(", ") }, func(_ int, tp ir.TypeParam) {
					p.printf("typename %s", tp.Var)
				})
				p.printf(">\n")
			}
			p.printf("struct %s {\n", id)
			p.printf("  %s() = delete;\n", id)
			p.printf("  using Self = ")
			p.withBindPosition(func() { p.printType(impl.TypeName) })
			p.printf(";\n")
			for _, m := range impl.Methods {
				p.printf("  static inline ")
				p.withBindPosition(func() { p.printType(m.RetType) })
				p.printf(" %s(", m.ID)
				ir.Interleave(m.Args, func() { p.printf(", ") }, func(_ int, arg ir.FunctionArg) {
					p.printType(arg.Type)
					p.printf(" %s", arg.ID)
				})
				p.printf(") ")
				p.withLastTerm(true, func() { p.PrintTerm(m.Body) })
				p.printf("\n")
			}
			p.printf("};\n")
		})
		return
	}

	traitName, err := impl.TraitType.TraitName()
	if err != nil {
		panic(err)
	}

	p.printInNamespace(traitCppName(traitName), func(id string) {
		if len(impl.TypeParams) > 0 {
			p.printf("template <")
			ir.Interleave(impl.TypeParams, func() { p.printf(", ") }, func(_ int, tp ir.TypeParam) {
				p.printf("typename %s", tp.Var)
			})
			p.printf(">\n")
		} else {
			p.printf("template <>\n")
		}
		p.printf("struct %s<", id)
		p.withBindPosition(func() { p.printType(impl.TypeName) })
		for _, arg := range impl.TraitType.AppArgs() {
			p.printf(", ")
			p.withBindPosition(func() { p.printType(arg) })
		}
		p.printf("> {\n")
		p.printf("  using Self = ")
		p.withBindPosition(func() { p.printType(impl.TypeName) })
		p.printf(";\n")

		for _, m := range impl.Methods {
			p.printf("  static inline ")
			p.withBindPosition(func() { p.printType(m.RetType) })
			p.printf(" %s(", m.ID)
			ir.Interleave(m.Args, func() { p.printf(", ") }, func(_ int, arg ir.FunctionArg) {
				p.printType(arg.Type)
				p.printf(" %s", arg.ID)
			})
			p.printf(") ")
			p.withLastTerm(true, func() { p.PrintTerm(m.Body) })
			p.printf("\n")
		}
		p.printf("};\n")
	})
}

func traitCppName(bapelName string) string {
	parts := strings.Split(bapelName, ir.NamespaceSeparator)
	runes := []string{"traits"}
	parts = append(parts[:len(parts)-1], append(runes, parts[len(parts)-1])...)
	return strings.Join(parts, ir.NamespaceSeparator)
}

func (p *CppPrinter) printTemplateParams(tvars []ir.TypeParam, isDefinition bool) {
	if len(tvars) == 0 {
		return
	}
	p.printf("template <")
	ir.Interleave(tvars, func() { p.printf(", ") }, func(_ int, vk ir.TypeParam) {
		p.printf("typename %s", vk.Var)
	})

	var constraints []string
	for _, vk := range tvars {
		for _, bound := range vk.Bounds {
			constraints = append(constraints, p.sfinaeConstraint(vk.Var, bound))
		}
	}

	if len(constraints) > 0 {
		if isDefinition {
			p.printf(", typename")
		} else {
			p.printf(", typename = std::enable_if_t<")
			p.printf("%s", strings.Join(constraints, " && "))
			p.printf(">")
		}
	}
	p.printf("> ")
}

func (p *CppPrinter) sfinaeConstraint(vkVar string, bound ir.IrType) string {
	var traitName string
	var args []ir.IrType
	if bound.Is(ir.NameType) {
		traitName = bound.Name
	} else if bound.Is(ir.AppType) {
		traitName = bound.App.Fun.Name
		args = bound.AppArgs()
	} else {
		panic(fmt.Errorf("invalid trait bound type: %T", bound))
	}

	cppName := traitCppName(traitName)
	cppName = "::" + cppName

	var argsStr []string
	argsStr = append(argsStr, vkVar)
	for _, arg := range args {
		var sb strings.Builder
		tempPrinter := newCppPrinter(&sb, p.Mode, p.unit)
		tempPrinter.printType(arg)
		argsStr = append(argsStr, sb.String())
	}

	return fmt.Sprintf("(sizeof(%s<%s>) > 0)", cppName, strings.Join(argsStr, ", "))
}




