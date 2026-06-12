package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
)

func main() {
	symbolFlag := flag.String("symbol", "SourceFile", "Initial symbol to parse")
	workspaceFlag := flag.Bool("workspace", false, "Parse as a workspace file")
	filenameFlag := flag.String("filename", "", "Filename to use for position info when reading from stdin")
	formatFlag := flag.String("format", "json", "Output format (json, flat)")
	flag.Parse()

	args := flag.Args()
	var input io.Reader
	filename := *filenameFlag

	if len(args) == 0 {
		input = os.Stdin
		if filename == "" {
			filename = "<stdin>"
		}
	} else if len(args) == 1 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
		if filename == "" {
			filename = args[0]
		}
	} else {
		fmt.Fprintln(os.Stderr, "Expected at most one argument")
		os.Exit(1)
	}

	symbol := *symbolFlag
	if *workspaceFlag {
		symbol = "Workspace"
	}

	astData, err := ParseSymbol(symbol, filename, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Validation and adjustments
	if symbol == "SourceFile" {
		sf := astData.(ast.SourceFile)
		sf.Header.Filename = ir.NewFilename(filename, ir.Pos{})
		if validation := ast.ValidateSourceFile(&sf); !validation.OK() {
			fmt.Fprintln(os.Stderr, validation.Err())
			os.Exit(1)
		}
		astData = sf
	} else if symbol == "Workspace" {
		ws := astData.(ast.Workspace)
		if validation := ast.ValidateWorkspace(ws); !validation.OK() {
			fmt.Fprintln(os.Stderr, validation.Err())
			os.Exit(1)
		}
		astData = ws
	}

	if *formatFlag == "flat" {
		if symbol == "SourceFile" {
			sf := astData.(ast.SourceFile)
			for _, imp := range sf.Imports.IDs {
				fmt.Printf("IMPORT %s\n", imp.Name)
			}
			for _, impl := range sf.Impls.Filenames {
				fmt.Printf("IMPL %s\n", impl.Value)
			}
			return
		} else if symbol == "Workspace" {
			ws := astData.(ast.Workspace)
			for _, pkg := range ws.Packages.Packages {
				if pkg.Case == ast.PrefixPackage {
					fmt.Printf("PREFIX %s %s\n", pkg.Prefix.Prefix.Name, pkg.Filename.Value)
				} else if pkg.Case == ast.ModulePackage {
					fmt.Printf("MODULE %s %s\n", pkg.Module.ModuleID.Name, pkg.Filename.Value)
				}
			}
			return
		}
		fmt.Fprintf(os.Stderr, "Flat format not supported for symbol %q\n", symbol)
		os.Exit(1)
	}

	jsonData, err := json.Marshal(astData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal AST to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}
