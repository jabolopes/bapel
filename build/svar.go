package build

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
)

type result[T any] struct {
	value T
	err   error
}

type svar[T any] struct {
	mutex  sync.Mutex
	result *result[T]
	c      chan result[T]
}

func (v *svar[T]) get() (T, error) {
	<-v.c
	return v.result.value, v.result.err
}

func (v *svar[T]) getErr() error {
	_, err := v.get()
	return err
}

func (v *svar[T]) set(value T) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if v.result != nil {
		return
	}

	v.result = &result[T]{value: value}
	close(v.c)
}

func (v *svar[T]) fail(err error) {
	if err == nil {
		glog.Error("attempted to set nil error on svar")
		return
	}

	v.mutex.Lock()
	defer v.mutex.Unlock()

	if v.result != nil {
		return
	}

	v.result = &result[T]{err: err}
	close(v.c)
}

func newSvar[T any]() *svar[T] {
	channel := make(chan result[T], 1)
	return &svar[T]{
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

func getSvar[T any](svar *svar[any]) (T, error) {
	var t T

	anyValue, err := svar.get()
	if err != nil {
		return t, err
	}

	value, ok := anyValue.(T)
	if !ok {
		return t, fmt.Errorf("expected type %T; got type %T", t, anyValue)
	}

	return value, nil
}
