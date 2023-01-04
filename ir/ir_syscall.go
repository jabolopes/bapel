package ir

import "syscall"

var syscallNames = map[string]uint32{
	"time": syscall.SYS_TIME,
}

func GetSyscall(name string) (uint32, bool) {
	num, ok := syscallNames[name]
	return num, ok
}
