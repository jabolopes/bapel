package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jabolopes/bapel/ast"
)

func FindParserBin() (string, error) {
	// Try relative to CWD first
	path := "bootstrap/parser"
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		return filepath.Abs(path)
	}

	// Walk up to find it
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path = filepath.Join(cwd, "bootstrap/parser")
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			return path, nil
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break // reached root
		}
		cwd = parent
	}

	return "", fmt.Errorf("could not find bootstrap/parser binary")
}

func ParseSymbol[T any](symbol string, filename string, reader io.Reader) (T, error) {
	var t T

	parserBin, err := FindParserBin()
	if err != nil {
		return t, err
	}

	cmd := exec.Command(parserBin, "--symbol", symbol, "--filename", filename)
	cmd.Stdin = reader
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return t, errors.New(strings.TrimSpace(stderr.String()))
		}
		return t, fmt.Errorf("parser failed: %v", err)
	}

	if err := json.Unmarshal(stdout.Bytes(), &t); err != nil {
		return t, fmt.Errorf("failed to unmarshal AST: %v", err)
	}

	return t, nil
}

func ParseSourceFile(inputFilename string) (ast.SourceFile, error) {
	file, err := os.Open(inputFilename)
	if err != nil {
		return ast.SourceFile{}, err
	}
	defer file.Close()
	return ParseSymbol[ast.SourceFile]("SourceFile", inputFilename, file)
}

func ParseWorkspace(inputFilename string) (ast.Workspace, error) {
	file, err := os.Open(inputFilename)
	if err != nil {
		return ast.Workspace{}, err
	}
	defer file.Close()
	return ParseSymbol[ast.Workspace]("Workspace", inputFilename, file)
}
