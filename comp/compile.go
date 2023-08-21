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

type Instruction struct {
	matches  func(*string) bool
	callback func(*Context, []string) error
}

type Context struct {
	instructions []Instruction
	compiler     *ir.Compiler
}

func prefix(token string) func(*string) bool {
	return func(line *string) bool {
		if !strings.HasPrefix(*line, token) {
			return false
		}
		*line = strings.TrimPrefix(*line, token)
		*line = strings.TrimPrefix(*line, " ")
		return true
	}
}

func suffix(token string) func(*string) bool {
	return func(line *string) bool {
		return strings.HasSuffix(*line, token)
	}
}

func contains(token string) func(*string) bool {
	return func(line *string) bool {
		return strings.Contains(*line, token)
	}
}

func always() func(*string) bool {
	return func(line *string) bool {
		return true
	}
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

func compilePrint(sign ir.Sign) func(*Context, []string) error {
	return func(context *Context, args []string) error {
		switch len(args) {
		case 1:
			return context.compiler.PrintVar(sign, args[0])
		case 2:
			return compilePrintImmediate(context, args[0], sign, args[1])
		default:
			return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
		}
	}
}

func compileDeclaration(context *Context, args []string) error {
	decl, args, err := bplparser.ParseDecl(args, false /* named */)
	if err != nil {
		return err
	}

	if err := parser.EOL(args); err != nil {
		return err
	}

	return context.compiler.Declare(decl)
}

func compileFunc(context *Context, args []string) error {
	args, err := parser.ShiftTokenEnd(args, "{")
	if err != nil {
		return err
	}

	id, vars, args, err := bplparser.ParseFunc(args)
	if err != nil {
		return err
	}

	if err := parser.EOL(args); err != nil {
		return err
	}

	return context.compiler.Function(id, vars)
}

func compileAny(context *Context, args []string) error {
	for _, section := range []string{"imports", "decls", "exports"} {
		if args, err := parser.ShiftTokens(args, []string{section, "{"}); err == nil {
			if err := parser.EOL(args); err != nil {
				return err
			}

			switch section {
			case "imports":
				return context.compiler.Imports()
			case "decls":
				return context.compiler.Decls()
			case "exports":
				return context.compiler.Exports()
			}
		}
	}

	if id, vars, _, err := bplparser.ParseFunc(args); err == nil {
		return context.compiler.Function(id, vars)
	}

	if decl, _, err := bplparser.ParseLet(args); err == nil {
		return context.compiler.DefineLocal(decl)
	}

	if len(args) >= 2 && args[1] == ":" {
		decl, args, err := bplparser.ParseDecl(args, false /* named */)
		if err != nil {
			return err
		}

		if err := parser.EOL(args); err != nil {
			return err
		}

		return context.compiler.Declare(decl)
	}

	if args, err := parser.ShiftToken(args, "if"); err == nil {
		args, err := parser.ShiftTokenEnd(args, "{")
		if err != nil {
			return err
		}

		then := true
		if args, err = parser.ShiftTokenEnd(args, "else"); err == nil {
			then = false
		}

		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument; got %q", args)
		}

		if then {
			return context.compiler.IfThen(args[0])
		}
		return context.compiler.IfElse(args[0])
	}

	if args, err := parser.ShiftTokens(args, []string{"}", "else", "{"}); err == nil {
		if err := parser.EOL(args); err != nil {
			return err
		}

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

	if args[0] == "entity" {
		id, args, err := bplparser.ParseEntity(args)
		if err != nil {
			return err
		}

		if err := parser.EOL(args); err != nil {
			return err
		}

		return context.compiler.Entity(id)
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

func compileInstruction(context *Context, line string) error {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil
	}

	for _, instruction := range context.instructions {
		matchLine := line
		if instruction.matches(&matchLine) {
			err := instruction.callback(context, parser.Words(matchLine))
			if err != nil {
				err = fmt.Errorf("in line\n  %s\n%v\n", line, err)
			}
			return err
		}
	}

	return fmt.Errorf("Unknown instruction line %q", line)
}

func compileFile(context *Context, input *os.File) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if err := compileInstruction(context, scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func CompileFile(inputFile *os.File, output io.Writer) error {
	compiler := ir.NewCompiler(output)

	context := &Context{
		[]Instruction{
			{prefix("printU "), compilePrint(ir.Unsigned)},
			{prefix("printS "), compilePrint(ir.Signed)},

			{always(), compileAny},
		},
		compiler,
	}

	if err := compiler.Module(); err != nil {
		return err
	}

	if err := compileFile(context, inputFile); err != nil {
		return err
	}

	return compiler.End()
}
