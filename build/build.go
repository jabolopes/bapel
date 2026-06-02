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
	targets         []BazelTarget
}

func (b *Builder) moduleActionImpl(a *action) error {
	a.addFieldVar("moduleFlags")

	moduleID, err := getConstant[ir.ModuleID](a, "moduleID")
	if err != nil {
		return err
	}

	moduleQuery, err := b.querier.QueryModule(moduleID)
	if err != nil {
		return fmt.Errorf("failed to query module %q: %v", moduleID, err)
	}

	var srcs []string
	var hdrs []string
	var deps []string

	baseFilename := b.querier.BaseSourceFilename(moduleID)
	baseOutputBasename := toBaseOutputFilename(moduleID)

	// Base BPL generated files are always expected.
	srcs = append(srcs, baseOutputBasename+".cc")
	hdrs = append(hdrs, baseOutputBasename+".h", baseOutputBasename+"_private.h")

	moduleBuilder := newModuleBuilder(b, a)

	moduleFlags := make([]string, 0, len(moduleQuery.Flags))
	for _, flag := range moduleQuery.Flags {
		moduleFlags = append(moduleFlags, flag.Value)
	}
	a.fieldVar("moduleFlags").set(moduleFlags)

	for _, id := range moduleQuery.Imports {
		depAction := b.buildModule(a, id)
		moduleBuilder.allDeps.add(depAction)

		// Collect dependency
		deps = append(deps, ":"+sanitizeTargetName(id.Name))
	}

	{
		// Compile base BPL first. This will set "allHeadersDone" when finished.
		moduleBuilder.compileBPL(baseFilename.Value, baseOutputBasename, true /* isBase */)

		for _, relativeImplFilename := range moduleQuery.Impls {
			implFilename := b.querier.ImplSourceFilename(baseFilename, relativeImplFilename)
			ext := path.Ext(implFilename.Value)

			if ext == ".bpl" {
				implOutputBasename := toImplOutputFilename(moduleID, implFilename)
				moduleBuilder.compileBPL(implFilename.Value, implOutputBasename, false /* isBase */)

				// Generated CC from impl BPL
				srcs = append(srcs, implOutputBasename+".cc")
			} else if ext == ".cc" || ext == ".cpp" {
				// Hand-written CC/CPP
				srcs = append(srcs, implFilename.Value)

				dst := path.Join(b.outputDirectory, implFilename.Value)
				if err := CopyFile(implFilename.Value, dst); err != nil {
					return fmt.Errorf("failed to copy %s to %s: %v", implFilename.Value, dst, err)
				}
			} else if ext == ".h" {
				// C++ headers don't need compilation, they are included.
				// Hand-written H
				hdrs = append(hdrs, implFilename.Value)

				dst := path.Join(b.outputDirectory, implFilename.Value)
				if err := CopyFile(implFilename.Value, dst); err != nil {
					return fmt.Errorf("failed to copy %s to %s: %v", implFilename.Value, dst, err)
				}
			}
		}
	}

	// Collect Bazel Target
	{
		isRoot, err := getConstant[bool](a, "isRoot")
		if err != nil {
			return err
		}

		targetName := sanitizeTargetName(moduleID.Name)
		targetType := "cc_library"
		if isRoot {
			targetType = "cc_binary"
			srcs = append(srcs, hdrs...)
			hdrs = nil
		}

		copts := append([]string{"-std=c++17"}, moduleFlags...)

		target := BazelTarget{
			Type:  targetType,
			Name:  targetName,
			Srcs:  srcs,
			Hdrs:  hdrs,
			Deps:  deps,
			Copts: copts,
		}
		b.addTarget(target)
	}

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
		addConstant("isRoot", parentAction == nil).
		setImpl(b.moduleActionImpl).
		build()

	b.moduleActions[moduleIDNoPos] = moduleAction

	return moduleAction
}

func (b *Builder) Build(moduleID ir.ModuleID) error {
	if err := EnsureBazelWorkspaceSetup(b.outputDirectory); err != nil {
		return err
	}

	moduleAction := b.buildModule(nil /* parentAction */, moduleID)
	if err := moduleAction.getErr(); err != nil {
		return fmt.Errorf("failed to build module %q: %v", moduleID, err)
	}

	if err := GenerateBuildFile(b.outputDirectory, b.targets); err != nil {
		return fmt.Errorf("failed to generate BUILD file: %v", err)
	}

	targetName := sanitizeTargetName(moduleID.Name)
	if err := RunBazelBuild(b.outputDirectory, targetName); err != nil {
		return fmt.Errorf("failed to build module %q with Bazel: %v", moduleID, err)
	}

	// Copy binary from bazel-bin to output location
	bazelBinPath := path.Join(b.outputDirectory, "bazel-bin", targetName)
	outputPath := path.Join(b.outputDirectory, moduleID.Name)
	if err := CopyFile(bazelBinPath, outputPath); err != nil {
		return fmt.Errorf("failed to copy built binary to %s: %v", outputPath, err)
	}

	return nil
}

func (b *Builder) addTarget(target BazelTarget) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.targets = append(b.targets, target)
}

func NewBuilder(querier query.Querier) *Builder {
	return &Builder{
		querier,
		sync.Mutex{},
		map[ir.ModuleID]*action{},
		"out", /* outputDirectory */
		nil,   /* targets */
	}
}
