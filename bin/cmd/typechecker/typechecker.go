package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
)

func main() {
	formatFlag := flag.String("format", "flat", "Output format (flat, json, ir)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: typechecker [-format=flat|json|ir] <input_file>")
		os.Exit(1)
	}
	inputFilename := args[0]

	querier, err := comp.NewQuerier()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create querier: %v\n", err)
		os.Exit(1)
	}

	unit, err := comp.TypecheckSourceFile(querier, comp.TypecheckOptions{}, inputFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to typecheck %q: %v\n", inputFilename, err)
		os.Exit(1)
	}

	switch *formatFlag {
	case "json":
		jsonData, err := json.MarshalIndent(unit, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal IrUnit to JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))

	case "ir":
		fmt.Printf("%s\n", unit)

	case "flat":
		fmt.Printf("MODULE %s\n", unit.ModuleID)
		if unit.Case == ir.BaseUnit {
			fmt.Println("CASE base")
		} else {
			fmt.Println("CASE impl")
		}

		for _, imp := range unit.Imports {
			fmt.Printf("IMPORT %s\n", imp.ModuleID)
		}
		for _, impl := range unit.Impls {
			fmt.Printf("IMPL %s\n", impl.RelativeFilename.Value)
		}
		for _, decl := range unit.Decls {
			s := fmt.Sprintf("%s", decl)
			fmt.Printf("DECL %s\n", strings.ReplaceAll(s, "\n", "\\n"))
		}
		for _, traitImpl := range unit.TraitImpls {
			s := fmt.Sprintf("%s", traitImpl)
			fmt.Printf("TRAIT_IMPL %s\n", strings.ReplaceAll(s, "\n", "\\n"))
		}
		for _, fn := range unit.Functions {
			s := fmt.Sprintf("%s", fn)
			fmt.Printf("FUNC %s\n", strings.ReplaceAll(s, "\n", "\\n"))
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown format %q (expected flat, json, or ir)\n", *formatFlag)
		os.Exit(1)
	}
}
