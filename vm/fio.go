package vm

import (
	"fmt"
	"sync"
	"syscall"

	"github.com/jabolopes/bapel/ir"
)

type ioOp struct {
	ch   chan struct{}
	data []byte
}

func (f *ioOp) get() []byte {
	<-f.ch
	data := f.data
	f.ch = nil
	f.data = nil
	return data
}

func (f *ioOp) put(data []byte) {
	f.data = data
	close(f.ch)
}

func newIOOp() *ioOp {
	return &ioOp{make(chan struct{}), nil /* data */}
}

var (
	ioOps      = []*ioOp{}
	ioOpsMutex sync.Mutex
)

func allocIoOp() (int, *ioOp) {
	ioOpsMutex.Lock()
	defer ioOpsMutex.Unlock()

	if len(ioOps) >= cap(ioOps) {
		for i, op := range ioOps {
			if op == nil {
				ioOps[i] = newIOOp()
				return i, op
			}
		}
	}

	i := len(ioOps)
	op := newIOOp()
	ioOps = append(ioOps, op)
	return i, op
}

func freeIoOp(id int) *ioOp {
	var op *ioOp

	ioOpsMutex.Lock()
	op, ioOps[id] = ioOps[id], nil
	ioOpsMutex.Unlock()

	return op
}

func opWaitIO(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			ioID := machine.Stack().PopI64()
			ioOp := freeIoOp(int(ioID))
			machine.Stack().PushN(ioOp.get())
			return nil
		},
	}
}

func getErrorI64(err error) uint64 {
	if err == nil {
		return 0
	}

	errno, ok := err.(syscall.Errno)
	if !ok {
		panic(fmt.Errorf("Expected error of type %T; got %w (%T)", syscall.Errno(0), err, err))
	}

	return uint64(errno)
}

func opDoIO(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			ioID, ioOp := allocIoOp()
			machine.Stack().PushI64(uint64(ioID))

			go func() {
				time, err := syscall.Time(nil)
				data := ir.NewByteArrayEncoder().
					PutI64(uint64(time)).
					PutI64(getErrorI64(err)).
					Data()
				ioOp.put(data)
			}()

			return nil
		},
	}
}
