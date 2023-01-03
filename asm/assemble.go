package asm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type Instruction struct {
	matches  func(*string) bool
	callback func(*Context, []string) error
}

type Context struct {
	instructions []Instruction
	assembler    *ir.IrGenerator
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

func shift(args []string, err error) (string, []string, error) {
	if len(args) == 0 {
		return "", nil, err
	}
	return args[0], args[1:], nil
}

func shiftIf(args []string, token string, err error) ([]string, error) {
	if len(args) == 0 || args[0] != token {
		return nil, err
	}
	return args[1:], nil
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

func pushImmediateOrVar(context *Context, destType ir.IrType, token string) error {
	if value, err := ir.ParseNumber[uint64](token); err == nil {
		// Push immediate.
		return context.assembler.PushImmediate(destType, value)
	}

	// Push variable.
	sourceVar, err := context.assembler.LookupVar(token)
	if err != nil {
		return err
	}

	if sourceVar.Type != destType {
		return fmt.Errorf("Variable %q has type %d instead of %d", token, destType, sourceVar.Type)
	}

	return context.assembler.PushVar(token)
}

func assemblePrint2Args(context *Context, typ string, sign ir.Sign, token string) error {
	optype, err := ir.ParseType(typ)
	if err != nil {
		return err
	}

	value, err := ir.ParseNumber[uint64](token)
	if err != nil {
		return err
	}

	return context.assembler.PrintImmediate(optype, sign, value)
}

func assemblePrint(sign ir.Sign) func(*Context, []string) error {
	return func(context *Context, args []string) error {
		switch len(args) {
		case 1:
			return context.assembler.PrintVar(sign, args[0])
		case 2:
			return assemblePrint2Args(context, args[0], sign, args[1])
		default:
			return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
		}
	}
}

func assembleDeclaration(context *Context, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("expected at least 3 arguments; got %q", args)
	}

	if args[1] != ":" {
		return fmt.Errorf("expected ':' as the second argument; got %q", args)
	}
	args = append(args[0:1], args[2:]...)

	id := args[0]
	args = args[1:]

	fmt.Fprintf(os.Stderr, "DEBUG HERE decl %v\n", args)

	vars, err := parseType(strings.Join(args, " "), false /* namedVars */)
	if err != nil {
		return err
	}

	var argTypes []ir.IrType
	var retTypes []ir.IrType
	for _, irvar := range vars {
		if irvar.VarType == ir.ArgVar {
			argTypes = append(argTypes, irvar.Type)
		} else {
			retTypes = append(retTypes, irvar.Type)
		}
	}

	return context.assembler.Declare(id, argTypes, retTypes)
}

func assembleFunc(context *Context, args []string) error {
	if len(args) == 0 || args[len(args)-1] != "{" {
		return fmt.Errorf("expected '{' before end of line of the 'func' instruction; got %q", args)
	}
	args = args[:len(args)-1]

	id, args, err := shift(args, fmt.Errorf("expected identifier after the 'func' token; got %v", args))
	if err != nil {
		return err
	}

	args, err = shiftIf(args, ":", fmt.Errorf("expected token ':' after the function's identifier; got %v", args))
	if err != nil {
		return err
	}

	vars, err := parseType(strings.Join(args, " "), true /* namedVars */)
	if err != nil {
		return err
	}

	if err := context.assembler.Function(id, vars); err != nil {
		return err
	}

	return nil
}

func assembleCall(context *Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected at least 1 argument; got %q", args)
	}

	return context.assembler.Call(args[0], args[1:], nil /* rets */)
}

func assembleAssignCall(context *Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected at least 1 argument; got %q", args)
	}

	var rets []string
	for ; len(args) > 0; args = args[1:] {
		if args[0] == "<-" {
			break
		}

		rets = append(rets, args[0])
	}

	if len(args) < 2 || args[0] != "<-" || args[1] != "call" {
		return fmt.Errorf("expected tokens '<- call' in assignment from call; got %q", args)
	}
	args = args[2:]

	if len(args) == 0 {
		return fmt.Errorf("expected function after '<- call' in assignment from call; got %q", args)
	}
	id := args[0]
	args = args[1:]

	return context.assembler.Call(id, args, rets)
}

func assembleDefineLocal(context *Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expects 2 argument; got %q", args)
	}

	typ, err := ir.ParseType(args[1])
	if err != nil {
		return err
	}

	return context.assembler.DefineLocal(args[0], typ)
}

