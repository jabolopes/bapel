package build

type barrier struct {
	waitVars []*svar[any]
	doneVar  *svar[any]
}

func (s *barrier) startImpl() *barrier {
	go func() {
		values := make([]any, 0, len(s.waitVars))
		for _, svar := range s.waitVars {
			value, err := svar.get()
			if err != nil {
				s.doneVar.fail(err)
				return
			}

			values = append(values, value)
		}
		s.doneVar.set(values)
	}()

	return s
}

// func (s *barrier) done() *svar[any] {
// 	return s.doneVar
// }

func newBarrier(waitVars []*svar[any], doneVar *svar[any]) *barrier {
	return (&barrier{waitVars, doneVar}).startImpl()
}

type barrierBuilder struct {
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

	barrier := newBarrier(s.waitVars, s.doneVar)
	s.builtBarrier = barrier
	return barrier
}

func newBarrierBuilder() *barrierBuilder {
	return &barrierBuilder{nil, nil, nil}
}
