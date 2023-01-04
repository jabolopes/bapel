package vm

import (
	"fmt"
	"syscall"

	"github.com/jabolopes/bapel/ir"
)

func getErrno(err error) uint64 {
	if err == nil {
		return 0
	}

	errno, ok := err.(syscall.Errno)
	if !ok {
		panic(fmt.Errorf("Expected error of type %T; got %w (%T)", syscall.Errno(0), err, err))
	}

	return uint64(errno)
}

var syscalls = []opFunction{
	syscall.SYS_TIME: func(machine *Machine) error {
		time, err := syscall.Time(nil)
		machine.Stack().
			PushI64(uint64(time)).
			PushI64(getErrno(err))
		return nil
	},
}

func opSyscall(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			return syscalls[machine.Tape().GetI32()](machine)
		},
	}
}
