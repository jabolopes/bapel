package asm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/vm"
	"golang.org/x/exp/constraints"
)

type Instruction struct {
	token    string
	callback func(*Context, []string) error
}

type Context struct {
	instructions []Instruction
	assembler    *ir.OpAssembler
}

func noargs(callback func() error) func(*Context, []string) error {
	return func(_ *Context, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("expected no arguments; got %q", args)
		}
		return callback()
	}
}

func family(callback func(ir.OpType) error) func(*Context, []string) error {
	return func(_ *Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument; got %q", args)
		}

		optype, err := ir.ParseOpType(args[0])
		if err != nil {
			return err
		}

		return callback(optype)
	}
}

func assemblePush(context *Context, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
	}

	if len(args) == 1 {
		// Push local.
		return context.assembler.PushLocal(args[0])
	}

	// Push immediate.
	optype, err := ir.ParseOpType(args[0])
	if err != nil {
		return err
	}

	value, err := ir.ParseNumber[uint64](args[1])
	if err != nil {
		return err
	}

	return context.assembler.PushImmediate(optype, value)
}

func assemblePopLocal(context *Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	return context.assembler.PopLocal(args[0])
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

func assembleDefineVar[T constraints.Integer]() func(*Context, []string) error {
	var value T
	size := uint16(unsafe.Sizeof(value))
	return func(context *Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expects 1 argument; got %q", args)
		}

		return context.assembler.DefineVar(args[0], size)
	}
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
			{"pop", assemblePopLocal},

			{"add", family(assembler.Add)},
			{"print", family(assembler.Print)},

			{"func", assembleFunc},

			{"args {", noargs(assembler.Args)},
			{"rets {", noargs(assembler.Rets)},
			{"locals {", noargs(assembler.Locals)},
			{"i8", assembleDefineVar[byte]()},
			{"i16", assembleDefineVar[uint16]()},
			{"i32", assembleDefineVar[uint32]()},
			{"i64", assembleDefineVar[uint64]()},

			{"if else {", noargs(assembler.IfElse)},
			{"if {", noargs(assembler.IfThen)},
			{"} else {", noargs(assembler.Else)},
			{"}", noargs(assembler.End)},
		},
		assembler,
	}

	if err := assembleFile(context, file); err != nil {
		return vm.OpProgram{}, err
	}

	return context.assembler.Program(), nil
}
