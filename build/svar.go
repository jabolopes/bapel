package build

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
)

type svar[T any] struct {
	name   string
	mutex  sync.Mutex
	result *T
	c      chan struct{}
}

func (v *svar[T]) setName(name string) *svar[T] {
	v.name = name
	return v
}

func (v *svar[T]) isSet() bool {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	return v.result != nil
}

func (v *svar[T]) get() T {
loop:
	for {
		select {
		case <-v.c:
			break loop
		case <-time.After(10 * time.Second):
			glog.V(1).Infof("Waiting for svar %q", v.name)
		}
	}

	return *v.result
}

func (v *svar[T]) getCtx(ctx context.Context) (T, error) {
	var t T

	select {
	case <-ctx.Done():
		return t, ctx.Err()
	case <-v.c:
		return *v.result, nil
	}
}

func (v *svar[T]) getErrCtx(ctx context.Context) error {
	_, err := v.getCtx(ctx)
	return err
}

func (v *svar[T]) set(value T) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if v.result != nil {
		return
	}

	v.result = &value
	close(v.c)
}

func newSvar[T any]() *svar[T] {
	channel := make(chan struct{})
	return &svar[T]{
		"",
		sync.Mutex{},
		nil, /* result */
		channel,
	}
}

func newValueSvar[T any](value T) *svar[T] {
	svar := newSvar[T]()
	svar.set(value)
	return svar
}

func getSvar[T any](svar *svar[any]) T {
	return svar.get().(T)
}

func getSvarCtx[T any](ctx context.Context, svar *svar[any]) (T, error) {
	var t T

	anyValue, err := svar.getCtx(ctx)
	if err != nil {
		return t, err
	}

	value, ok := anyValue.(T)
	if !ok {
		return t, fmt.Errorf("expected type %T; got type %T", t, anyValue)
	}

	return value, nil
}
