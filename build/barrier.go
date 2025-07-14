package build

import (
	"context"
)

type barrier struct {
	ctx      context.Context
	waitVars []*svar[any]
	doneVar  *svar[any]
}

func (s *barrier) startImpl() *barrier {
	go func() {
		for _, svar := range s.waitVars {
			svar.getCtx(s.ctx)
		}
		s.doneVar.set(struct{}{})
	}()

	return s
}

func newBarrier(ctx context.Context, waitVars []*svar[any], doneVar *svar[any]) *barrier {
	return (&barrier{ctx, waitVars, doneVar}).startImpl()
}

type barrierBuilder struct {
	ctx          context.Context
	builtBarrier *barrier
	waitVars     []*svar[any]
	doneVar      *svar[any]
}

func (s *barrierBuilder) add(svar *svar[any]) *barrierBuilder {
	if s.builtBarrier != nil {
		panic("barrier is already built")
	}

	s.waitVars = append(s.waitVars, svar)
	return s
}

func (s *barrierBuilder) setDone(doneVar *svar[any]) *barrierBuilder {
	if s.builtBarrier != nil {
		panic("barrier is already built")
	}

	s.doneVar = doneVar
	return s
}

func (s *barrierBuilder) build() *barrier {
	if s.builtBarrier != nil {
		return s.builtBarrier
	}

	if s.doneVar == nil {
		s.doneVar = newSvar[any]()
	}

	barrier := newBarrier(s.ctx, s.waitVars, s.doneVar)
	s.builtBarrier = barrier
	return barrier
}

func newBarrierBuilder(ctx context.Context) *barrierBuilder {
	return &barrierBuilder{ctx, nil, nil, nil}
}
