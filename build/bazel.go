package build

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/golang/glog"
)

// sanitizeTargetName replaces '.' and '/' with '_' to make it a valid Bazel target name.
func sanitizeTargetName(name string) string {
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "/", "_", -1)
	return name
}

// BazelTarget represents a cc_library or cc_binary target in a BUILD file.
type BazelTarget struct {
	Type  string   // "cc_library" or "cc_binary"
	Name  string   // Target name
	Srcs  []string // Source files relative to out/
	Hdrs  []string // Header files relative to out/
	Deps  []string // Dependency targets
	Copts []string // Compiler options
}

// CopyFile copies a file from src to dst. It creates destination directories if they don't exist.
func CopyFile(src, dst string) error {
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

// EnsureBazelWorkspaceSetup ensures that WORKSPACE and MODULE.bazel exist and are configured.
func EnsureBazelWorkspaceSetup(outputDirectory string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDirectory, 0750); err != nil {
		return fmt.Errorf("failed to create output directory %s: %v", outputDirectory, err)
	}

	// Ensure WORKSPACE exists and has workspace name
	workspacePath := path.Join(outputDirectory, "WORKSPACE")
	file, err := os.OpenFile(workspacePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create WORKSPACE file: %v", err)
	}
	
	workspaceContent := `workspace(name = "bapel_out")
`
	_, err = file.WriteString(workspaceContent)
	closeErr := file.Close()
	if err != nil {
		return fmt.Errorf("failed to write to WORKSPACE: %v", err)
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close WORKSPACE: %v", closeErr)
	}

	// Ensure MODULE.bazel exists and has rules_cc dependency
	modulePath := path.Join(outputDirectory, "MODULE.bazel")
	mfile, err := os.OpenFile(modulePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create MODULE.bazel file: %v", err)
	}
	defer mfile.Close()

	content := `module(name = "bapel_out")
bazel_dep(name = "rules_cc", version = "0.2.17")
`
	if _, err := mfile.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to MODULE.bazel: %v", err)
	}
	return nil
}

const buildTemplate = `load("@rules_cc//cc:defs.bzl", "cc_binary", "cc_library")

{{range .}}
{{.Type}}(
    name = "{{.Name}}",
{{if .Srcs}}    srcs = [
{{range .Srcs}}        "{{.}}",
{{end}}    ],
{{end}}{{if .Hdrs}}    hdrs = [
{{range .Hdrs}}        "{{.}}",
{{end}}    ],
{{end}}{{if .Copts}}    copts = [
{{range .Copts}}        "{{.}}",
{{end}}    ],
{{end}}{{if .Deps}}    deps = [
{{range .Deps}}        "{{.}}",
{{end}}    ],
{{end}})
{{end}}`

// GenerateBuildFile generates a BUILD file in outputDirectory with the given targets.
func GenerateBuildFile(outputDirectory string, targets []BazelTarget) error {
	tmpl, err := template.New("BUILD").Parse(buildTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse BUILD template: %v", err)
	}

	buildPath := path.Join(outputDirectory, "BUILD")
	file, err := os.OpenFile(buildPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open BUILD file for writing: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, targets); err != nil {
		return fmt.Errorf("failed to execute BUILD template: %v", err)
	}

	return nil
}

// RunBazelBuild runs bazel build //:<targetName> inside outputDirectory.
func RunBazelBuild(outputDirectory, targetName string) error {
	cmd := exec.Command("bazel", "build", fmt.Sprintf("//:%s", targetName))
	cmd.Dir = outputDirectory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CC=clang", "CXX=clang++")

	glog.V(1).Infof("Running bazel build //:%s inside %s", targetName, outputDirectory)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("bazel build failed: %v", err)
	}

	return nil
}

