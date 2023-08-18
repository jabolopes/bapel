package asm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
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

		typ, err := ir.ParseIntType(typStr)
		if err != nil {
			return nil, err
		}

		vars = append(vars, ir.IrVar{Id: id, VarType: ir.ArgVar, Type: ir.NewIntType(typ)})
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

		typ, err := ir.ParseIntType(typStr)
		if err != nil {
			return nil, err
		}

		vars = append(vars, ir.IrVar{Id: id, VarType: ir.RetVar, Type: ir.NewIntType(typ)})
	}

	return vars, nil
}

func pushImmediateOrVar(context *Context, destType ir.IrIntType, token string) error {
	if value, err := parser.ParseNumber[uint64](token); err == nil {
		// Push immediate.
		return context.assembler.PushImmediate(destType, value)
	}

	// Push variable.
	sourceVar, err := context.assembler.LookupVar(token)
	if err != nil {
		return err
	}

	// TODO: Check that type is actual IntType.
	if sourceVar.Type.IntType != destType {
		return fmt.Errorf("Variable %q has type %s instead of %s", token, destType, sourceVar.Type)
	}

	return context.assembler.PushVar(token)
}

func assemblePrintImmediate(context *Context, typ string, sign ir.Sign, token string) error {
	optype, err := ir.ParseIntType(typ)
	if err != nil {
		return err
	}

	value, err := parser.ParseNumber[uint64](token)
	if err != nil {
		return err
	}

	return context.assembler.PrintImmediate(optype, sign, value)
}

func assembleIOWait(context *Context, rets, args []string) error {
	opID, args, err := parser.Shift(args, fmt.Errorf("expected exactly 1 argument in call to 'io.wait'; got %v", args))
	if err != nil {
		return err
	}

	if len(args) > 0 {
		return fmt.Errorf("too many arguments given to 'io.wait'; got %v", args)
	}

	errID, rets, err := parser.Shift(rets, fmt.Errorf("expected exactly 2 return values in call to 'io.wait'; got %v", args))
	if err != nil {
		return err
	}

	valueID, rets, err := parser.Shift(rets, fmt.Errorf("expected exactly 2 return values in call to 'io.wait'; got %v", args))
	if err != nil {
		return err
	}

	if len(rets) > 0 {
		return fmt.Errorf("too many return values given to 'io.wait'; got %v", args)
	}

	return context.assembler.IOWait(opID, errID, valueID)
}

func assembleIODo(context *Context, rets, args []string) error {
	funID, args, err := parser.Shift(args, fmt.Errorf("expected exactly 1 argument in call to 'io.do'; got %v", args))
	if err != nil {
		return err
	}

	if len(args) > 0 {
		return fmt.Errorf("too many arguments given to 'io.do'; got %v", args)
	}

	retID, rets, err := parser.Shift(rets, fmt.Errorf("expected exactly 1 argument in call to 'io.do'; got %v", rets))
	if err != nil {
		return err
	}

	if len(rets) > 0 {
		return fmt.Errorf("too many return values given to 'io.do'; got %v", rets)
	}

	return context.assembler.IODo(funID, retID)
}

func assemblePrint(sign ir.Sign) func(*Context, []string) error {
	return func(context *Context, args []string) error {
		switch len(args) {
		case 1:
			return context.assembler.PrintVar(sign, args[0])
		case 2:
			return assemblePrintImmediate(context, args[0], sign, args[1])
		default:
			return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
		}
	}
}

func assembleDeclaration(context *Context, args []string) error {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in declaration; got %v", args))
	if err != nil {
		return err
	}

	args, err = parser.ShiftIf(args, ":", fmt.Errorf("expected token ':' after the declaration's identifier; got %v", args))
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("expected type in declaration; got %v", args)
	}

	vars, err := parseType(strings.Join(args, " "), false /* namedVars */)
	if err != nil {
		return err
	}

	var argTypes []ir.IrIntType
	var retTypes []ir.IrIntType
	for _, irvar := range vars {
		if irvar.VarType == ir.ArgVar {
			// TODO: Check that the type is actually IntType.
			argTypes = append(argTypes, irvar.Type.IntType)
		} else {
			// TODO: Check that the type is actually IntType.
			retTypes = append(retTypes, irvar.Type.IntType)
		}
	}

	return context.assembler.Declare(id, argTypes, retTypes)
}

