package build

import (
	"fmt"
	"os"
	"path"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/query"
)

type moduleActionDependencies struct {
	querier      query.Querier
	pcmSequencer *sequencer
}

func (d moduleActionDependencies) compileToCCMActionImpl(a *action) error {
	inputFilename, err := getSvar[string](a.inputVar("inputFilename"))
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
	if _, err := a.inputVar("waitDepsPCMs").get(); err != nil {
		return err
	}

	if _, err := a.inputVar("childDone").get(); err != nil {
		return err
	}

	{
		sequence, err := getConstant[int](a, "sequence")
		if err != nil {
			return err
		}

		d.pcmSequencer.wait(sequence)
		defer d.pcmSequencer.next()
	}

	inputFilename, err := getSvar[string](a.inputVar("inputFilename"))
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

	moduleFlags, err := getSvar[[]string](a.inputVar("moduleFlags"))
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
	if _, err := a.inputVar("childDone").get(); err != nil {
		return err
	}

	inputFilename, err := getSvar[string](a.inputVar("inputFilename"))
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

func newModuleActionDependencies(moduleBuilder *moduleBuilder) moduleActionDependencies {
	return moduleActionDependencies{moduleBuilder.builder.querier, moduleBuilder.pcmSequencer}
}
