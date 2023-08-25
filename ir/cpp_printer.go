package ir

import (
	"fmt"
	"io"

	"github.com/jabolopes/bapel/parser"
)

type CppPrinter struct {
	output       io.Writer
	bindPosition bool
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
	switch typ.Case {
	case ArrayType:
		fmt.Fprintf(p.out(), "std::array<")
		p.printType(typ.ArrayType.ElemType)
		fmt.Fprintf(p.out(), ", %d>", typ.ArrayType.Size)
	case FunType:
		panic(fmt.Errorf("printType: Unimplemented function type"))
	case IntType:
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
	case IDType:
		fmt.Fprintf(p.out(), "struct %s", toID(typ.IDType))
	default:
		panic(fmt.Errorf("printType: Unhandled case %d", typ.Case))
	}
}

func (p *CppPrinter) PrintTerm(term IrTerm) {
	switch term.Case {
	case AssignTerm:
		p.bindPosition = true
		p.PrintTerm(term.Assign.Ret)
		p.bindPosition = false
		p.printf(" = ")
		p.PrintTerm(term.Assign.Arg)

	case CallTerm:
		p.printCall(term.Call.ID, term.Call.Args)

	case IfTerm:

	case TokenTerm:
		p.printToken(*term.Token)

	case TupleTerm:
		if p.bindPosition {
			p.printf("std::tie(")
		}

		if len(term.Tuple) > 0 {
			p.PrintTerm(term.Tuple[0])
			for _, term := range term.Tuple[1:] {
				p.printf(", ")
				p.PrintTerm(term)
			}
		}

		if p.bindPosition {
			p.printf(")")
		}
	}
}

func NewCppPrinter(output io.Writer) *CppPrinter {
	return &CppPrinter{
		output,
		false, /* bindPosition */
	}
}
