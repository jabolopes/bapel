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
	source, err := context.parser.ParseAny()
	if err != nil {
		return err
	}

	switch source.Case {
	case bplparser.SectionSource:
		return context.compiler.Section(source.Section)
	case bplparser.DeclSource:
		return context.compiler.Declare(*source.Decl)
	case bplparser.EntitySource:
		return context.compiler.Entity(source.Entity)
	case bplparser.FunctionSource:
		return context.compiler.Function(source.Function.ID, source.Function.Args, source.Function.Rets)
	case bplparser.TermSource:
		return context.compiler.Term(*source.Term)
	case bplparser.ElseSource:
		return context.compiler.Else()
	case bplparser.EndSource:
		return context.compiler.End()
	case bplparser.PrintSource:
		return compilePrint(context, source.Print.Sign, source.Print.Args)
	default:
		return fmt.Errorf("unhandled source case %d", source.Case)
	}
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
