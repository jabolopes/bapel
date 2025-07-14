package build

import (
	"fmt"
	"path"
	"sync"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
)

func toOutputFilename(inputFilename, outputDirectory, outputBasename string) string {
	var extension string
	switch path.Ext(inputFilename) {
	case ".bpl":
		extension = ".ccm"
	case ".ccm":
		extension = ".pcm"
	case ".pcm":
		extension = ".o"
	case ".o":
		return inputFilename
	}

	return fmt.Sprintf("%s%s", path.Join(outputDirectory, outputBasename), extension)
}

type Builder struct {
	querier         query.Querier
	mutex           sync.Mutex
	moduleActions   map[ir.ModuleID]*action
	outputDirectory string
}

func (b *Builder) linkObjFiles(moduleID ir.ModuleID, allObjFiles, allFlags []string) error {
	// TODO: Extract this filename computation to a centralized place.
	outputFilename := path.Join(b.outputDirectory, moduleID.Name)
	if _, err := LinkObjsToExecutable(allObjFiles, allFlags, outputFilename); err != nil {
		return err
	}

	return nil
}

func (b *Builder) moduleActionImpl(a *action) error {
	a.addFieldVar("moduleFlags").
		addFieldVar("waitDepsPCMs")

	moduleID, err := getConstant[ir.ModuleID](a, "moduleID")
	if err != nil {
		return err
	}

	moduleQuery, err := b.querier.QueryModule(moduleID)
	if err != nil {
		return fmt.Errorf("failed to query module %q: %v", moduleID, err)
	}

	moduleBuilder := newModuleBuilder(
		b,
		a,
		a.outputVar("allPCMsDone"))

	moduleFlags := make([]string, 0, len(moduleQuery.Flags))
	for _, flag := range moduleQuery.Flags {
		moduleFlags = append(moduleFlags, flag.Value)
	}
	a.fieldVar("moduleFlags").set(moduleFlags)

	waitDepsPCMs := a.addBarrier().setDone(a.fieldVar("waitDepsPCMs"))
	for _, id := range moduleQuery.Imports {
		depAction := b.buildModule(a, id)
		moduleBuilder.allDeps.add(depAction)
		waitDepsPCMs.add(depAction.outputVar("allPCMsDone"))
	}

	{
		baseFilename := b.querier.BaseSourceFilename(moduleID)

		for i, relativeImplFilename := range moduleQuery.Impls {
			implFilename := b.querier.ImplSourceFilename(baseFilename, relativeImplFilename)

			outputBasename := fmt.Sprintf("%s-%s", moduleID.Name, parse.TrimExtension(path.Base(implFilename.Value)))
			moduleBuilder.compileToObj(implFilename.Value, outputBasename, i)
		}

		moduleBuilder.compileToObj(baseFilename.Value, moduleID.Name, len(moduleQuery.Impls))
	}

	moduleBuilder.computeAllObjs(a.outputVar("allObjFiles"), a.outputVar("allFlags"))

	return nil
}

func (b *Builder) buildModule(parentAction *action, moduleID ir.ModuleID) *action {
	moduleIDNoPos := moduleID
	moduleIDNoPos.Pos = ir.Pos{}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if moduleAction, ok := b.moduleActions[moduleIDNoPos]; ok {
		glog.V(1).Infof("Already built module %q", moduleID)
		return moduleAction
	}

	glog.V(1).Infof("Found new module %q", moduleID)

	moduleAction := parentAction.addChild(fmt.Sprintf("buildModule(%s)", moduleID)).
		addConstant("moduleID", moduleID).
		addConstant("outputDirectory", b.outputDirectory).
		addOutputVar("allPCMsDone").
		addOutputVar("allFlags").
		addOutputVar("allObjFiles").
		setImpl(b.moduleActionImpl).
		build()

	b.moduleActions[moduleIDNoPos] = moduleAction

	return moduleAction
}

func (b *Builder) Build(moduleID ir.ModuleID) error {
	moduleAction := b.buildModule(nil /* parentAction */, moduleID)
	if err := moduleAction.getErr(); err != nil {
		return fmt.Errorf("failed to build module %q: %v", moduleID, err)
	}

	allObjFiles, err := getOutputVar[[]string](moduleAction, "allObjFiles")
	if err != nil {
		return fmt.Errorf("failed to build module %q: %v", moduleID, err)
	}

	allFlags, err := getOutputVar[[]string](moduleAction, "allFlags")
	if err != nil {
		return fmt.Errorf("failed to build module %q: %v", moduleID, err)
	}

	if err := b.linkObjFiles(moduleID, allObjFiles, allFlags); err != nil {
		return fmt.Errorf("failed to program module %q: %v", moduleID, err)
	}

	return nil
}

func NewBuilder(querier query.Querier) *Builder {
	return &Builder{
		querier,
		sync.Mutex{},
		map[ir.ModuleID]*action{},
		"out", /* outputDirectory */
	}
}
