package build

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
)

func toOutputFilename(inputFilename, outputDirectory, outputBasename string) string {
	var extension string
	switch path.Ext(inputFilename) {
	case ".bpl":
		extension = ".cc"
	case ".cc", ".cpp":
		extension = ".o"
	case ".o":
		return inputFilename
	case ".h":
		return ""
	default:
		panic(fmt.Errorf("unhandled extension %q", path.Ext(inputFilename)))
	}

	if extension == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", path.Join(outputDirectory, outputBasename), extension)
}

func toBaseOutputFilename(moduleID ir.ModuleID) string {
	return strings.Replace(moduleID.Name, ir.ModuleIDSeparator, "/", -1)
}

func toImplOutputFilename(moduleID ir.ModuleID, implFilename ir.Filename) string {
	return fmt.Sprintf("%s-%s", toBaseOutputFilename(moduleID), parse.TrimExtension(path.Base(implFilename.Value)))
}

type Builder struct {
	querier         query.Querier
	mutex           sync.Mutex
	builtModules    map[ir.ModuleID]error
	outputDirectory string
	targets         []BazelTarget
}

func NewBuilder(querier query.Querier) *Builder {
	return &Builder{
		querier:         querier,
		builtModules:    make(map[ir.ModuleID]error),
		outputDirectory: "out",
	}
}

func (b *Builder) Build(moduleID ir.ModuleID) error {
	if err := EnsureBazelWorkspaceSetup(b.outputDirectory); err != nil {
		return err
	}

	if err := b.buildModule(moduleID, true /* isRoot */); err != nil {
		return fmt.Errorf("failed to build module %q: %v", moduleID, err)
	}

	if err := GenerateBuildFile(b.outputDirectory, b.targets); err != nil {
		return fmt.Errorf("failed to generate BUILD file: %v", err)
	}

	targetName := sanitizeTargetName(moduleID.Name)
	if err := RunBazelBuild(b.outputDirectory, targetName); err != nil {
		return fmt.Errorf("failed to build module %q with Bazel: %v", moduleID, err)
	}

	// Copy binary from bazel-bin to output location
	bazelBinPath := path.Join(b.outputDirectory, "bazel-bin", targetName)
	outputPath := path.Join(b.outputDirectory, moduleID.Name)
	if err := copyFile(bazelBinPath, outputPath); err != nil {
		return fmt.Errorf("failed to copy built binary to %s: %v", outputPath, err)
	}

	return nil
}

func (b *Builder) buildModule(moduleID ir.ModuleID, isRoot bool) error {
	moduleIDNoPos := moduleID
	moduleIDNoPos.Pos = ir.Pos{}

	b.mutex.Lock()
	if err, ok := b.builtModules[moduleIDNoPos]; ok {
		b.mutex.Unlock()
		if err != nil {
			glog.V(1).Infof("Module %q already failed to build: %v", moduleID, err)
		} else {
			glog.V(1).Infof("Already built module %q", moduleID)
		}
		return err
	}
	// Mark as building to detect cycles.
	b.builtModules[moduleIDNoPos] = fmt.Errorf("cycle detected: module %q is already building", moduleID)
	b.mutex.Unlock()

	err := b.buildModuleImpl(moduleID, isRoot)

	b.mutex.Lock()
	b.builtModules[moduleIDNoPos] = err
	b.mutex.Unlock()

	return err
}

