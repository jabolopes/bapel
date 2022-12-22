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
// Types: i8 i16 i32 i64 u8 u16 u32 u64
//
// push <type> <value>
// print <type>

type OpCode struct {
	token    string
	callback func(*Machine, []string) error
}

type OpLocal struct {
	offset uint16
	size   uint16
}

type OpFunction struct {
	locals        map[string]OpLocal
	currentOffset uint16
}

func (f *OpFunction) Local(id string, size uint16) {
	f.locals[id] = OpLocal{f.currentOffset, size}
	f.currentOffset += size
}

func NewFunction() *OpFunction {
	return &OpFunction{map[string]OpLocal{}, 0}
}

type Machine struct {
	opcodes      []OpCode
	mainFunction *OpFunction
	assembler    *asm.Assembler
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

		machine.mainFunction.Local(args[0], size)
		return nil
	}
}

func assembleLocalGet(machine *Machine, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expects 1 argument; got %q", args)
	}

	local, ok := machine.mainFunction.locals[args[0]]
	if !ok {
		return fmt.Errorf("Undeclared local %q", args[0])
	}

	switch local.size {
	case 1:
		machine.assembler.PutOpCode(vm.PushL8)
	case 2:
		machine.assembler.PutOpCode(vm.PushL16)
	case 4:
		machine.assembler.PutOpCode(vm.PushL32)
	case 8:
		machine.assembler.PutOpCode(vm.PushL64)
	}

	machine.assembler.PutU16(local.offset)

	return nil
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
			{"local u8", assembleLocal[byte]()},
			{"local u16", assembleLocal[uint16]()},
			{"local u32", assembleLocal[uint32]()},
			{"local u64", assembleLocal[uint64]()},

			{"local get", assembleLocalGet},

			{"push u8", assemblePushU8},
			{"push u16", assemblePushU16},
			{"push u32", assemblePushU32},
			{"push u64", assemblePushU64},
			{"print u8", assemblePrintU8},
			{"print u16", assemblePrintU16},
			{"print u32", assemblePrintU32},
			{"print u64", assemblePrintU64},
		},
		NewFunction(),
		asm.New(),
	}
	if err := assembleFile(machine, os.Stdin); err != nil {
		return err
	}

	program := vm.OpProgram{
		machine.assembler.Data(),
		[]vm.OpFunction{},
	}

	vmachine := vm.New(program)
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
