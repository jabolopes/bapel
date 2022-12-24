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

func family(callback func(ir.IrType) error) func(*Context, []string) error {
	return func(_ *Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument; got %q", args)
		}

		optype, err := ir.ParseType(args[0])
		if err != nil {
			return err
		}

		return callback(optype)
	}
}

func pushImmediateOrVar(context *Context, destType ir.IrType, token string) error {
	if value, err := ir.ParseNumber[uint64](token); err == nil {
		// Push immediate.
		return context.assembler.PushImmediate(destType, value)
	}

	// Push variable.
	sourceVar, ok := context.assembler.LookupVar(token)
	if !ok {
		return fmt.Errorf("Undefined variable %q", token)
	}

	if sourceVar.Type != destType {
		return fmt.Errorf("Variable %q has type %d instead of %d", token, destType, sourceVar.Type)
	}

	return context.assembler.PushVar(token)
}

func assemblePush(context *Context, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
	}

	if len(args) == 1 {
		// Push local.
		return context.assembler.PushVar(args[0])
	}

	// Push immediate.
	optype, err := ir.ParseType(args[0])
	if err != nil {
		return err
	}

	value, err := ir.ParseNumber[uint64](args[1])
	if err != nil {
		return err
	}

	return context.assembler.PushImmediate(optype, value)
}

func assemblePopVar(context *Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	return context.assembler.PopVar(args[0])
}

func assembleFunc(context *Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments; got %q", args)
	}

	if args[1] != "{" {
		return fmt.Errorf("expected '{' after the function's identifier; got %q", args)
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

// assembleAssign2Args assembles an assign op where the right side is
// either a variable or a literal. The token '<-' should not be passed
// in 'args'.
//
// x <- y
// x <- 123
func assembleAssign2Args(context *Context, args []string) error {
	dest := args[0]
	source := args[1]

	destVar, ok := context.assembler.LookupVar(dest)
	if !ok {
		return fmt.Errorf("Undefined variable %q", dest)
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

	destVar, ok := context.assembler.LookupVar(dest)
	if !ok {
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

	destVar, ok := context.assembler.LookupVar(dest)
	if !ok {
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
			{"push", assemblePush},
			{"pop", assemblePopVar},

			{"add", family(assembler.Add)},
			{"print", family(assembler.Print)},

			{"func", assembleFunc},

			{"args {", noargs(assembler.Args)},
			{"rets {", noargs(assembler.Rets)},
			{"locals {", noargs(assembler.Locals)},
			{"i8", assembleDefineVar(ir.I8)},
			{"i16", assembleDefineVar(ir.I16)},
			{"i32", assembleDefineVar(ir.I32)},
			{"i64", assembleDefineVar(ir.I64)},

			{"if else {", noargs(assembler.IfElse)},
			{"if {", noargs(assembler.IfThen)},
			{"} else {", noargs(assembler.Else)},
			{"}", noargs(assembler.End)},
			{"", assembleFallback},
		},
		assembler,
	}

	if err := assembleFile(context, file); err != nil {
		return vm.OpProgram{}, err
	}

	return context.assembler.Program(), nil
}
