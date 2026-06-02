package build

import (
	"fmt"
)


type moduleBuilder struct {
	builder      *Builder
	moduleAction *action
	allDeps      *groupBuilder
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



	return action
}


func newModuleBuilder(builder *Builder, moduleAction *action) *moduleBuilder {
	return &moduleBuilder{
		builder,
		moduleAction,
		moduleAction.addGroup(),
	}
}
