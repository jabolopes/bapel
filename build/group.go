package build

type group struct {
	actions []*action
	doneVar *svar[any]
}

func (s *group) startImpl() *group {
	go func() {
		for _, action := range s.actions {
			if _, err := action.done().get(); err != nil {
				s.doneVar.fail(err)
				return
			}
		}
		s.doneVar.set(s.actions)
	}()

	return s
}

func (s *group) done() *svar[any] {
	return s.doneVar
}

func newGroup(actions []*action, doneVar *svar[any]) *group {
	return (&group{actions, doneVar}).startImpl()
}

func getGroup(group *group) ([]*action, error) {
	return getSvar[[]*action](group.done())
}

type groupBuilder struct {
	built   bool
	actions []*action
	doneVar *svar[any]
}

func (s *groupBuilder) add(action *action) *groupBuilder {
	if s.built {
		panic("group is already built")
	}

	s.actions = append(s.actions, action)
	return s
}

func (s *groupBuilder) setDone(doneVar *svar[any]) *groupBuilder {
	if s.built {
		panic("group is already built")
	}

	s.doneVar = doneVar
	return s
}

func (s *groupBuilder) build() *group {
	if s.built {
		panic("group is already built")
	}

	s.built = true

	if s.doneVar == nil {
		s.doneVar = newSvar[any]()
	}

	return newGroup(s.actions, s.doneVar)
}

func newGroupBuilder() *groupBuilder {
	return &groupBuilder{false, nil, nil}
}
