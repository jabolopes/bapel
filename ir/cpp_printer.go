package ir

import (
	"fmt"
	"io"

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

func (p *CppPrinter) printf(format string, args ...any) {
	fmt.Fprintf(p.out(), format, args...)
}

func (p *CppPrinter) printCall(id string, arg IrTerm) {
	p.printf("%s(", toID(id))
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
		panic(fmt.Errorf("unhandled token %d", token.Case))
	}
}

func (p *CppPrinter) printType(typ IrType) {
	switch {
	case typ.Case == ArrayType:
		fmt.Fprintf(p.out(), "std::array<")
		p.printType(typ.Array.ElemType)
		fmt.Fprintf(p.out(), ", %d>", typ.Array.Size)

	case typ.Case == IntType:
		switch typ.Int {
		case I8:
			fmt.Fprintf(p.out(), "char")
		case I16:
			fmt.Fprintf(p.out(), "int16_t")
		case I32:
			fmt.Fprintf(p.out(), "int32_t")
		case I64:
			fmt.Fprintf(p.out(), "int64_t")
		}

	case typ.Case == NameType:
		p.printf("struct %s", toID(typ.Name))

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
		panic(fmt.Errorf("printType: Unhandled case %d", typ.Case))
	}
}

func (p *CppPrinter) printDecl(decl IrDecl) {
	switch decl.Type.Case {
	case FunType:
		typ := decl.Type.Fun

		p.withBindPosition(func() { p.printType(typ.Ret) })

		// Print id.
		//
		// TODO: Handle namespacing.
		p.printf(" %s(", decl.ID)

		p.printType(typ.Arg)
		p.printf(")")

	case IntType:
		p.printType(decl.Type)
		p.printf(" %s", decl.ID)

	case StructType:
		// TODO: Handle namespacing.
		p.printf("struct %s", decl.ID)

	default:
		panic(fmt.Errorf("unhandled IrType %d", decl.Type.Case))
	}
}

func (p *CppPrinter) PrintDef(decl IrDecl) {
	switch decl.Type.Case {
	case StructType:
		p.printf("struct %s {\n", decl.ID)
		for _, field := range decl.Type.Fields() {
			p.printType(field.Type)
			p.printf(" %s;\n", field.ID)
		}
		p.printf("};\n")

	case IntType:
		p.printType(decl.Type)
		p.printf(" %s;\n", decl.ID)

	default:
		panic(fmt.Errorf("unhandled IrType %d", decl.Type.Case))
	}
}

func (p *CppPrinter) PrintTerm(term IrTerm) {
	switch term.Case {
	case AssignTerm:
		p.withBindPosition(func() { p.PrintTerm(term.Assign.Ret) })
		p.printf(" = ")
		p.PrintTerm(term.Assign.Arg)

	case CallTerm:
		p.printCall(term.Call.ID, term.Call.Arg)

	case IfTerm:
		p.printf("if (")
		if !term.If.Then {
			p.printf("!")
		}
		p.PrintTerm(term.If.Condition)
		p.printf(") {\n")

	case IndexGetTerm:
		if len(term.IndexGet.Field) == 0 {
			p.PrintTerm(term.IndexGet.Term)
			p.printf("[")
			p.PrintTerm(term.IndexGet.Index)
			p.printf("]")
		} else {
			p.PrintTerm(term.IndexGet.Term)
			p.printf(".%s", term.IndexGet.Field)
		}

	case IndexSetTerm:
		if len(term.IndexSet.Field) == 0 {
			p.PrintTerm(term.IndexSet.Ret)
			p.printf("[")
			p.PrintTerm(term.IndexSet.Index)
			p.printf("] = ")
			p.PrintTerm(term.IndexSet.Arg)
		} else {
			p.PrintTerm(term.IndexSet.Ret)
			p.printf(".%s = ", term.IndexSet.Field)
			p.PrintTerm(term.IndexSet.Arg)
		}

	case StatementTerm:
		p.PrintTerm(term.Statement.Term)
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
