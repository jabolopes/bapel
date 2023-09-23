package vm

import (
	"fmt"
	"syscall"
)

func getErrno(err error) uint64 {
	if err == nil {
		return 0
	}

	errno, ok := err.(syscall.Errno)
	if !ok {
		panic(fmt.Errorf("expected error of type %T; got %w (%T)", syscall.Errno(0), err, err))
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

func opSyscall(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			return syscalls[machine.Tape().GetI32()](machine)
		},
	}
}
