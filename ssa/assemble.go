package ssa

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/vm"
)

type Instruction struct {
	token    string
	callback func(*Context, []string) error
}

type Context struct {
	instructions []Instruction
	assembler    *ir.IrGenerator
}

func noargs(callback func() error) func(*Context, []string) error {
	return func(_ *Context, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("expected no arguments; got %q", args)
		}
		return callback()
	}
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

func assemblePrint2Args(context *Context, typ, token string) error {
	optype, err := ir.ParseType(typ)
	if err != nil {
		return err
	}

	value, err := ir.ParseNumber[uint64](token)
	if err != nil {
		return err
	}

	return context.assembler.PrintImmediate(optype, value)
}

func assemblePrint(context *Context, args []string) error {
	switch len(args) {
	case 1:
		return context.assembler.PrintVar(args[0])
	case 2:
		return assemblePrint2Args(context, args[0], args[1])
	default:
		return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
	}
}

func assembleFunc(context *Context, args []string) error {
	if len(args) == 0 || args[len(args)-1] != "{" {
		return fmt.Errorf("expected '{' before end of line of the 'func' instruction; got %q", args)
	}
	args = args[:len(args)-1]

	if len(args) != 1 {
		return fmt.Errorf("expected 1 arguments; got %q", args)
	}

	return context.assembler.Function(args[0])
}

func assembleDefineVar(typ ir.IrType) func(*Context, []string) error {
	return func(context *Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expects 1 argument; got %q", args)
		}

		return context.assembler.DefineVar(args[0], typ)
	}
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
		return fmt.Errorf("Undefined variable %q", dest)
	}

	if err := pushImmediateOrVar(context, destVar.Type, source); err != nil {
		return err
	}

	switch op {
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
		return fmt.Errorf("Undefined variable %q", dest)
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
		if strings.HasPrefix(line, instruction.token) {
			line = strings.TrimPrefix(line, instruction.token)
			line = strings.TrimPrefix(line, " ")
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

func AssembleFile(file *os.File) (vm.OpProgram, error) {
	assembler := ir.New()

	context := &Context{
		[]Instruction{
			{"print", assemblePrint},

			{"func", assembleFunc},

			{"args {", noargs(assembler.Args)},
			{"rets {", noargs(assembler.Rets)},
			{"locals {", noargs(assembler.Locals)},
			{"i8", assembleDefineVar(ir.I8)},
			{"i16", assembleDefineVar(ir.I16)},
			{"i32", assembleDefineVar(ir.I32)},
			{"i64", assembleDefineVar(ir.I64)},

			{"if", assembleIf},
			{"} else {", noargs(assembler.Else)},
			{"}", noargs(assembler.End)},
			{"", assembleFallback}, // Used for assign (<-) also.
		},
		assembler,
	}

	if err := assembleFile(context, file); err != nil {
		return vm.OpProgram{}, err
	}

	return context.assembler.Program(), nil
}
