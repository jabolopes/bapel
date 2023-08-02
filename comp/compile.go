package comp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/shift"
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

func parseType(token string, namedVars bool) ([]ir.IrVar, error) {
	splits := strings.SplitN(token, " -> ", 2)
	if len(splits) != 2 {
		return nil, fmt.Errorf("invalid type; expected '(arg1 type1, ...) -> (ret1 type1, ...)'; got %q", token)
	}

	arg := splits[0]
	ret := splits[1]

	if err := trimPrefix(&arg, "(", fmt.Errorf("expected argument list in type; got %v", token)); err != nil {
		return nil, err
	}

	if err := trimSuffix(&arg, ")", fmt.Errorf("expected argument list in type; got %v", token)); err != nil {
		return nil, err
	}

	if err := trimPrefix(&ret, "(", fmt.Errorf("expected return value list in type; got %v", token)); err != nil {
		return nil, err
	}

	if err := trimSuffix(&ret, ")", fmt.Errorf("expected return value list in type; got %v", token)); err != nil {
		return nil, err
	}

	var args []string
	if len(arg) > 0 {
		args = strings.Split(arg, ", ")
	}

	var rets []string
	if len(ret) > 0 {
		rets = strings.Split(ret, ", ")
	}

	var vars []ir.IrVar
	for _, arg := range args {
		var id string
		var typStr string
		if namedVars {
			splits := strings.SplitN(arg, " ", 2)
			if len(splits) != 2 {
				return nil, fmt.Errorf("expected return value list in type; got %v", arg)
			}
			id = splits[0]
			typStr = splits[1]
		} else {
			typStr = arg
		}

		typ, err := ir.ParseType(typStr)
		if err != nil {
			return nil, err
		}

		vars = append(vars, ir.IrVar{Id: id, VarType: ir.ArgVar, Type: typ})
	}

	for _, ret := range rets {
		var id string
		var typStr string
		if namedVars {
			splits := strings.SplitN(ret, " ", 2)
			if len(splits) != 2 {
				return nil, fmt.Errorf("expected return value list in type; got %v", ret)
			}
			id = splits[0]
			typStr = splits[1]
		} else {
			typStr = ret
		}

		typ, err := ir.ParseType(typStr)
		if err != nil {
			return nil, err
		}

		vars = append(vars, ir.IrVar{Id: id, VarType: ir.RetVar, Type: typ})
	}

	return vars, nil
}

func compilePrintImmediate(context *Context, typ string, sign ir.Sign, token string) error {
	optype, err := ir.ParseType(typ)
	if err != nil {
		return err
	}

	value, err := ir.ParseNumber[uint64](token)
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
	decl, err := ir.ParseDecl(args)
	if err != nil {
		return err
	}

	return context.compiler.Declare(decl)
}

func compileFunc(context *Context, args []string) error {
	args, err := shift.ShiftIfEnd(args, "{", fmt.Errorf("expected '{' before end of line of the 'func' instruction; got %q", args))
	if err != nil {
		return err
	}

	id, args, err := shift.Shift(args, fmt.Errorf("expected identifier after the 'func' token; got %v", args))
	if err != nil {
		return err
	}

	args, err = shift.ShiftIf(args, ":", fmt.Errorf("expected token ':' after the function's identifier; got %v", args))
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("expected type in function definition; got %v", args)
	}

	vars, err := parseType(strings.Join(args, " "), true /* namedVars */)
	if err != nil {
		return err
	}

	return context.compiler.Function(id, vars)
}

func compileCall(context *Context, args []string) error {
	id, args, err := shift.Shift(args, fmt.Errorf("expected identifier as first argument to call; got %v", args))
	if err != nil {
		return err
	}

	return context.compiler.Call(id, args, nil /* rets */)
}

func compileDefineLocal(context *Context, args []string) error {
	id, args, err := shift.Shift(args, fmt.Errorf("expected identifier as first token in variable definition; got %v", args))
	if err != nil {
		return err
	}

	typStr, args, err := shift.Shift(args, fmt.Errorf("expected type as second token in variable definition; got %v", args))
	if err != nil {
		return err
	}

	if len(args) > 0 {
		return fmt.Errorf("too many tokens in variable definition; got %v", args)
	}

	typ, err := ir.ParseType(typStr)
	if err != nil {
		return err
	}

	return context.compiler.DefineLocal(id, typ)
}

func compileIf(context *Context, args []string) error {
	args, err := shift.ShiftIfEnd(args, "{", fmt.Errorf("expected '{' before end of line of the 'if' instruction; got %q", args))
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
	args, err = shift.ShiftIf(args, "<-", fmt.Errorf("expected token '<-' as second token in assignment; got %v", args))
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("expected at least 1 argument; got %q", args)
	}

	return context.compiler.Assign(args, rets)
}

func compileInstruction(context *Context, line string) error {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil
	}

	for _, instruction := range context.instructions {
		if instruction.matches(&line) {
			var args []string
			if line != "" {
				args = strings.Split(line, " ")
			}
			return instruction.callback(context, args)
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
