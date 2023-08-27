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

func compileAny(context *Context, args []string) error {
	if section, err := context.parser.ParseSection(); err == nil {
		return context.compiler.Section(section)
	}

	if id, argTuple, retTuple, err := context.parser.ParseFunc(); err == nil {
		return context.compiler.Function(id, argTuple, retTuple)
	}

	if decl, err := context.parser.ParseLet(); err == nil {
		return context.compiler.DefineLocal(decl)
	}

	if decl, _, err := context.parser.ParseDecl(args, false /* named */); err == nil {
		return context.compiler.Declare(decl)
	}

	if ifTerm, err := context.parser.ParseIf(); err == nil {
		return context.compiler.If(ifTerm)
	}

	if err := context.parser.ParseElse(); err == nil {
		return context.compiler.Else()
	}

	if err := context.parser.ParseEnd(); err == nil {
		return context.compiler.End()
	}

	if id, typ, _, err := context.parser.ParseStruct(args); err == nil {
		return context.compiler.Struct(id, typ)
	}

	if id, _, err := context.parser.ParseEntity(args); err == nil {
		return context.compiler.Entity(id)
	}

	// PrintU/S.
	if args, err := parser.ShiftToken(args, "printU"); err == nil {
		return compilePrint(context, ir.Unsigned, args)
	}

	if args, err := parser.ShiftToken(args, "printS"); err == nil {
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
		if err := compileAny(context, context.parser.Words()); err != nil {
			return fmt.Errorf("in line\n  %s\n%v\n", context.parser.Line(), err)
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
