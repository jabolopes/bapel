package build

import (
	"errors"
	"fmt"
)

type actionImpl func(*action) error

type action struct {
	doneVar    *svar[any]
	impl       actionImpl
	constants  map[string]any
	inputVars  map[string]*svar[any]
	fieldVars  map[string]*svar[any]
	outputVars map[string]*svar[any]
	children   *groupBuilder
}

func (a *action) runImpl() {
	implErr := errors.Join(a.impl(a), a.children.build().done().getErr())
	actionErr := errors.Join(implErr, errors.New("done"))

	for _, svar := range a.inputVars {
		svar.fail(actionErr)
	}
	for _, svar := range a.fieldVars {
		svar.fail(actionErr)
	}
	for _, svar := range a.outputVars {
		svar.fail(actionErr)
	}

	if implErr != nil {
		a.doneVar.fail(implErr)
		return
	}

	a.doneVar.set(struct{}{})
}

func (a *action) startImpl() *action {
	go func() { a.runImpl() }()
	return a
}

func (a *action) inputVar(name string) *svar[any] {
	svar, ok := a.inputVars[name]
	if !ok {
		panic(fmt.Sprintf("undefined input variable %q", name))
	}
	return svar
}

func (a *action) addFieldVar(name string) *action {
	a.fieldVars[name] = newSvar[any]()
	return a
}

func (a *action) fieldVar(name string) *svar[any] {
	svar, ok := a.fieldVars[name]
	if !ok {
		panic(fmt.Sprintf("undefined field variable %q", name))
	}
	return svar
}

func (a *action) outputVar(name string) *svar[any] {
	svar, ok := a.outputVars[name]
	if !ok {
		panic(fmt.Sprintf("undefined output variable %q", name))
	}
	return svar
}

func (a *action) addChild() *actionBuilder {
	return newActionBuilder().
		addGroupBuilder(a.children)
}

func (a *action) done() *svar[any] {
	return a.doneVar
}

func getConstant[T any](a *action, name string) (T, error) {
	var t T

	anyConstant, ok := a.constants[name]
	if !ok {
		return t, fmt.Errorf("undefined constant %q", name)
	}

	value, ok := anyConstant.(T)
	if !ok {
		return t, fmt.Errorf("expected type %T; got type %T", t, anyConstant)
	}

	return value, nil
}

type actionBuilder struct {
	builtAction   *action
	impl          actionImpl
	constants     map[string]any
	inputVars     map[string]*svar[any]
	outputVars    map[string]*svar[any]
	groupBuilders []*groupBuilder
}

func (a *actionBuilder) addConstant(name string, value any) *actionBuilder {
	if a.builtAction != nil {
		panic("action is already built")
	}

	if _, ok := a.constants[name]; ok {
		panic(fmt.Errorf("constant %q already defined", name))
	}

	a.constants[name] = value
	return a
}

func (a *actionBuilder) addInputVar(name string, svar *svar[any]) *actionBuilder {
	if a.builtAction != nil {
		panic("action is already built")
	}

	if _, ok := a.inputVars[name]; ok {
		panic(fmt.Errorf("input variable %q already defined", name))
	}

	a.inputVars[name] = svar
	return a
}

func (a *actionBuilder) addOutputVarTo(name string, svar *svar[any]) *actionBuilder {
	if a.builtAction != nil {
		panic("action is already built")
	}

	if _, ok := a.outputVars[name]; ok {
		panic(fmt.Errorf("output variable %q already defined", name))
	}

	if svar == nil {
		svar = newSvar[any]()
	}

	a.outputVars[name] = svar
	return a
}

func (a *actionBuilder) addOutputVar(name string) *actionBuilder {
	return a.addOutputVarTo(name, nil)
}

func (a *actionBuilder) addGroupBuilder(groupBuilder *groupBuilder) *actionBuilder {
	if a.builtAction != nil {
		panic("action is already built")
	}

	a.groupBuilders = append(a.groupBuilders, groupBuilder)
	return a
}

func (a *actionBuilder) setImpl(impl actionImpl) *actionBuilder {
	if a.builtAction != nil {
		panic("action is already built")
	}

	a.impl = impl
	return a
}

func (a *actionBuilder) build() *action {
	if a.builtAction != nil {
		return a.builtAction
	}

	newAction := &action{
		newSvar[any](),
		a.impl,
		a.constants,
		a.inputVars,
		map[string]*svar[any]{}, /* fieldVars */
		a.outputVars,
		newGroupBuilder(),
	}

	for _, groupBuilder := range a.groupBuilders {
		groupBuilder.add(newAction)
	}

	newAction.startImpl()
	return newAction
}

func newActionBuilder() *actionBuilder {
	return &actionBuilder{
		nil, /* builtAction */
		func(*action) error { return nil },
		map[string]any{},        /* constants */
		map[string]*svar[any]{}, /* inputVars */
		map[string]*svar[any]{}, /* outputVars */
		nil,
	}
}
