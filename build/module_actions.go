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
	querier query.Querier
}

func (d moduleActionDependencies) compileBPLActionImpl(a *action) error {
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

	if err := comp.CompileBPL(d.querier, inputFilename, outputFilename); err != nil {
		return fmt.Errorf("failed to compile %q to %q: %v", inputFilename, outputFilename, err)
	}

	return nil
}


func newModuleActionDependencies(moduleBuilder *moduleBuilder) moduleActionDependencies {
	return moduleActionDependencies{moduleBuilder.builder.querier}
}
