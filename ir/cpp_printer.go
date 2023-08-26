package ir

import (
	"fmt"
	"io"

	"github.com/jabolopes/bapel/parser"
	"github.com/zyedidia/generic/stack"
)

type CppPrinter struct {
	output       io.Writer
	bindPosition *stack.Stack[bool]
}

func (p *CppPrinter) out() io.Writer {
	return p.output
}

func (p *CppPrinter) printf(format string, args ...any) {
	fmt.Fprintf(p.out(), format, args...)
}

func (p *CppPrinter) printCall(id string, args []IrTerm) {
	p.printf("%s(", toID(id))

	if len(args) > 0 {
		p.PrintTerm(args[0])
		for _, arg := range args[1:] {
			p.printf(", ")
			p.PrintTerm(arg)
		}
	}

	p.printf(")")
}

func (p *CppPrinter) printToken(token parser.Token) {
	switch token.Case {
	case parser.IDToken:
		fmt.Fprintf(p.out(), "%s", toID(token.Text))
	case parser.NumberToken:
		fmt.Fprintf(p.out(), "%s", token.Text)
	default:
		panic(fmt.Errorf("Unhandled token %d", token.Case))
	}
}

func (p *CppPrinter) printType(typ IrType) {
	switch {
	case typ.Case == ArrayType:
		fmt.Fprintf(p.out(), "std::array<")
		p.printType(typ.ArrayType.ElemType)
		fmt.Fprintf(p.out(), ", %d>", typ.ArrayType.Size)
	case typ.Case == FunType:
		panic(fmt.Errorf("printType: Unimplemented function type"))
	case typ.Case == IntType:
		switch typ.IntType {
		case I8:
			fmt.Fprintf(p.out(), "char")
		case I16:
			fmt.Fprintf(p.out(), "int16_t")
		case I32:
			fmt.Fprintf(p.out(), "int32_t")
		case I64:
			fmt.Fprintf(p.out(), "int64_t")
		}

	case typ.Case == TupleType && p.bindPosition.Peek():
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

	case typ.Case == TupleType && !p.bindPosition.Peek():
		tuple := typ.Tuple
		if len(tuple) > 0 {
			p.printType(tuple[0])
			for _, elem := range tuple[1:] {
				p.printf(", ")
				p.printType(elem)
			}
		}

	case typ.Case == IDType:
		p.printf("struct %s", toID(typ.IDType))
	default:
		panic(fmt.Errorf("printType: Unhandled case %d", typ.Case))
	}
}

func (p *CppPrinter) printDecl(decl IrDecl) {
	switch decl.Type.Case {
	case FunType:
		typ := decl.Type.FunType

		p.bindPosition.Push(true)
		p.printType(NewTupleType(typ.Rets))
		p.bindPosition.Pop()

		// Print id.
		//
		// TODO: Handle namespacing.
		p.printf(" %s(", decl.ID)

		p.printType(NewTupleType(typ.Args))
		p.printf(")")

	case IntType:
		p.printType(decl.Type)
		p.printf(" %s", decl.ID)

	case StructType:
		// TODO: Handle namespacing.
		p.printf("struct %s", decl.ID)

	default:
		panic(fmt.Errorf("Unhandled IrType %d", decl.Type.Case))
	}
}

func (p *CppPrinter) PrintTerm(term IrTerm) {
	switch term.Case {
	case AssignTerm:
		p.bindPosition.Push(true)
		p.PrintTerm(term.Assign.Ret)
		p.bindPosition.Pop()
		p.printf(" = ")
		p.PrintTerm(term.Assign.Arg)

	case CallTerm:
		p.printCall(term.Call.ID, term.Call.Args)

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

	case OpUnaryTerm:
		p.printf("%s ", term.OpUnary.ID)
		p.PrintTerm(term.OpUnary.Term)

	case OpBinaryTerm:
		p.PrintTerm(term.OpBinary.Left)
		p.printf(" %s ", term.OpBinary.ID)
		p.PrintTerm(term.OpBinary.Right)

	case StatementTerm:
		p.PrintTerm(term.Statement.Expr)
		p.printf(";\n")

	case TokenTerm:
		p.printToken(*term.Token)

	case TupleTerm:
		if p.bindPosition.Peek() {
			p.printf("std::tie(")
		}

		if len(term.Tuple) > 0 {
			p.PrintTerm(term.Tuple[0])
			for _, term := range term.Tuple[1:] {
				p.printf(", ")
				p.PrintTerm(term)
			}
		}

		if p.bindPosition.Peek() {
			p.printf(")")
		}

	case WidenTerm:
		// TODO: Insert a cast.
		p.PrintTerm(term.Widen.Term)

	default:
		panic(fmt.Errorf("Unhandled IrTerm %d", term.Case))
	}
}

func NewCppPrinter(output io.Writer) *CppPrinter {
	printer := &CppPrinter{
		output,
		stack.New[bool](), /* bindPosition */
	}
	printer.bindPosition.Push(false)
	return printer
}
