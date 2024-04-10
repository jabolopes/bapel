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

func (p *CppPrinter) out() io.Writer {
	return p.output
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
	p.printf(" }\n")
}

func (p *CppPrinter) printf(format string, args ...any) {
	fmt.Fprintf(p.out(), format, args...)
}

func (p *CppPrinter) printCall(id string, types []IrType, arg IrTerm) {
	p.printf("%s", toID(id))
	if !isOperator(id) && len(types) > 0 {
		p.printf("<")
		p.printType(types[0])
		for _, typ := range types[1:] {
			p.printf(", ")
			p.printType(typ)
		}
		p.printf(">")
	}
	p.printf("(")
	p.PrintTerm(arg)
	p.printf(")")
}

func (p *CppPrinter) printToken(token parser.Token) {
	switch token.Case {
	case parser.IDToken:
		fmt.Fprintf(p.out(), "%s", toID(token.Text))
	case parser.NumberToken:
		fmt.Fprintf(p.out(), "%s", token.Text)
	default:
		panic(fmt.Errorf("unhandled %d %d", token.Case, token.Case))
	}
}

func (p *CppPrinter) printNamedType(name, value IrType) {
	switch value.Case {
	case StructType:
		p.printf("struct ")
		p.printType(name)
		p.printf(" {\n")
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
	case typ.Case == AliasType:
		p.printNamedType(typ.Alias.Name, typ.Alias.Value)

	case typ.Case == ArrayType:
		fmt.Fprintf(p.out(), "std::array<")
		p.printType(typ.Array.ElemType)
		fmt.Fprintf(p.out(), ", %d>", typ.Array.Size)

	case typ.Case == ForallType:
		// Print type variables.
		p.printf("template <typename %s", typ.Forall.Vars[0])
		for _, tvar := range typ.Forall.Vars[1:] {
			p.printf(", typename %s", tvar)
		}
		p.printf("> ")
		p.printType(typ.Forall.Type)

	case typ.Case == NameType:
		switch typ.Name {
		case "i8":
			fmt.Fprintf(p.out(), "char")
		case "i16":
			fmt.Fprintf(p.out(), "int16_t")
		case "i32":
			fmt.Fprintf(p.out(), "int32_t")
		case "i64":
			fmt.Fprintf(p.out(), "int64_t")
		default:
			p.printf("%s", toID(typ.Name))
		}

	case typ.Case == TupleType && p.position == TypePosition:
		tuple := typ.Tuple
		if len(tuple) > 0 {
			p.printType(tuple[0])
			for _, elem := range tuple[1:] {
				p.printf(", ")
				p.printType(elem)
			}
		}

	case typ.Case == TupleType && p.position == BindPosition:
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

	case typ.Case == VarType:
		p.printf("%s", typ.Var)

	default:
		panic(fmt.Errorf("printType: unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) printDecl(decl IrDecl) {
	if decl.Case == TypeDecl {
		p.printType(decl.Type())
		return
	}

	switch typ := decl.Type(); typ.Case {
	case ForallType:
		p.printInNamespace(decl.Term.ID, func(id string) {
			// Print type variables.
			p.printf("template <typename %s", typ.Forall.Vars[0])
			for _, tvar := range typ.Forall.Vars[1:] {
				p.printf(", typename %s", tvar)
			}
			p.printf("> ")
			p.printDecl(NewTermDecl(id, typ.Forall.Type))
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

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) PrintDef(decl IrDecl) {
	if decl.Case == TypeDecl {
		p.printType(decl.Type())
		return
	}

	switch typ := decl.Type(); typ.Case {
	case NameType:
		p.printType(typ)
		p.printf(" %s;\n", decl.Term.ID)

	case StructType:
		p.printf("struct %s {\n", decl.Term.ID)
		for _, field := range typ.Fields() {
			p.printType(field.Type)
			p.printf(" %s;\n", field.ID)
		}
		p.printf("};\n")

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) PrintTerm(term IrTerm) {
	switch term.Case {
	case AssignTerm:
		p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
		p.printf(" = ")
		p.PrintTerm(term.Assign.Arg)

	case BlockTerm:
		c := term.Block
		p.printf("{\n")
		for _, term := range c.Terms {
			p.PrintTerm(term)
		}
		p.printf("}\n")

	case CallTerm:
		p.printCall(term.Call.ID, term.Call.Types, term.Call.Arg)

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
		if len(term.IndexGet.Field) == 0 {
			p.PrintTerm(term.IndexGet.Obj)
			p.printf("[")
			p.PrintTerm(term.IndexGet.Index)
			p.printf("]")
		} else {
			p.PrintTerm(term.IndexGet.Obj)
			p.printf(".%s", term.IndexGet.Field)
		}

	case IndexSetTerm:
		if len(term.IndexSet.Field) == 0 {
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
		p.printDecl(c.Decl)

	case StatementTerm:
		c := term.Statement
		p.PrintTerm(c.Term)
		p.printf(";\n")

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

	case WidenTerm:
		// TODO: Insert a cast.
		p.PrintTerm(term.Widen.Term)

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
