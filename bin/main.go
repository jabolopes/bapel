package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/jabolopes/go-vm/asm"
	"github.com/jabolopes/go-vm/vm"
	"golang.org/x/exp/constraints"
)

// Instruction set:
//
// Types: i8 i16 i32 i64
//
// push <type> <value>
// print <type>

type OpCode struct {
	token    string
	callback func(*Machine, []string) error
}

type Machine struct {
	opcodes   []OpCode
	assembler *asm.OpAssembler
}

func parseNumber[T constraints.Integer](line string) (T, error) {
	var value T

	if strings.HasPrefix(line, "0x") {
		// Hexadecimal
		_, err := fmt.Sscanf(line, "0x%x", &value)

		return value, err
	}

	// Decimal.
	_, err := fmt.Sscanf(line, "%d", &value)
	return value, err
}

func assembleLocal[T constraints.Integer]() func(*Machine, []string) error {
	var value T
	size := uint16(unsafe.Sizeof(value))
	return func(machine *Machine, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expects 1 argument; got %q", args)
		}

		return machine.assembler.LocalDefine(args[0], size)
	}
}

func assemblePushI8(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	value, err := parseNumber[byte](args[0])
	if err != nil {
		return err
	}

	return machine.assembler.PushI8(value)
}

func assemblePushI16(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	value, err := parseNumber[uint16](args[0])
	if err != nil {
		return err
	}

	return machine.assembler.PushI16(value)
}

func assemblePushI32(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	value, err := parseNumber[uint32](args[0])
	if err != nil {
		return err
	}

	return machine.assembler.PushI32(value)
}

func assemblePushI64(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	value, err := parseNumber[uint64](args[0])
	if err != nil {
		return err
	}

	return machine.assembler.PushI64(value)
}

func assemblePushLocal(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	return machine.assembler.PushLocal(args[0])
}

func assemblePopLocal(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument; got %q", args)
	}

	return machine.assembler.PopLocal(args[0])
}

func assemblePrintI8(machine *Machine, args []string) error {
	machine.assembler.PrintI8()
	return nil
}

func assemblePrintI16(machine *Machine, args []string) error {
	machine.assembler.PrintI16()
	return nil
}

func assemblePrintI32(machine *Machine, args []string) error {
	machine.assembler.PrintI32()
	return nil
}

func assemblePrintI64(machine *Machine, args []string) error {
	machine.assembler.PrintI64()
	return nil
}

func assembleAddI8(machine *Machine, args []string) error {
	machine.assembler.AddI8()
	return nil
}

func assembleAddI16(machine *Machine, args []string) error {
	machine.assembler.AddI16()
	return nil
}

func assembleAddI32(machine *Machine, args []string) error {
	machine.assembler.AddI32()
	return nil
}

func assembleAddI64(machine *Machine, args []string) error {
	machine.assembler.AddI64()
	return nil
}

func assembleFunc(machine *Machine, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("expected 3 arguments; got %q", args)
	}

	id := args[0]

	argBytes, err := parseNumber[uint16](args[1])
	if err != nil {
		return err
	}

	localBytes, err := parseNumber[uint16](args[2])
	if err != nil {
		return err
	}

	return machine.assembler.Function(id, argBytes, localBytes)
}

func assembleEnd(machine *Machine, args []string) error {
	return machine.assembler.End()
}

func assembleOp(machine *Machine, line string) error {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil
	}

	for _, opcode := range machine.opcodes {
		if strings.HasPrefix(line, opcode.token) {
			line = strings.TrimPrefix(line, opcode.token)
			line = strings.TrimPrefix(line, " ")
			args := strings.Split(line, " ")
			return opcode.callback(machine, args)
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
	machine := &Machine{
		[]OpCode{
			{"local i8", assembleLocal[byte]()},
			{"local i16", assembleLocal[uint16]()},
			{"local i32", assembleLocal[uint32]()},
			{"local i64", assembleLocal[uint64]()},

			{"push i8", assemblePushI8},
			{"push i16", assemblePushI16},
			{"push i32", assemblePushI32},
			{"push i64", assemblePushI64},

			{"push", assemblePushLocal},
			{"pop", assemblePopLocal},

			{"print i8", assemblePrintI8},
			{"print i16", assemblePrintI16},
			{"print i32", assemblePrintI32},
			{"print i64", assemblePrintI64},

			{"add i8", assembleAddI8},
			{"add i16", assembleAddI16},
			{"add i32", assembleAddI32},
			{"add i64", assembleAddI64},

			{"func", assembleFunc},
			{"end", assembleEnd},
		},
		asm.New(),
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
