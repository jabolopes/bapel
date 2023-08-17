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

func noargs(callback func() error) func(*Context, []string) error {
	return func(_ *Context, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("expected no arguments; got %q", args)
		}
		return callback()
	}
}

func trimPrefix(arg *string, token string, err error) error {
	if !strings.HasPrefix(*arg, token) {
		return err
	}
	*arg = strings.TrimPrefix(*arg, token)
	return nil
}

func trimSuffix(arg *string, token string, err error) error {
	if !strings.HasSuffix(*arg, token) {
		return err
	}
	*arg = strings.TrimSuffix(*arg, token)
	return nil
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
	args, err := parser.ShiftIfEnd(args, "{", fmt.Errorf("expected '{' before end of line of the 'func' instruction; got %q", args))
	if err != nil {
		return err
	}

	id, vars, args, err := bplparser.ParseDef(args)
	if err != nil {
		return err
	}

	if err := parser.EOL(args); err != nil {
		return err
	}

	return context.compiler.Function(id, vars)
}

func compileCall(context *Context, args []string) error {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first argument to call; got %v", args))
	if err != nil {
		return err
	}

	argTokens, err := parser.ParseTokens(args)
	if err != nil {
		return err
	}

	return context.compiler.Call(id, argTokens, nil /* rets */)
}

func compileDefineLocal(context *Context, args []string) error {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in variable definition; got %v", args))
	if err != nil {
		return err
	}

	typ, args, err := bplparser.ParseType(args, false /* named */)
	if err != nil {
		return err
	}

	if err := parser.EOL(args); err != nil {
		return err
	}

	return context.compiler.DefineLocal(ir.NewDecl(id, typ))
}

func compileIf(context *Context, args []string) error {
	args, err := parser.ShiftIfEnd(args, "{", fmt.Errorf("expected '{' before end of line of the 'if' instruction; got %q", args))
	if err != nil {
		return err
	}

	then := true
	if len(args) > 0 && args[len(args)-1] == "else" {
		args = args[:len(args)-1]
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

func compileAssign(context *Context, args []string) error {
	var rets []string
	for ; len(args) > 0; args = args[1:] {
		if args[0] == "<-" {
			break
		}

		rets = append(rets, args[0])
	}

	if len(rets) == 0 {
		return fmt.Errorf("expected at least 1 return variable; got %q", args)
	}

	var err error
	args, err = parser.ShiftIf(args, "<-", fmt.Errorf("expected token '<-' as second token in assignment; got %v", args))
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("expected at least 1 argument; got %q", args)
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
		if instruction.matches(&line) {
			return instruction.callback(context, parser.Words(line))
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

func CompileFile(inputFile *os.File, output io.Writer) (ir.IrProgram, error) {
	compiler := ir.NewCompiler(output)

	context := &Context{
		[]Instruction{
			{prefix("imports {"), noargs(compiler.Imports)},
			{prefix("exports {"), noargs(compiler.Exports)},
			{prefix("decls {"), noargs(compiler.Decls)},
			{prefix("func "), compileFunc},

			{contains(" : "), compileDeclaration},

			{suffix(" i8"), compileDefineLocal},
			{suffix(" i16"), compileDefineLocal},
			{suffix(" i32"), compileDefineLocal},
			{suffix(" i64"), compileDefineLocal},

			{prefix("call "), compileCall},
			{contains(" <- "), compileAssign},

			{prefix("if "), compileIf},
			{prefix("} else {"), noargs(compiler.Else)},

			{prefix("printU "), compilePrint(ir.Unsigned)},
			{prefix("printS "), compilePrint(ir.Signed)},

			{prefix("}"), noargs(compiler.End)},
		},
		compiler,
	}

	if err := compiler.Module(); err != nil {
		return ir.IrProgram{}, err
	}

	if err := compileFile(context, inputFile); err != nil {
		return ir.IrProgram{}, err
	}

	if err := compiler.End(); err != nil {
		return ir.IrProgram{}, err
	}

	return context.compiler.Program(), nil
}
