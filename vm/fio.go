package vm

import (
	"fmt"
	"sync"
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

func opWaitIO(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			ioID := machine.Stack().PopI64()
			ioOp := freeIoOp(int(ioID))
			machine.Stack().PushN(ioOp.get())
			return nil
		},
	}
}

func opDoIO(base OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			offset := machine.Tape().GetI16()
			greenThreadPc := machine.pc
			machine.pc += uint64(offset)

			ioID, ioOp := allocIoOp()
			machine.Stack().PushI64(uint64(ioID))

			go func() {
				greenThread := New(machine.program)
				greenThread.pc = greenThreadPc
				if err := greenThread.Run(); err != nil && err != errHalt {
					panic(fmt.Errorf("green thread failed: %w", err))
				}
				ioOp.put(greenThread.stack)
			}()

			return nil
		},
	}
}
