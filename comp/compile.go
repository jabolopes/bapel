package comp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

type Context struct {
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
	if section, _, err := bplparser.ParseSection(args); err == nil {
		return context.compiler.Section(section)
	}

	if id, argTuple, retTuple, _, err := bplparser.ParseFunc(args); err == nil {
		return context.compiler.Function(id, argTuple, retTuple)
	}

	if decl, _, err := bplparser.ParseLet(args); err == nil {
		return context.compiler.DefineLocal(decl)
	}

	if decl, _, err := bplparser.ParseDecl(args, false /* named */); err == nil {
		return context.compiler.Declare(decl)
	}

	if then, argTokens, _, err := bplparser.ParseIf(args); err == nil {
		// TODO: Validate that argTokens[0] is an ID.
		return context.compiler.If(then, argTokens[0].Text, argTokens[1:])
	}

	if _, err := bplparser.ParseElse(args); err == nil {
		return context.compiler.Else()
	}

	if args, err := parser.ShiftToken(args, "}"); err == nil {
		if err := parser.EOL(args); err != nil {
			return err
		}

		return context.compiler.End()
	}

	if id, typ, _, err := bplparser.ParseStruct(args); err == nil {
		return context.compiler.Struct(id, typ)
	}

	if id, _, err := bplparser.ParseEntity(args); err == nil {
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
	args, rets, err := bplparser.ParseCallAssign(args)
	if err != nil {
		return err
	}

	argTokens, err := parser.ParseTokens(args)
	if err != nil {
		return err
	}

	return context.compiler.Assign(argTokens, rets)
}

func compileFile(context *Context, input *os.File) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := compileAny(context, parser.Words(line)); err != nil {
			return fmt.Errorf("in line\n  %s\n%v\n", line, err)
		}
	}

	return scanner.Err()
}

func CompileFile(inputFile *os.File, output io.Writer) error {
	compiler := ir.NewCompiler(output)

	context := &Context{compiler}

	if err := compiler.Module(); err != nil {
		return err
	}

	if err := compileFile(context, inputFile); err != nil {
		return err
	}

	return compiler.End()
}
