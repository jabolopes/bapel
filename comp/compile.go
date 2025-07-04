package comp

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/query"
)

func formatFile(filename string) error {
	cmd := exec.Command("clang-format", "-i", filename)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run %s: %s", cmd, output)
	}

	return nil
}

func CompileBPLToCCM(querier query.Querier, inputFilename, outputFilename string) error {
	glog.V(1).Infof("Compiling %q to %q...", inputFilename, outputFilename)

	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	module, err := CheckModule(querier, inputFilename)
	if err != nil {
		return err
	}

	if err := printModuleToCpp(module, outputFile); err != nil {
		return err
	}

	if err := outputFile.Close(); err != nil {
		return err
	}

	return formatFile(outputFile.Name())
}
