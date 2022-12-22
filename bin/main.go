package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jabolopes/go-vm/asm"
	"github.com/jabolopes/go-vm/vm"
	"golang.org/x/exp/constraints"
)

// Instruction set:
//
// Types: i8 i16 i32 i64 u8 u16 u32 u64
//
// ALLOC <size:number>
// GET <type> <address>
// SET <type> <address> <value>

type OpCode struct {
	token    string
	callback func(*Machine, []string) error
}

type Machine struct {
	opcodes   []OpCode
	assembler *asm.Assembler
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

func assemblePushU8(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("'push u8' expects 1 argument; got %q", args)
	}

	value, err := parseNumber[byte](args[0])
	if err != nil {
		return err
	}

	machine.assembler.PutOpCode(vm.PushU8)
	machine.assembler.PutU8(value)
	return nil
}

func assemblePushU16(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("'push u16' expects 1 argument; got %q", args)
	}

	value, err := parseNumber[uint16](args[0])
	if err != nil {
		return err
	}

	machine.assembler.PutOpCode(vm.PushU16)
	machine.assembler.PutU16(value)
	return nil
}

func assemblePushU32(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("'push u32' expects 1 argument; got %q", args)
	}

	value, err := parseNumber[uint32](args[0])
	if err != nil {
		return err
	}

	machine.assembler.PutOpCode(vm.PushU32)
	machine.assembler.PutU32(value)
	return nil
}

func assemblePushU64(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("'push u64' expects 1 argument; got %q", args)
	}

	value, err := parseNumber[uint64](args[0])
	if err != nil {
		return err
	}

	machine.assembler.PutOpCode(vm.PushU64)
	machine.assembler.PutU64(value)
	return nil
}

func assemblePrintU8(machine *Machine, args []string) error {
	machine.assembler.PutOpCode(vm.PrintU8)
	return nil
}

func assemblePrintU16(machine *Machine, args []string) error {
	machine.assembler.PutOpCode(vm.PrintU16)
	return nil
}

func assemblePrintU32(machine *Machine, args []string) error {
	machine.assembler.PutOpCode(vm.PrintU32)
	return nil
}

func assemblePrintU64(machine *Machine, args []string) error {
	machine.assembler.PutOpCode(vm.PrintU64)
	return nil
}

func assembleOp(machine *Machine, line string) error {
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

func runFromFile(machine *Machine, input *os.File) error {
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
			{"push u8", assemblePushU8},
			{"push u16", assemblePushU16},
			{"push u32", assemblePushU32},
			{"push u64", assemblePushU64},
			{"print u8", assemblePrintU8},
			{"print u16", assemblePrintU16},
			{"print u32", assemblePrintU32},
			{"print u64", assemblePrintU64},
		},
		asm.New(),
	}
	if err := runFromFile(machine, os.Stdin); err != nil {
		return err
	}

	vmachine := vm.New(machine.assembler.Data())
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