func (b *Builder) buildModuleImpl(moduleID ir.ModuleID, isRoot bool) error {
	glog.V(1).Infof("Building module %q", moduleID)

	moduleQuery, err := b.querier.QueryModule(moduleID)
	if err != nil {
		return fmt.Errorf("failed to query module %q: %v", moduleID, err)
	}

	var srcs []string
	var hdrs []string
	var deps []string

	baseFilename := b.querier.BaseSourceFilename(moduleID)
	baseOutputBasename := toBaseOutputFilename(moduleID)

	// Base BPL generated files are always expected.
	srcs = append(srcs, baseOutputBasename+".cc")
	hdrs = append(hdrs, baseOutputBasename+".h", baseOutputBasename+"_private.h")

	moduleFlags := make([]string, 0, len(moduleQuery.Flags))
	for _, flag := range moduleQuery.Flags {
		moduleFlags = append(moduleFlags, flag.Value)
	}

	// Build imports sequentially
	for _, id := range moduleQuery.Imports {
		if err := b.buildModule(id, false /* isRoot */); err != nil {
			return err
		}
		deps = append(deps, ":"+sanitizeTargetName(id.Name))
	}

	// Compile base BPL
	_, err = b.compileBPL(baseFilename.Value, baseOutputBasename)
	if err != nil {
		return err
	}

	// Compile impls in parallel
	var wg sync.WaitGroup
	var errs []error
	var errsMu sync.Mutex
	var srcsMu sync.Mutex
	var hdrsMu sync.Mutex

	for _, relativeImplFilename := range moduleQuery.Impls {
		implFilename := b.querier.ImplSourceFilename(baseFilename, relativeImplFilename)
		ext := path.Ext(implFilename.Value)

		wg.Add(1)
		go func(implFilename ir.Filename, ext string) {
			defer wg.Done()

			if ext == ".bpl" {
				implOutputBasename := toImplOutputFilename(moduleID, implFilename)
				_, err := b.compileBPL(implFilename.Value, implOutputBasename)
				if err != nil {
					errsMu.Lock()
					errs = append(errs, err)
					errsMu.Unlock()
					return
				}
				srcsMu.Lock()
				srcs = append(srcs, implOutputBasename+".cc")
				srcsMu.Unlock()
			} else if ext == ".cc" || ext == ".cpp" {
				srcsMu.Lock()
				srcs = append(srcs, implFilename.Value)
				srcsMu.Unlock()

				dst := path.Join(b.outputDirectory, implFilename.Value)
				if err := copyFile(implFilename.Value, dst); err != nil {
					errsMu.Lock()
					errs = append(errs, fmt.Errorf("failed to copy %s to %s: %v", implFilename.Value, dst, err))
					errsMu.Unlock()
				}
			} else if ext == ".h" {
				hdrsMu.Lock()
				hdrs = append(hdrs, implFilename.Value)
				hdrsMu.Unlock()

				dst := path.Join(b.outputDirectory, implFilename.Value)
				if err := copyFile(implFilename.Value, dst); err != nil {
					errsMu.Lock()
					errs = append(errs, fmt.Errorf("failed to copy %s to %s: %v", implFilename.Value, dst, err))
					errsMu.Unlock()
				}
			}
		}(implFilename, ext)
	}
	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	// Add Bazel Target
	targetName := sanitizeTargetName(moduleID.Name)
	targetType := "cc_library"
	if isRoot {
		targetType = "cc_binary"
		srcs = append(srcs, hdrs...)
		hdrs = nil
	}

	copts := append([]string{"-std=c++17"}, moduleFlags...)

	target := BazelTarget{
		Type:  targetType,
		Name:  targetName,
		Srcs:  srcs,
		Hdrs:  hdrs,
		Deps:  deps,
		Copts: copts,
	}
	b.addTarget(target)

	return nil
}

func (b *Builder) compileBPL(inputFilename, outputBasename string) (string, error) {
	if path.Ext(inputFilename) != ".bpl" {
		return "", fmt.Errorf("expected file with extension '.bpl'; got %q", inputFilename)
	}

	outputFilename := toOutputFilename(inputFilename, b.outputDirectory, outputBasename)

	if err := os.MkdirAll(path.Dir(outputFilename), 0750); err != nil {
		return "", err
	}

	glog.V(1).Infof("Compiling %q to %q", inputFilename, outputFilename)

	if err := comp.CompileBPL(b.querier, inputFilename, outputFilename); err != nil {
		return "", fmt.Errorf("failed to compile %q to %q: %v", inputFilename, outputFilename, err)
	}

	return outputFilename, nil
}

func (b *Builder) addTarget(target BazelTarget) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.targets = append(b.targets, target)
}

func copyFile(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if err := os.MkdirAll(path.Dir(dst), 0750); err != nil {
		return err
	}

	// Remove destination file if it exists to avoid permission denied on read-only files.
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing destination file %s: %v", dst, err)
	}

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcStat.Mode())
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}
