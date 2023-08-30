package comp

import (
	"fmt"
	"io"
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

type Context struct {
	parser   *bplparser.Parser
	compiler *ir.Compiler
}

func compilePrintImmediate(context *Context, typ string, sign ir.Sign, token string) error {
	optype, err := ir.ParseIntType(typ)
	if err != nil {
		return err
	}

	value, err := parser.ParseNumber[uint64](token)
	if err != nil {
		return err
	}

	return context.compiler.PrintImmediate(optype, sign, value)
}

func compilePrint(context *Context, sign ir.Sign, args []string) error {
	switch len(args) {
	case 1:
		return context.compiler.PrintVar(sign, args[0])
	case 2:
		return compilePrintImmediate(context, args[0], sign, args[1])
	default:
		return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
	}
}

func compileAny(context *Context) error {
	if section, err := context.parser.ParseSection(); err == nil {
		return context.compiler.Section(section)
	}

	if context.parser.PeekToken("func") {
		id, argTuple, retTuple, err := context.parser.ParseFunc()
		if err != nil {
			return err
		}

		return context.compiler.Function(id, argTuple, retTuple)
	}

	if context.parser.PeekToken("struct") {
		decl, err := context.parser.ParseStruct()
		if err != nil {
			return err
		}

		return context.compiler.Define(decl)
	}

	if context.parser.PeekToken("let") {
		decl, err := context.parser.ParseLet()
		if err != nil {
			return err
		}

		return context.compiler.Define(decl)
	}

	if context.parser.PeekToken("if") {
		ifTerm, err := context.parser.ParseIf()
		if err != nil {
			return err
		}

		return context.compiler.If(ifTerm)
	}

	if context.parser.PeekToken("}") {
		if err := context.parser.ParseElse(); err == nil {
			return context.compiler.Else()
		}

		if err := context.parser.ParseEnd(); err != nil {
			return err
		}

		return context.compiler.End()
	}

	if context.parser.PeekToken("entity") {
		id, err := context.parser.ParseEntity()
		if err != nil {
			return err
		}

		return context.compiler.Entity(id)
	}

	if decl, err := context.parser.ParseDecl(false /* named */); err == nil {
		return context.compiler.Declare(decl)
	}

	// PrintU/S.
	if context.parser.PeekToken("printU") {
		args, err := parser.ShiftToken(context.parser.Words(), "printU")
		if err != nil {
			return err
		}

		return compilePrint(context, ir.Unsigned, args)
	}

	if context.parser.PeekToken("printS") {
		args, err := parser.ShiftToken(context.parser.Words(), "printS")
		if err != nil {
			return err
		}

		return compilePrint(context, ir.Signed, args)
	}

	// Parse call / assignment.
	callAssignTerm, err := context.parser.ParseCallAssign()
	if err != nil {
		return err
	}

	return context.compiler.Statement(ir.NewStatementTerm(callAssignTerm))
}

func compileFile(context *Context, input *os.File) error {
	if err := context.compiler.Module(); err != nil {
		return err
	}

	context.parser.Open(input)
	for context.parser.Scan() {
		if err := compileAny(context); err != nil {
			return fmt.Errorf("in line\n  %s\n%v", context.parser.Line(), err)
		}
	}

	return context.parser.ScanErr()
}

func CompileFile(inputFile *os.File, output io.Writer) error {
	compiler := ir.NewCompiler(output)
	context := &Context{bplparser.NewParser(compiler), compiler}
	if err := compileFile(context, inputFile); err != nil {
		return err
	}

	return compiler.End()
}
