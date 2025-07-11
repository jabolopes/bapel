package build

import (
	"path"
)

type moduleBuilder struct {
	builder      *Builder
	moduleAction *action
	allDeps      *groupBuilder
	allPCMs      *groupBuilder
	allObjs      *groupBuilder
	pcmSequencer *sequencer
}

func (b *moduleBuilder) compileToCCM(inputFilename, outputBasename string) *action {
	return newActionBuilder().
		addConstant("outputDirectory", b.builder.outputDirectory).
		addConstant("outputBasename", outputBasename).
		addInputVar("inputFilename", newValueSvar[any](inputFilename)).
		addOutputVar("outputFilename").
		setImpl(func(a *action) error {
			return newModuleActionDependencies(b).compileToCCMActionImpl(a)
		}).
		build()
}

func (b *moduleBuilder) compileToPCM(inputFilename, outputBasename string, sequence int) *action {
	var childDoneVar *svar[any]
	var inputFilenameVar *svar[any]
	if path.Ext(inputFilename) == ".ccm" {
		childDoneVar = newValueSvar[any](struct{}{})
		inputFilenameVar = newValueSvar[any](inputFilename)
	} else {
		ccmAction := b.compileToCCM(inputFilename, outputBasename)
		childDoneVar = ccmAction.done()
		inputFilenameVar = ccmAction.outputVar("outputFilename")
	}

	pcmAction := newActionBuilder().
		addConstant("outputDirectory", b.builder.outputDirectory).
		addConstant("outputBasename", outputBasename).
		addConstant("sequence", sequence).
		addInputVar("waitDepsPCMs", b.moduleAction.fieldVar("waitDepsPCMs")).
		addInputVar("childDone", childDoneVar).
		addInputVar("inputFilename", inputFilenameVar).
		addInputVar("moduleFlags", b.moduleAction.fieldVar("moduleFlags")).
		addOutputVar("outputFilename").
		setImpl(func(a *action) error {
			return newModuleActionDependencies(b).compileToPCMActionImpl(a)
		}).
		build()

	b.allPCMs.add(pcmAction)

	return pcmAction
}

func (b *moduleBuilder) compileToObj(inputFilename, outputBasename string, sequence int) *action {
	var childDoneVar *svar[any]
	var inputFilenameVar *svar[any]
	if path.Ext(inputFilename) == ".pcm" {
		childDoneVar = newValueSvar[any](struct{}{})
		inputFilenameVar = newValueSvar[any](inputFilename)
	} else {
		pcmAction := b.compileToPCM(inputFilename, outputBasename, sequence)
		childDoneVar = pcmAction.done()
		inputFilenameVar = pcmAction.outputVar("outputFilename")
	}

	objAction := newActionBuilder().
		addConstant("outputDirectory", b.builder.outputDirectory).
		addConstant("outputBasename", outputBasename).
		addInputVar("childDone", childDoneVar).
		addInputVar("inputFilename", inputFilenameVar).
		addOutputVar("outputFilename").
		setImpl(func(a *action) error {
			return newModuleActionDependencies(b).compileToObjActionImpl(a)
		}).
		build()

	b.allObjs.add(objAction)

	return objAction
}

func (b *moduleBuilder) computeAllObjs(allObjsVar, allFlagsVar *svar[any]) *action {
	return newActionBuilder().
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

func newModuleBuilder(builder *Builder, moduleAction *action, allPCMs *svar[any]) *moduleBuilder {
	return &moduleBuilder{
		builder,
		moduleAction,
		newGroupBuilder(),
		newGroupBuilder().setDone(allPCMs),
		newGroupBuilder(),
		newSequencer(),
	}
}
