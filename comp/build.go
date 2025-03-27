package comp

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
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

func toOutputFilename(inputFilename, outputDirectory, moduleName, extension string) string {
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

func (b *Builder) precompile(moduleName string, flags []string, inputFilename string) (string, error) {
	if strings.HasSuffix(inputFilename, ".cc") {
		// Example:
		// $ clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out -Ientt/single_include -ISDL/include game_impl.cc --precompile -o out/game-game_impl.pcm

		outputFilename := toOutputFilename(inputFilename, b.outputDirectory, moduleName, ".pcm")

		glog.V(1).Infof("Compiling %q to %q...", inputFilename, outputFilename)

		args := []string{"-std=c++20", "-x", "c++-module", fmt.Sprintf("-fprebuilt-module-path=%s", b.outputDirectory), inputFilename, "--precompile", "-o", outputFilename}
		args = append(args, flags...)
		cmd := exec.Command("clang++", args...)

		glog.V(1).Infof("Calling %s", cmd)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run %s: %s", cmd, output)
		}

		return outputFilename, nil
	}

	if strings.HasSuffix(inputFilename, ".bpl") {
		outputFilename := toOutputFilename(inputFilename, b.outputDirectory, moduleName, ".cc")

		glog.V(1).Infof("Compiling %q to %q...", inputFilename, outputFilename)

		input, err := os.Open(inputFilename)
		if err != nil {
			return "", err
		}
		defer input.Close()

		outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return "", err
		}
		defer outputFile.Close()

		if err := CompileModule(inputFilename, input, outputFile); err != nil {
			return "", err
		}

		if err := outputFile.Close(); err != nil {
			return "", err
		}

		return b.precompile(moduleName, flags, outputFilename)
	}

	return "", fmt.Errorf("don't know how to precompile file %q with unknown extension", inputFilename)
}

func (b *Builder) compileImpl(moduleName string, flags []string, inputFilename string) error {
	if strings.HasSuffix(inputFilename, ".o") {
		glog.V(1).Infof("Found %q...", inputFilename)

		b.allObjFiles = append(b.allObjFiles, inputFilename)
		_, err := os.Stat(inputFilename)
		return err
	}

	if strings.HasSuffix(inputFilename, ".pcm") {
		// Example:
		// $ clang++ -std=c++20 -fprebuilt-module-path=out -c out/game-game_impl.pcm -o out/game-game_impl.o

		outputFilename := toOutputFilename(inputFilename, b.outputDirectory, moduleName, ".o")

		glog.V(1).Infof("Compiling %q to %q...", inputFilename, outputFilename)

		args := []string{"-std=c++20", fmt.Sprintf("-fprebuilt-module-path=%s", b.outputDirectory), "-c", inputFilename, "-o", outputFilename}
		cmd := exec.Command("clang++", args...)

		glog.V(1).Infof("Calling %s", cmd)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run %s: %s", cmd, output)
		}

		return b.compileImpl(moduleName, flags, outputFilename)
	}

	return fmt.Errorf("don't know how to compile file %q with unknown extension", inputFilename)
}

func (b *Builder) linkObjFiles(moduleName string) error {
	if len(b.allObjFiles) == 0 {
		return fmt.Errorf("no cc files to build")
	}

	// Example:
	// clang++ -std=c++20 -fprebuilt-module-path=out -o out/program \
	//   -Wl,-rpath,SDL/build \
	//   -LSDL/build -lSDL3 \
	//   out/arr-arr_impl.o \
	//   ...

	outputFilename := path.Join(b.outputDirectory, moduleName)

	args := []string{"-std=c++20", fmt.Sprintf("-fprebuilt-module-path=%s", b.outputDirectory), "-o", outputFilename}
	args = append(args, b.allFlags...)
	args = append(args, b.allObjFiles...)
	cmd := exec.Command("clang++", args...)

	glog.V(1).Infof("Building program %q with %s", outputFilename, cmd)

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
	pcms := make([]string, 0, len(module.Impls.IDs))
	for _, impl := range module.Impls.IDs {
		pcm, err := b.precompile(module.Header.Name, moduleFlags, impl.Value)
		if err != nil {
			return err
		}

		pcms = append(pcms, pcm)
	}

	// Compile modules to object files.
	for _, pcm := range pcms {
		if err := b.compileImpl(module.Header.Name, moduleFlags, pcm); err != nil {
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

	// Precompile base module source file tp a C++ precompiled module.
	pcm, err := b.precompile(module.Header.Name, moduleFlags, inputFilename)
	if err != nil {
		return err
	}

	// Compule module to object file.
	return b.compileImpl(module.Header.Name, moduleFlags, pcm)
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
