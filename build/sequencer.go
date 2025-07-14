package build

import (
	"fmt"
	"sync"
)

type sequencer struct {
	mutex   sync.Mutex
	current int
	waiters map[int]*svar[any]
}

func (s *sequencer) waitImpl(i int, svar *svar[any]) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.current == i {
		svar.set(struct{}{})
		return
	}

	if _, ok := s.waiters[i]; ok {
		panic(fmt.Errorf("waiter %d is already registered", i))
	}

	s.waiters[i] = svar
}

func (s *sequencer) wait(i int, svar *svar[any]) error {
	s.waitImpl(i, svar)
	return svar.getErr()
}

func (s *sequencer) nextImpl() *svar[any] {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.current++

	svar, ok := s.waiters[s.current]
	if !ok {
		return nil
	}

	delete(s.waiters, s.current)
	return svar
}

func (s *sequencer) next() {
	svar := s.nextImpl()
	if svar != nil {
		svar.set(struct{}{})
	}
}

func newSequencer() *sequencer {
	return &sequencer{
		sync.Mutex{},
		0,
		map[int]*svar[any]{},
	}
}
