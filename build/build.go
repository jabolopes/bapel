package build

import (
	"fmt"
	"path"
	"strings"
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
		extension = ".cc"
	case ".cc", ".cpp":
		extension = ".o"
	case ".o":
		return inputFilename
	case ".h":
		return ""
	default:
		panic(fmt.Errorf("unhandled extension %q", path.Ext(inputFilename)))
	}

	if extension == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", path.Join(outputDirectory, outputBasename), extension)
}

func toBaseOutputFilename(moduleID ir.ModuleID) string {
	return strings.Replace(moduleID.Name, ir.ModuleIDSeparator, "/", -1)
}

func toImplOutputFilename(moduleID ir.ModuleID, implFilename ir.Filename) string {
	return fmt.Sprintf("%s-%s", toBaseOutputFilename(moduleID), parse.TrimExtension(path.Base(implFilename.Value)))
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
	cmd, err := LinkObjsToExecutable(allObjFiles, allFlags, outputFilename)
	if err != nil {
		return err
	}

	_, err = runCommand(b.outputDirectory, cmd)
	return err
}

func (b *Builder) moduleActionImpl(a *action) error {
	a.addFieldVar("moduleFlags").
		addFieldVar("waitDepsHeaders")

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
		a.outputVar("allHeadersDone"))

	moduleFlags := make([]string, 0, len(moduleQuery.Flags))
	for _, flag := range moduleQuery.Flags {
		moduleFlags = append(moduleFlags, flag.Value)
	}
	a.fieldVar("moduleFlags").set(moduleFlags)

	waitDepsHeaders := a.addBarrier().setDone(a.fieldVar("waitDepsHeaders"))
	for _, id := range moduleQuery.Imports {
		depAction := b.buildModule(a, id)
		moduleBuilder.allDeps.add(depAction)
		waitDepsHeaders.add(depAction.outputVar("allHeadersDone"))
	}

	{
		baseFilename := b.querier.BaseSourceFilename(moduleID)
		baseOutputBasename := toBaseOutputFilename(moduleID)

		// Compile base BPL first. This will set "allHeadersDone" when finished.
		baseBplAction := moduleBuilder.compileBPL(baseFilename.Value, baseOutputBasename, true /* isBase */)

		for _, relativeImplFilename := range moduleQuery.Impls {
			implFilename := b.querier.ImplSourceFilename(baseFilename, relativeImplFilename)
			ext := path.Ext(implFilename.Value)

			if ext == ".bpl" {
				implOutputBasename := toImplOutputFilename(moduleID, implFilename)
				implBplAction := moduleBuilder.compileBPL(implFilename.Value, implOutputBasename, false /* isBase */)
				moduleBuilder.compileBplCcToObj(implBplAction, baseBplAction, implOutputBasename)
			} else if ext == ".cc" || ext == ".cpp" {
				implOutputBasename := toImplOutputFilename(moduleID, implFilename)
				moduleBuilder.compileCppToObj(implFilename.Value, baseBplAction, implOutputBasename)
			} else if ext == ".h" {
				// C++ headers don't need compilation, they are included.
			}
		}

		// Compile base CC to Obj
		moduleBuilder.compileBplCcToObj(baseBplAction, baseBplAction, baseOutputBasename)
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
		addOutputVar("allHeadersDone").
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
