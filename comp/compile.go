package comp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func CompileBPLDirect(querier query.Querier, inputFilename, outputFilename string) error {
	glog.V(1).Infof("Compiling %q with output base %q...", inputFilename, outputFilename)

	unit, err := TypecheckSourceFile(querier, TypecheckOptions{}, inputFilename)
	if err != nil {
		return err
	}

	baseOutputPath := strings.TrimSuffix(outputFilename, filepath.Ext(outputFilename))

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

func findWorkspaceRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("failed to find workspace root (go.mod)")
}

func CompileBPL(querier query.Querier, inputFilename, outputFilename string) error {
	workspaceRoot, err := findWorkspaceRoot()
	if err != nil {
		return err
	}
	compilerPath := filepath.Join(workspaceRoot, "bootstrap/compiler")
	if _, err := os.Stat(compilerPath); err != nil {
		return fmt.Errorf("compiler not found at %s; run 'make bootstrap/compiler' first", compilerPath)
	}

	absInput, err := filepath.Abs(inputFilename)
	if err != nil {
		return err
	}
	absOutput, err := filepath.Abs(outputFilename)
	if err != nil {
		return err
	}

	relInput, err := filepath.Rel(workspaceRoot, absInput)
	if err != nil {
		return err
	}
	relOutput, err := filepath.Rel(workspaceRoot, absOutput)
	if err != nil {
		return err
	}

	glog.V(1).Infof("Shelling out to compiler %s -o %s %s (Dir: %s)", compilerPath, relOutput, relInput, workspaceRoot)
	cmd := exec.Command(compilerPath, "-o", relOutput, relInput)
	cmd.Dir = workspaceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run %s: %s: %w", cmd, output, err)
	}

	return nil
}

