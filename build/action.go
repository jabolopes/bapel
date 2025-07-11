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
}

func (a *action) runImpl() {
	if err := a.impl(a); err != nil {
		for _, svar := range a.inputVars {
			svar.fail(err)
		}
		for _, svar := range a.fieldVars {
			svar.fail(err)
		}
		for _, svar := range a.outputVars {
			svar.fail(err)
		}
		a.doneVar.fail(err)
		return
	}

	doneErr := errors.New("done")
	for _, svar := range a.inputVars {
		svar.fail(doneErr)
	}
	for _, svar := range a.fieldVars {
		svar.fail(doneErr)
	}
	for _, svar := range a.outputVars {
		svar.fail(doneErr)
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
	impl       actionImpl
	constants  map[string]any
	inputVars  map[string]*svar[any]
	outputVars map[string]*svar[any]
}

func (a *actionBuilder) addConstant(name string, value any) *actionBuilder {
	if _, ok := a.constants[name]; ok {
		panic(fmt.Errorf("constant %q already defined", name))
	}

	a.constants[name] = value
	return a
}

func (a *actionBuilder) addInputVar(name string, svar *svar[any]) *actionBuilder {
	if _, ok := a.inputVars[name]; ok {
		panic(fmt.Errorf("input variable %q already defined", name))
	}

	a.inputVars[name] = svar
	return a
}

func (a *actionBuilder) addOutputVarTo(name string, svar *svar[any]) *actionBuilder {
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

func (a *actionBuilder) setImpl(impl actionImpl) *actionBuilder {
	a.impl = impl
	return a
}

func (a *actionBuilder) build() *action {
	return (&action{
		newSvar[any](),
		a.impl,
		a.constants,
		a.inputVars,
		map[string]*svar[any]{}, /* fieldVars */
		a.outputVars,
	}).startImpl()
}

func newActionBuilder() *actionBuilder {
	return &actionBuilder{
		func(*action) error { return nil },
		map[string]any{},        /* constants */
		map[string]*svar[any]{}, /* inputVars */
		map[string]*svar[any]{}, /* outputVars */
	}
}
