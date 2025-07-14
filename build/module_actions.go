package build

import (
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/query"
)

type moduleActionDependencies struct {
	querier      query.Querier
	pcmSequencer *sequencer
}

func (d moduleActionDependencies) compileToCCMActionImpl(a *action) error {
	inputFilename, err := getInputVar[string](a, "inputFilename")
	if err != nil {
		return err
	}

	if path.Ext(inputFilename) != ".bpl" {
		return fmt.Errorf("expected file with extension '.bpl'; got %q", inputFilename)
	}

	outputDirectory, err := getConstant[string](a, "outputDirectory")
	if err != nil {
		return err
	}

	outputBasename, err := getConstant[string](a, "outputBasename")
	if err != nil {
		return err
	}

	outputFilename := toOutputFilename(inputFilename, outputDirectory, outputBasename)
	a.outputVar("outputFilename").set(outputFilename)

	if err := os.MkdirAll(path.Dir(outputFilename), 0750); err != nil {
		return err
	}

	glog.V(1).Infof("Compiling %q to %q", inputFilename, outputFilename)

	return comp.CompileBPLToCCM(d.querier, inputFilename, outputFilename)
}

func (d moduleActionDependencies) compileToPCMActionImpl(a *action) error {
	_ = a.addFieldVar("sequencer")

	{
		sequence, err := getConstant[int](a, "sequence")
		if err != nil {
			return err
		}

		if err := d.pcmSequencer.wait(sequence, a.fieldVar("sequencer")); err != nil {
			return err
		}
		defer d.pcmSequencer.next()
	}

	if err := getInputVarErr(a, "waitDepsPCMs"); err != nil {
		return err
	}

	if err := getInputVarErr(a, "childDone"); err != nil {
		return err
	}

	inputFilename, err := getInputVar[string](a, "inputFilename")
	if err != nil {
		return err
	}

	outputDirectory, err := getConstant[string](a, "outputDirectory")
	if err != nil {
		return err
	}

	outputBasename, err := getConstant[string](a, "outputBasename")
	if err != nil {
		return err
	}

	outputFilename := toOutputFilename(inputFilename, outputDirectory, outputBasename)
	a.outputVar("outputFilename").set(outputFilename)

	if err := os.MkdirAll(path.Dir(outputFilename), 0750); err != nil {
		return err
	}

	glog.V(1).Infof("Compiling %q to %q", inputFilename, outputFilename)

	moduleFlags, err := getInputVar[[]string](a, "moduleFlags")
	if err != nil {
		return err
	}

	cmd, err := CompileCCMToPCMCommand(inputFilename, moduleFlags, outputFilename)
	if err != nil {
		return err
	}

	_, err = runCommand(cmd)
	return err
}

func (d moduleActionDependencies) compileToObjActionImpl(a *action) error {
	if err := getInputVarErr(a, "childDone"); err != nil {
		return err
	}

	inputFilename, err := getInputVar[string](a, "inputFilename")
	if err != nil {
		return err
	}

	outputDirectory, err := getConstant[string](a, "outputDirectory")
	if err != nil {
		return err
	}

	outputBasename, err := getConstant[string](a, "outputBasename")
	if err != nil {
		return err
	}

	outputFilename := toOutputFilename(inputFilename, outputDirectory, outputBasename)
	a.outputVar("outputFilename").set(outputFilename)

	if err := os.MkdirAll(path.Dir(outputFilename), 0750); err != nil {
		return err
	}

	glog.V(1).Infof("Compiling %q to %q", inputFilename, outputFilename)

	cmd, err := CompilePCMToObjCommand(inputFilename, outputFilename)
	if err != nil {
		return err
	}

	_, err = runCommand(cmd)
	return err
}

func (d moduleActionDependencies) computeAllObjs(a *action) error {
	var allObjFiles []string
	{
		// Compute output variable 'allFlags'.
		//
		// Partially compute output variable 'allObjFiles'.
		allFlags, err := getInputVar[[]string](a, "moduleFlags")
		if err != nil {
			return err
		}

		allDepsActions, err := getGroupInputVar(a, "allDepsGroupDone")
		if err != nil {
			return err
		}

		for _, depAction := range allDepsActions {
			objFiles, err := getOutputVar[[]string](depAction, "allObjFiles")
			if err != nil {
				return err
			}

			flags, err := getOutputVar[[]string](depAction, "allFlags")
			if err != nil {
				return err
			}

			allObjFiles = append(allObjFiles, objFiles...)
			allFlags = append(allFlags, flags...)
		}

		a.outputVar("allFlags").set(allFlags)
	}

	{
		// Compute output variable 'allObjFiles'.
		allObjsActions, err := getGroupInputVar(a, "allObjsGroupDone")
		if err != nil {
			return err
		}

		for _, objAction := range allObjsActions {
			objFile, err := getOutputVar[string](objAction, "outputFilename")
			if err != nil {
				return err
			}

			allObjFiles = append(allObjFiles, objFile)
		}

		slices.Sort(allObjFiles)
		allObjFiles = slices.Compact(allObjFiles)

		a.outputVar("allObjFiles").set(allObjFiles)
	}

	return nil
}

func newModuleActionDependencies(moduleBuilder *moduleBuilder) moduleActionDependencies {
	return moduleActionDependencies{moduleBuilder.builder.querier, moduleBuilder.pcmSequencer}
}
