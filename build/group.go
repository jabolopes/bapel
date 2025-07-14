package build

type group struct {
	actions []*action
	doneVar *svar[any]
}

func (s *group) startImpl() *group {
	go func() {
		for _, action := range s.actions {
			_ = action.getErr()
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

type groupBuilder struct {
	builtGroup *group
	actions    []*action
	doneVar    *svar[any]
}

func (s *groupBuilder) add(action *action) *groupBuilder {
	if s.builtGroup != nil {
		panic("group is already built")
	}

	s.actions = append(s.actions, action)
	return s
}

func (s *groupBuilder) setDone(doneVar *svar[any]) *groupBuilder {
	if s.builtGroup != nil {
		panic("group is already built")
	}

	s.doneVar = doneVar
	return s
}

func (s *groupBuilder) build() *group {
	if s.builtGroup != nil {
		return s.builtGroup
	}

	if s.doneVar == nil {
		s.doneVar = newSvar[any]()
	}

	group := newGroup(s.actions, s.doneVar)
	s.builtGroup = group
	return group
}

func newGroupBuilder() *groupBuilder {
	return &groupBuilder{nil, nil, nil}
}
