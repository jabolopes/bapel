package comp

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ir"
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

func compileOneFile(unit ir.IrUnit, mode PrinterMode, outputFilename string) error {
	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if err := PrintUnitToCpp(unit, mode, outputFile); err != nil {
		return err
	}

	if err := outputFile.Close(); err != nil {
		return err
	}

	return formatFile(outputFile.Name())
}

func deepCopyUnit(unit ir.IrUnit) ir.IrUnit {
	unitCopy := unit
	unitCopy.Decls = make([]ir.IrDecl, len(unit.Decls))
	for i := range unit.Decls {
		unitCopy.Decls[i] = unit.Decls[i].Clone()
	}
	unitCopy.Functions = make([]ir.IrFunction, len(unit.Functions))
	copy(unitCopy.Functions, unit.Functions)
	return unitCopy
}

func CompileBPL(querier query.Querier, inputFilename, outputFilename string) error {
	glog.V(1).Infof("Compiling %q with output base %q...", inputFilename, outputFilename)

	unit, err := TypecheckSourceFile(querier, TypecheckOptions{}, inputFilename)
	if err != nil {
		return err
	}

	baseOutputPath := strings.TrimSuffix(outputFilename, path.Ext(outputFilename))

	if unit.Case == ir.BaseUnit {
		// Generate M.h
		if err := compileOneFile(deepCopyUnit(unit), ModePublicHeader, baseOutputPath+".h"); err != nil {
			return err
		}
		// Generate M_private.h
		if err := compileOneFile(deepCopyUnit(unit), ModePrivateHeader, baseOutputPath+"_private.h"); err != nil {
			return err
		}
		// Generate M.cc
		if err := compileOneFile(deepCopyUnit(unit), ModeSource, baseOutputPath+".cc"); err != nil {
			return err
		}
	} else {
		// ImplUnit
		// Generate M_impl.cc
		if err := compileOneFile(deepCopyUnit(unit), ModeSource, baseOutputPath+".cc"); err != nil {
			return err
		}
	}

	return nil
}
