package build

import (
	"context"
	"fmt"

	"github.com/golang/glog"
)

type actionImpl func(*action) error

type action struct {
	ctx           context.Context
	cancelFunc    func()
	name          string
	actionDoneVar *svar[any]
	actionErrVar  *svar[any]
	impl          actionImpl
	constants     map[string]any
	inputVars     map[string]*svar[any]
	fieldVars     map[string]*svar[any]
	outputVars    map[string]*svar[any]
	barriers      []*barrierBuilder
	groups        []*groupBuilder
	children      *groupBuilder
}

func (a *action) runImpl() {
	glog.V(1).Infof("%q: action started", a.name)

	implErr := a.impl(a)

	glog.V(1).Infof("%q: action finished with %v", a.name, implErr)

	for _, barrier := range a.barriers {
		_ = barrier.build()
	}

	for _, group := range a.groups {
		_ = group.build()
	}

	_ = a.children.build()

	if implErr != nil {
		a.cancel()
	}

	{
		actions := getSvar[[]*action](a.children.build().done())
		for _, action := range actions {
			if err := action.getErr(); err != nil {
				implErr = JoinErrors(implErr, err)
			}
		}
	}

	glog.V(1).Infof("%q: action's children finished with %v", a.name, implErr)

	if implErr != nil {
		a.actionDoneVar.set(struct{}{})
		a.actionErrVar.set(implErr)
		a.cancel()
		return
	}

	for name, svar := range a.outputVars {
		if !svar.isSet() {
			panic(fmt.Errorf("%q: output variable %q was not set", a.name, name))
		}
	}

	a.actionDoneVar.set(struct{}{})
	a.actionErrVar.set(nil)
}

func (a *action) startImpl() *action {
	go func() { a.runImpl() }()
	return a
}

func (a *action) cancel() {
	a.cancelFunc()
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

func (a *action) addBarrier() *barrierBuilder {
	barrierBuilder := newBarrierBuilder(a.ctx)
	a.barriers = append(a.barriers, barrierBuilder)
	return barrierBuilder
}

func (a *action) addGroup() *groupBuilder {
	groupBuilder := newGroupBuilder()
	a.groups = append(a.groups, groupBuilder)
	return groupBuilder
}

func (a *action) addChild(name string) *actionBuilder {
	return newActionBuilder(a, name)
}

func (a *action) doneVar() *svar[any] {
	return a.actionDoneVar
}

func (a *action) getErr() error {
	anyValue := a.actionErrVar.get()
	if anyValue == nil {
		return nil
	}

	return anyValue.(error)
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

func getInputVar[T any](a *action, name string) (T, error) {
	var t T

	value, err := getSvarCtx[T](a.ctx, a.inputVar(name))
	if err != nil {
		return t, errCancelled
	}

	return value, nil
}

func getInputVarErr(a *action, name string) error {
	if err := a.inputVar(name).getErrCtx(a.ctx); err != nil {
		return errCancelled
	}

	return nil
}

func getOutputVar[T any](a *action, name string) (T, error) {
	var t T

	value, err := getSvarCtx[T](a.ctx, a.outputVar(name))
	if err != nil {
		return t, errCancelled
	}

	return value, nil
}

func getGroupInputVar(a *action, name string) ([]*action, error) {
	return getInputVar[[]*action](a, name)
}

type actionBuilder struct {
	builtAction   *action
	ctx           context.Context
	cancelFunc    context.CancelFunc
	name          string
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
		a.ctx,
		a.cancelFunc,
		a.name,
		newSvar[any]().setName(fmt.Sprintf("%s.done", a.name)),
		newSvar[any]().setName(fmt.Sprintf("%s.doneErr", a.name)),
		a.impl,
		a.constants,
		a.inputVars,
		map[string]*svar[any]{}, /* fieldVars */
		a.outputVars,
		nil, /* barriers */
		nil, /* groups */
		newGroupBuilder(),
	}

	a.builtAction = newAction

	for _, groupBuilder := range a.groupBuilders {
		groupBuilder.add(newAction)
	}

	newAction.startImpl()
	return newAction
}

func newActionBuilder(parentAction *action, name string) *actionBuilder {
	if parentAction == nil {
		// TODO: Replace context.TODO.
		ctx, cancelFunc := context.WithCancel(context.TODO())

		return &actionBuilder{
			nil,
			ctx,
			cancelFunc,
			name,
			func(*action) error { return nil }, /* impl */
			map[string]any{},                   /* constants */
			map[string]*svar[any]{},            /* inputVars */
			map[string]*svar[any]{},            /* outputVars */
			nil,
		}
	}

	return &actionBuilder{
		nil, /* builtAction */
		parentAction.ctx,
		parentAction.cancelFunc,
		fmt.Sprintf("%s->%s", parentAction.name, name), /* name */
		func(*action) error { return nil },             /* impl */
		map[string]any{},                               /* constants */
		map[string]*svar[any]{},                        /* inputVars */
		map[string]*svar[any]{},                        /* outputVars */
		[]*groupBuilder{parentAction.children},
	}
}
