package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/jabolopes/bapel/asm"
	"github.com/jabolopes/bapel/vm"
	"golang.org/x/exp/constraints"
)

type Instruction struct {
	token    string
	callback func(*Machine, []string) error
}

type Machine struct {
	instructions []Instruction
	assembler    *asm.OpAssembler
}

func noargs(callback func() error) func(*Machine, []string) error {
	return func(_ *Machine, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("expected no arguments; got %q", args)
		}
		return callback()
	}
}

func family(callback func(asm.OpType) error) func(*Machine, []string) error {
	return func(_ *Machine, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument; got %q", args)
		}

		optype, err := asm.ParseOpType(args[0])
		if err != nil {
			return err
		}

		return callback(optype)
	}
}

func assemblePush(machine *Machine, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return fmt.Errorf("expected 1 or 2 arguments; got %q", args)
	}

	if len(args) == 1 {
		// Push local.
		return machine.assembler.PushLocal(args[0])
	}

	// Push immediate.
	optype, err := asm.ParseOpType(args[0])
	if err != nil {
		return err
	}

	value, err := asm.ParseNumber[uint64](args[1])
	if err != nil {
		return err
	}

	return machine.assembler.PushImmediate(optype, value)
}

func assemblePopLocal(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	return machine.assembler.PopLocal(args[0])
}

func assembleFunc(machine *Machine, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments; got %q", args)
	}

	if args[1] != "{" {
		return fmt.Errorf("expected '{' after the function's identifier; got %q", args)
	}

	return machine.assembler.Function(args[0])
}

func assembleDefineVar[T constraints.Integer]() func(*Machine, []string) error {
	var value T
	size := uint16(unsafe.Sizeof(value))
	return func(machine *Machine, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expects 1 argument; got %q", args)
		}

		return machine.assembler.DefineVar(args[0], size)
	}
}

func assembleOp(machine *Machine, line string) error {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil
	}

	for _, instruction := range machine.instructions {
		if strings.HasPrefix(line, instruction.token) {
			line = strings.TrimPrefix(line, instruction.token)
			line = strings.TrimPrefix(line, " ")
			var args []string
			if line != "" {
				args = strings.Split(line, " ")
			}
			return instruction.callback(machine, args)
		}
	}

	return fmt.Errorf("Unknown instruction line %q", line)
}

func assembleFile(machine *Machine, input *os.File) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if err := assembleOp(machine, scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func run() error {
	assembler := asm.New()
	machine := &Machine{
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
	if err := assembleFile(machine, os.Stdin); err != nil {
		return err
	}

	vmachine := vm.New(machine.assembler.Program())
	if err := vmachine.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
