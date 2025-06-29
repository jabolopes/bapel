package build

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/emirpasic/gods/v2/sets"
	"github.com/emirpasic/gods/v2/sets/hashset"
	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/comp"
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
			module.Flags.IDs = append(module.Flags.IDs, implModule.Flags.IDs...)
			module.Errors = append(module.Errors, implModule.Errors...)
		}
	}

	slices.SortFunc(module.Imports.IDs, func(id1, id2 ast.ID) int {
		return cmp.Compare(id1.Value, id2.Value)
	})
	module.Imports.IDs = slices.Compact(module.Imports.IDs)

	module.Flags.IDs = slices.Compact(module.Flags.IDs)

	return module, nil
}

func addSlash(p string) string {
	if strings.HasSuffix(p, "/") {
		return p
	}
	return p + "/"
}

func toOutputFilename(inputFilename, outputDirectory, moduleName string) string {
	var extension string
	switch path.Ext(inputFilename) {
	case ".bpl":
		extension = ".cc"
	case ".cc":
		extension = ".pcm"
	case ".pcm":
		extension = ".o"
	case ".o":
		return inputFilename
	}

	dir, base := path.Split(inputFilename)
	base = bplparser2.ReplaceExtension(base, extension)

	if !strings.HasPrefix(base, fmt.Sprintf("%s.", moduleName)) &&
		!strings.HasPrefix(base, fmt.Sprintf("%s-", moduleName)) {
		base = fmt.Sprintf("%s-%s", moduleName, base)
	}

	if !strings.HasPrefix(inputFilename, addSlash(outputDirectory)) {
		dir = path.Join(outputDirectory, dir)
	}

	return path.Join(dir, base)
}

type Builder struct {
	foundModules    sets.Set[string]
	outputDirectory string
	allObjFiles     []string
	allFlags        []string
}

func (b *Builder) runAction(moduleName string, flags []string, inputFilename string) (string, error) {
	outputFilename := toOutputFilename(inputFilename, b.outputDirectory, moduleName)

	glog.V(1).Infof("Compiling %q to %q", inputFilename, outputFilename)

	if path.Ext(inputFilename) == ".bpl" && path.Ext(outputFilename) == ".cc" {
		if err := comp.CompileBPLToCC(inputFilename, outputFilename); err != nil {
			return "", err
		}

		return b.runAction(moduleName, flags, outputFilename)
	}

	if path.Ext(inputFilename) == ".cc" && path.Ext(outputFilename) == ".pcm" {
		if _, err := CompileCCToPCM(inputFilename, flags, outputFilename); err != nil {
			return "", err
		}

		return outputFilename, nil
	}

	if path.Ext(inputFilename) == ".pcm" && path.Ext(outputFilename) == ".o" {
		if _, err := CompilePCMToObj(inputFilename, outputFilename); err != nil {
			return outputFilename, err
		}

		b.allObjFiles = append(b.allObjFiles, outputFilename)
		return outputFilename, nil
	}

	return "", fmt.Errorf("don't know how to compile file %q to file %q", inputFilename, outputFilename)
}

func (b *Builder) linkObjFiles(moduleName string) error {
	outputFilename := path.Join(b.outputDirectory, moduleName)
	if _, err := LinkObjsToExecutable(b.allObjFiles, b.allFlags, outputFilename); err != nil {
		return err
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

	var moduleFlags []string
	for _, flag := range module.Flags.IDs {
		moduleFlags = append(moduleFlags, flag.Value)
		b.allFlags = append(b.allFlags, flag.Value)
	}

	for _, imp := range module.Imports.IDs {
		if err := b.buildModule(imp.Value); err != nil {
			return err
		}
	}

	// Precompile sources to C++ precompiled modules.
	pcms := make([]string, 0, len(module.Impls.IDs)+1)
	for _, impl := range module.Impls.IDs {
		pcm, err := b.runAction(module.Header.Name, moduleFlags, impl.Value)
		if err != nil {
			return err
		}

		pcms = append(pcms, pcm)
	}
	{
		// Precompile base module source file to a C++ precompiled module.
		pcm, err := b.runAction(module.Header.Name, moduleFlags, inputFilename)
		if err != nil {
			return err
		}

		pcms = append(pcms, pcm)
	}

	// Compile modules to object files.
	for _, pcm := range pcms {
		if _, err := b.runAction(module.Header.Name, moduleFlags, pcm); err != nil {
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

	return nil
}

func (b *Builder) Build(inputFilename string) error {
	b.allObjFiles = b.allObjFiles[:0]
	b.allFlags = b.allFlags[:0]

	moduleName := bplparser2.TrimExtension(inputFilename)
	if err := b.buildModule(moduleName); err != nil {
		return err
	}

	return b.linkObjFiles(moduleName)
}

func NewBuilder() *Builder {
	return &Builder{
		hashset.New[string](),
		"out", /* outputDirectory */
		nil,   /* allObjFiles */
		nil,   /* allFlags */
	}
}
