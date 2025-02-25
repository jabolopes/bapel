package comp

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/emirpasic/gods/v2/sets"
	"github.com/emirpasic/gods/v2/sets/hashset"
	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
)

func parseModuleNoBody(filename string) (ast.Module, error) {
	input, err := os.Open(filename)
	if err != nil {
		return ast.Module{}, err
	}
	defer input.Close()

	module, err := bplparser2.ParseFile(filename, input)
	if err != nil {
		return ast.Module{}, err
	}

	// TODO: At this stage, the builder only cares about the build graph, so we
	// could optimize the build process by not parsing the module body.
	module.Body = nil

	return module, nil
}

func parseModuleAndImplsNoBody(filename string) (ast.Module, error) {
	module, err := parseModuleNoBody(filename)
	if err != nil {
		return ast.Module{}, err
	}

	for _, filename := range module.Impls.IDs {
		if strings.HasSuffix(filename.Value, ".bpl") {
			implModule, err := parseModuleNoBody(filename.Value)
			if err != nil {
				return ast.Module{}, err
			}

			module.Imports.IDs = append(module.Imports.IDs, implModule.Imports.IDs...)
			module.Errors = append(module.Errors, implModule.Errors...)
		}
	}

	slices.SortFunc(module.Imports.IDs, func(id1, id2 ast.ID) int {
		return cmp.Compare(id1.Value, id2.Value)
	})
	module.Imports.IDs = slices.Compact(module.Imports.IDs)

	return module, nil
}

type Builder struct {
	foundModules sets.Set[string]
	allCcFiles   []string
}

func (b *Builder) compileImpl(inputFilename string) error {
	glog.V(1).Infof("Found module file %q", inputFilename)

	if strings.HasSuffix(inputFilename, ".cc") {
		b.allCcFiles = append(b.allCcFiles, inputFilename)
		_, err := os.Stat(inputFilename)
		return err
	}

	input, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer input.Close()

	outputFilename := bplparser2.ReplaceExtension(inputFilename, ".cc")
	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if err := CompileModule(inputFilename, input, outputFile); err != nil {
		return err
	}

	if err := outputFile.Close(); err != nil {
		return err
	}

	b.allCcFiles = append(b.allCcFiles, outputFilename)
	return nil
}

func (b *Builder) compileCcFiles(outputFilename string) error {
	if len(b.allCcFiles) == 0 {
		return fmt.Errorf("no cc files to build")
	}

	args := append([]string{"-std=c++20", "-fmodules-ts", "-o", outputFilename}, b.allCcFiles...)
	cmd := exec.Command("g++", args...)

	glog.V(1).Infof("Building program with %s", cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run %s: %s", cmd, output)
	}

	return nil
}

func (b *Builder) buildModule(moduleName string) error {
	if b.foundModules.Contains(moduleName) {
		glog.V(1).Infof("Already built module %q", moduleName)
		return nil
	}

	glog.V(1).Infof("Found new module %q", moduleName)
	b.foundModules.Add(moduleName)

	inputFilename := fmt.Sprintf("%s.bpl", moduleName)
	module, err := parseModuleAndImplsNoBody(inputFilename)
	if err != nil {
		return err
	}

	for _, imp := range module.Imports.IDs {
		if err := b.buildModule(imp.Value); err != nil {
			return err
		}
	}

	for _, impl := range module.Impls.IDs {
		if err := b.compileImpl(impl.Value); err != nil {
			return err
		}
	}

	if !module.Valid() {
		var str strings.Builder
		str.WriteString(fmt.Sprintf("Failed to build %q:\n", moduleName))

		firstErrors := module.Errors[:min(10, len(module.Errors))]
		interleave(firstErrors, func() { str.WriteString("\n\n") }, func(_ int, err ir.Error) {
			str.WriteString(err.String())
		})

		if len(module.Errors) > len(firstErrors) {
			str.WriteString("\n\nToo many errors to continue.")
		}

		return errors.New(str.String())
	}

	return b.compileImpl(inputFilename)
}

func (b *Builder) Build(inputFilename string) error {
	b.allCcFiles = b.allCcFiles[:0]

	moduleName := bplparser2.TrimExtension(inputFilename)
	if err := b.buildModule(moduleName); err != nil {
		return err
	}

	return b.compileCcFiles(moduleName)
}

func NewBuilder() *Builder {
	return &Builder{
		hashset.New[string](),
		nil,
	}
}
