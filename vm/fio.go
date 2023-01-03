package vm

import (
	"errors"
	"fmt"
	"sync"
	"syscall"

	"github.com/jabolopes/bapel/ir"
)

type ioOp struct {
	ch chan []byte
}

func (f ioOp) ok() bool {
	return f.ch != nil
}

func (f ioOp) get() ([]byte, error) {
	data, ok := <-f.ch
	if !ok {
		return nil, errors.New("invalid IO operation channel")
	}
	return data, nil
}

func (f ioOp) put(t []byte) {
	f.ch <- t
	close(f.ch)
}

func newIOOp() ioOp {
	return ioOp{make(chan []byte, 1)}
}

func ioOpcode() ir.OpCode {
	return 0
}

var (
	ioOps      = []ioOp{}
	ioOpsMutex sync.Mutex
)

func allocIoOp() (int, ioOp) {
	ioOpsMutex.Lock()
	defer ioOpsMutex.Unlock()

	if len(ioOps) >= cap(ioOps) {
		for i := range ioOps {
			op := &ioOps[i]
			if !op.ok() {
				*op = newIOOp()
				return i, *op
			}
		}
	}

	i := len(ioOps)
	op := newIOOp()
	ioOps = append(ioOps, op)
	return i, op
}

func freeIoOp(id int) ioOp {
	var op ioOp
	var defaultOp ioOp

	ioOpsMutex.Lock()
	op, ioOps[id] = ioOps[id], defaultOp
	ioOpsMutex.Unlock()

	return op
}

func opWaitIO(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(machine *Machine) error {
			ioID := machine.Stack().PopI64()
			ioOp := freeIoOp(int(ioID))

			data, err := ioOp.get()
			if err != nil {
				return err
			}

			machine.Stack().PushN(data)
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