func assembleIf(context *Context, args []string) error {
	if len(args) == 0 || args[len(args)-1] != "{" {
		return fmt.Errorf("expected '{' before end of line of the 'if' instruction; got %q", args)
	}
	args = args[:len(args)-1]

	then := true
	if len(args) > 0 && args[len(args)-1] == "else" {
		args = args[:len(args)-1]
		then = false
	}

	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	// TODO: Avoid pushing to stack. Instead, pass the variable offset
	// as immediate in the 'if'.
	if err := context.assembler.PushVar(args[0]); err != nil {
		return err
	}

	if then {
		return context.assembler.IfThen()
	}
	return context.assembler.IfElse()
}

// assembleAssign2Args assembles an assign op where the right side is
// either a variable or a literal. The token '<-' should not be passed
// in 'args'.
//
// x <- y
// x <- 123
func assembleAssign2Args(context *Context, args []string) error {
	dest := args[0]
	source := args[1]

	destVar, err := context.assembler.LookupVar(dest)
	if err != nil {
		return err
	}

	if err := pushImmediateOrVar(context, destVar.Type, source); err != nil {
		return err
	}

	return context.assembler.PopVar(dest)
}

// assembleAssign3Args assembles an assign op where the right side is
// a unary op on either a variable or a literal. The token '<-' should
// not be passed in 'args'.
//
// x <- <unaryOp> y
// x <- <unaryOp> 123
func assembleAssign3Args(context *Context, args []string) error {
	dest := args[0]
	op := args[1]
	source := args[2]

	destVar, err := context.assembler.LookupVar(dest)
	if err != nil {
		return err
	}

	if err := pushImmediateOrVar(context, destVar.Type, source); err != nil {
		return err
	}

	switch op {
	case "-":
		if err := context.assembler.Neg(destVar.Type); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Undefined op %q", op)
	}

	return context.assembler.PopVar(dest)
}

// assembleAssign3Args assembles an assign op where the right side is
// a binary op on a pair of variables and/or literals. The token '<-'
// should not be passed in 'args'.
//
// x <- y   <binaryOp> z
// x <- 123 <binaryOp> 456
// x <- y   <binaryOp> 123
// x <- 123 <binaryOp> y
func assembleAssign4Args(context *Context, args []string) error {
	dest := args[0]
	source1 := args[1]
	op := args[2]
	source2 := args[3]

	destVar, err := context.assembler.LookupVar(dest)
	if err != nil {
		return err
	}

	if err := pushImmediateOrVar(context, destVar.Type, source1); err != nil {
		return err
	}

	if err := pushImmediateOrVar(context, destVar.Type, source2); err != nil {
		return err
	}

	switch op {
	case "+":
		if err := context.assembler.Add(destVar.Type); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Undefined op %q", op)
	}

	return context.assembler.PopVar(dest)
}

func assembleAssign(context *Context, args []string) error {
	switch len(args) {
	case 2:
		return assembleAssign2Args(context, args)
	case 3:
		return assembleAssign3Args(context, args)
	case 4:
		return assembleAssign4Args(context, args)
	default:
		return fmt.Errorf("expected 1, 2 or 3 arguments; got %q", args)
	}
}

func assembleFallback(context *Context, args []string) error {
	if len(args) > 1 && args[1] == "<-" {
		return assembleAssign(context, append(args[:1], args[2:]...))
	}

	return fmt.Errorf("Unknown instruction %q", args)
}

func assembleInstruction(context *Context, line string) error {
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

func assembleFile(context *Context, input *os.File) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if err := assembleInstruction(context, scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func AssembleFile(inputFile *os.File) (ir.IrProgram, error) {
	assembler := ir.New()

	context := &Context{
		[]Instruction{
			{prefix("decls {"), noargs(assembler.Decls)},
			{prefix("func "), assembleFunc},

			{contains(" : "), assembleDeclaration},

			{suffix(" i8"), assembleDefineLocal},
			{suffix(" i16"), assembleDefineLocal},
			{suffix(" i32"), assembleDefineLocal},
			{suffix(" i64"), assembleDefineLocal},

			{prefix("call "), assembleCall},
			{contains(" <- call "), assembleAssignCall},

			{prefix("if "), assembleIf},
			{prefix("} else {"), noargs(assembler.Else)},

			{prefix("printU "), assemblePrint(ir.Unsigned)},
			{prefix("printS "), assemblePrint(ir.Signed)},

			{prefix("}"), noargs(assembler.End)},
			{prefix(""), assembleFallback}, // Used for assign (<-) also.
		},
		assembler,
	}

	if err := assembler.Module(); err != nil {
		return ir.IrProgram{}, err
	}

	if err := assembleFile(context, inputFile); err != nil {
		return ir.IrProgram{}, err
	}

	if err := assembler.End(); err != nil {
		return ir.IrProgram{}, err
	}

	return context.assembler.Program(), nil
}
