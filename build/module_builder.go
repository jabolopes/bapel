package build

import (
	"fmt"
)


type moduleBuilder struct {
	builder        *Builder
	moduleAction   *action
	allDeps        *groupBuilder
	allObjs        *groupBuilder
	allHeadersDone *svar[any]
}

func (b *moduleBuilder) compileBPL(inputFilename, outputBasename string, isBase bool) *action {
	action := b.moduleAction.addChild(fmt.Sprintf("compileBPL(%q)", inputFilename)).
		addConstant("outputDirectory", b.builder.outputDirectory).
		addConstant("outputBasename", outputBasename).
		addConstant("isBase", isBase).
		addInputVar("inputFilename", newValueSvar[any](inputFilename)).
		addOutputVar("outputFilename").
		setImpl(func(a *action) error {
			return newModuleActionDependencies(b).compileBPLActionImpl(a)
		}).
		build()

	if isBase {
		go func() {
			val := action.doneVar().get()
			b.allHeadersDone.set(val)
		}()
	}

	return action
}

func (b *moduleBuilder) compileBplCcToObj(bplAction, baseBplAction *action, outputBasename string) *action {
	return b.moduleAction.addChild(fmt.Sprintf("compileBplCcToObj(%s)", outputBasename)).
		addGroupBuilder(b.allObjs).
		addConstant("outputDirectory", b.builder.outputDirectory).
		addConstant("outputBasename", outputBasename).
		addInputVar("childDone", bplAction.doneVar()).
		addInputVar("inputFilename", bplAction.outputVar("outputFilename")).
		addInputVar("baseHeadersDone", baseBplAction.doneVar()).
		addInputVar("waitDepsHeaders", b.moduleAction.fieldVar("waitDepsHeaders")).
		addInputVar("moduleFlags", b.moduleAction.fieldVar("moduleFlags")).
		addOutputVar("outputFilename").
		setImpl(func(a *action) error {
			return newModuleActionDependencies(b).compileCcToObjActionImpl(a)
		}).
		build()
}

func (b *moduleBuilder) compileCppToObj(inputFilename string, baseBplAction *action, outputBasename string) *action {
	return b.moduleAction.addChild(fmt.Sprintf("compileCppToObj(%q)", inputFilename)).
		addGroupBuilder(b.allObjs).
		addConstant("outputDirectory", b.builder.outputDirectory).
		addConstant("outputBasename", outputBasename).
		addInputVar("childDone", newValueSvar[any](struct{}{})).
		addInputVar("inputFilename", newValueSvar[any](inputFilename)).
		addInputVar("baseHeadersDone", baseBplAction.doneVar()).
		addInputVar("waitDepsHeaders", b.moduleAction.fieldVar("waitDepsHeaders")).
		addInputVar("moduleFlags", b.moduleAction.fieldVar("moduleFlags")).
		addOutputVar("outputFilename").
		setImpl(func(a *action) error {
			return newModuleActionDependencies(b).compileCcToObjActionImpl(a)
		}).
		build()
}

func (b *moduleBuilder) computeAllObjs(allObjsVar, allFlagsVar *svar[any]) *action {
	return b.moduleAction.addChild("computeAllObjs").
		addInputVar("moduleFlags", b.moduleAction.fieldVar("moduleFlags")).
		addInputVar("allDepsGroupDone", b.allDeps.build().done()).
		addInputVar("allObjsGroupDone", b.allObjs.build().done()).
		addOutputVarTo("allFlags", allFlagsVar).
		addOutputVarTo("allObjFiles", allObjsVar).
		setImpl(func(a *action) error {
			return newModuleActionDependencies(b).computeAllObjs(a)
		}).
		build()
}

func newModuleBuilder(builder *Builder, moduleAction *action, allHeadersDone *svar[any]) *moduleBuilder {
	return &moduleBuilder{
		builder,
		moduleAction,
		moduleAction.addGroup(),
		moduleAction.addGroup(),
		allHeadersDone,
	}
}
