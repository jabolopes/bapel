package build

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/emirpasic/gods/v2/sets"
	"github.com/emirpasic/gods/v2/sets/hashset"
	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/query"
)

func toOutputFilename(moduleHeader ast.Header, inputFilename, outputDirectory string) string {
	var extension string
	switch path.Ext(inputFilename) {
	case ".bpl":
		extension = ".ccm"
	case ".ccm":
		extension = ".pcm"
	case ".pcm":
		extension = ".o"
	case ".o":
		return inputFilename
	}

	var basename string
	switch moduleHeader.Case {
	case ast.BaseFile:
		basename = moduleHeader.ModuleID.Name
	case ast.ImplementationFile:
		basename = path.Base(moduleHeader.Filename)
		basename = bplparser2.TrimExtension(basename)
		basename = fmt.Sprintf("%s-%s", moduleHeader.ModuleID.Name, basename)
	}

	return fmt.Sprintf("%s%s", path.Join(outputDirectory, basename), extension)
}

type Builder struct {
	querier         query.Querier
	foundModules    sets.Set[ast.ModuleID]
	outputDirectory string
	allObjFiles     []string
	allFlags        []string
}

// moduleName: name of the module (base file or implementation file),
// e.g., 'main', 'main_impl', etc.
func (b *Builder) runAction(moduleHeader ast.Header, flags []string, inputFilename string) (string, error) {
	outputFilename := toOutputFilename(moduleHeader, inputFilename, b.outputDirectory)

	glog.V(1).Infof("Compiling %q to %q", inputFilename, outputFilename)

	if err := os.MkdirAll(path.Dir(outputFilename), 0750); err != nil {
		return "", err
	}

	if path.Ext(inputFilename) == ".bpl" && path.Ext(outputFilename) == ".ccm" {
		if err := comp.CompileBPLToCCM(b.querier, inputFilename, outputFilename); err != nil {
			return "", err
		}

		return b.runAction(moduleHeader, flags, outputFilename)
	}

	if path.Ext(inputFilename) == ".ccm" && path.Ext(outputFilename) == ".pcm" {
		if _, err := CompileCCMToPCM(inputFilename, flags, outputFilename); err != nil {
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

func (b *Builder) linkObjFiles(moduleID ast.ModuleID) error {
	// TODO: Extract this filename computation to a centralized place.
	outputFilename := path.Join(b.outputDirectory, moduleID.Name)
	if _, err := LinkObjsToExecutable(b.allObjFiles, b.allFlags, outputFilename); err != nil {
		return err
	}

	return nil
}

func (b *Builder) buildModule(moduleID ast.ModuleID) error {
	moduleIDNoPos := moduleID
	moduleIDNoPos.Pos = ir.Pos{}

	if b.foundModules.Contains(moduleIDNoPos) {
		glog.V(1).Infof("Already built module %q", moduleID)
		return nil
	}

	glog.V(1).Infof("Found new module %q", moduleID)
	b.foundModules.Add(moduleIDNoPos)

	module, err := b.querier.QueryModuleMetadata(moduleID)
	if err != nil {
		return err
	}

	var moduleFlags []string
	for _, flag := range module.Flags.IDs {
		moduleFlags = append(moduleFlags, flag.Value)
		b.allFlags = append(b.allFlags, flag.Value)
	}

	for _, id := range module.Imports.IDs {
		if err := b.buildModule(id); err != nil {
			return err
		}
	}

	actions := make([]func() error, 0, len(module.Impls.IDs)+1)

	// Precompile sources to C++ precompiled modules.
	baseFilename := b.querier.ModuleBaseFilename(moduleID)

	for _, relativeImplFilename := range module.Impls.IDs {
		implFilename := b.querier.ModuleImplFilename(baseFilename, relativeImplFilename)

		header := module.Header
		header.Case = ast.ImplementationFile
		header.Filename = implFilename

		pcm, err := b.runAction(header, moduleFlags, implFilename)
		if err != nil {
			return err
		}

		actions = append(actions, func() error {
			_, err := b.runAction(header, moduleFlags, pcm)
			return err
		})
	}

	{
		// Precompile base module source file to a C++ precompiled module.
		pcm, err := b.runAction(module.Header, moduleFlags, baseFilename)
		if err != nil {
			return err
		}

		actions = append(actions, func() error {
			_, err := b.runAction(module.Header, moduleFlags, pcm)
			return err
		})
	}

	// Compile modules to object files.
	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}

	if !module.Valid() {
		var str strings.Builder
		str.WriteString(fmt.Sprintf("Failed to build %q:\n", moduleID))

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

func (b *Builder) Build(moduleID ast.ModuleID) error {
	b.allObjFiles = b.allObjFiles[:0]
	b.allFlags = b.allFlags[:0]

	if err := b.buildModule(moduleID); err != nil {
		return err
	}

	return b.linkObjFiles(moduleID)
}

func NewBuilder(querier query.Querier) *Builder {
	return &Builder{
		querier,
		hashset.New[ast.ModuleID](),
		"out", /* outputDirectory */
		nil,   /* allObjFiles */
		nil,   /* allFlags */
	}
}