func assembleFunc(context *Context, args []string) error {
	args, err := parser.ShiftIfEnd(args, "{", fmt.Errorf("expected '{' before end of line of the 'func' instruction; got %q", args))
	if err != nil {
		return err
	}

	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier after the 'func' token; got %v", args))
	if err != nil {
		return err
	}

	args, err = parser.ShiftIf(args, ":", fmt.Errorf("expected token ':' after the function's identifier; got %v", args))
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

	return context.assembler.Function(id, vars)
}

func assembleCall(context *Context, args []string) error {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first argument to call; got %v", args))
	if err != nil {
		return err
	}

	return context.assembler.Call(id, args, nil /* rets */)
}

func assembleAssignCall(context *Context, rets, args []string) error {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first argument to call; got %v", args))
	if err != nil {
		return err
	}

	return context.assembler.Call(id, args, rets)
}

func assembleAssignSyscall(context *Context, rets, args []string) error {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first argument to call; got %v", args))
	if err != nil {
		return err
	}

	return context.assembler.Syscall(id, args, rets)
}

func assembleDefineLocal(context *Context, args []string) error {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in variable definition; got %v", args))
	if err != nil {
		return err
	}

	typStr, args, err := parser.Shift(args, fmt.Errorf("expected type as second token in variable definition; got %v", args))
	if err != nil {
		return err
	}

	if len(args) > 0 {
		return fmt.Errorf("too many token in variable definition; got %v", args)
	}

	typ, err := ir.ParseIntType(typStr)
	if err != nil {
		return err
	}

	return context.assembler.DefineLocal(id, typ)
}

func assembleIf(context *Context, args []string) error {
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
func assembleAssign2Args(context *Context, dest, source string) error {
	destVar, err := context.assembler.LookupVar(dest)
	if err != nil {
		return err
	}

	// TODO: Check that type is actual IntType.
	if err := pushImmediateOrVar(context, destVar.Type.IntType, source); err != nil {
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
func assembleAssign3Args(context *Context, dest, op, source string) error {
	destVar, err := context.assembler.LookupVar(dest)
	if err != nil {
		return err
	}

	// TODO: Check that type is actual IntType.
	if err := pushImmediateOrVar(context, destVar.Type.IntType, source); err != nil {
		return err
	}

	switch op {
	case "-":
		// TODO: Check that type is actual IntType.
		if err := context.assembler.Neg(destVar.Type.IntType); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Undefined op %q", op)
	}

	return context.assembler.PopVar(dest)
}

// assembleAssign4Args assembles an assign op where the right side is
// a binary op on a pair of variables and/or literals. The token '<-'
// should not be passed in 'args'.
//
// x <- y   <binaryOp> z
// x <- 123 <binaryOp> 456
// x <- y   <binaryOp> 123
// x <- 123 <binaryOp> y
func assembleAssign4Args(context *Context, dest, source1, op, source2 string) error {
	destVar, err := context.assembler.LookupVar(dest)
	if err != nil {
		return err
	}

	// TODO: Check that type is actual IntType.
	if err := pushImmediateOrVar(context, destVar.Type.IntType, source1); err != nil {
		return err
	}

	// TODO: Check that type is actual IntType.
	if err := pushImmediateOrVar(context, destVar.Type.IntType, source2); err != nil {
		return err
	}

	switch op {
	case "+":
		// TODO: Check that type is actual IntType.
		if err := context.assembler.Add(destVar.Type.IntType); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Undefined op %q", op)
	}

	return context.assembler.PopVar(dest)
}

func assembleAssign(context *Context, args []string) error {
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

	if len(args) > 0 {
		switch args[0] {
		case "call":
			return assembleAssignCall(context, rets, args[1:])
		case "syscall":
			return assembleAssignSyscall(context, rets, args[1:])
		case "io.wait":
			return assembleIOWait(context, rets, args[1:])
		case "io.do":
			return assembleIODo(context, rets, args[1:])
		}
	}

	if len(rets) != 1 {
		return fmt.Errorf("expected at most 1 return variable; got %q", args)
	}

	switch len(args) {
	case 1:
		return assembleAssign2Args(context, rets[0], args[0])
	case 2:
		return assembleAssign3Args(context, rets[0], args[0], args[1])
	case 3:
		return assembleAssign4Args(context, rets[0], args[0], args[1], args[2])
	default:
		return fmt.Errorf("expected 1, 2 or 3 arguments; got %q", args)
	}
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
			{contains(" <- "), assembleAssign},

			{prefix("if "), assembleIf},
			{prefix("} else {"), noargs(assembler.Else)},

			{prefix("printU "), assemblePrint(ir.Unsigned)},
			{prefix("printS "), assemblePrint(ir.Signed)},

			{prefix("}"), noargs(assembler.End)},
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
