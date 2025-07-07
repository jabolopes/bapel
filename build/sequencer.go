package build

import (
	"sync"
)

type sequencer struct {
	mutex   sync.Mutex
	current int
	waiters map[int]chan struct{}
}

func (s *sequencer) waitImpl(i int) chan struct{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.current == i {
		return nil
	}

	if c, ok := s.waiters[i]; ok {
		return c
	}

	c := make(chan struct{})
	s.waiters[i] = c
	return c
}

func (s *sequencer) wait(i int) {
	if c := s.waitImpl(i); c != nil {
		<-c
	}
}

func (s *sequencer) next() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.current++

	if c, ok := s.waiters[s.current]; ok {
		delete(s.waiters, s.current)
		close(c)
	}
}

func newSequencer() *sequencer {
	return &sequencer{
		sync.Mutex{},
		0,
		map[int](chan struct{}){},
	}
}
