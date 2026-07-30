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
	outputFilename := flag.String("o", "", "File to write C++ output to (when compiling)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: typechecker [-o output_file] [-format=flat|json|ir] <input_file>")
		os.Exit(1)
	}
	inputFilename := args[0]

	querier, err := comp.NewQuerier()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create querier: %v\n", err)
		os.Exit(1)
	}

	if *outputFilename != "" {
		if err := comp.CompileBPLDirect(querier, inputFilename, *outputFilename); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to compile %q: %v\n", inputFilename, err)
			os.Exit(1)
		}
		return
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
			exportStr := "0"
			if decl.Export {
				exportStr = "1"
			}
			fmt.Printf("DECL_DEF %s %s %s\n", exportStr, decl.ID(), strings.ReplaceAll(s, "\n", "\\n"))
		}
		for _, traitImpl := range unit.TraitImpls {
			s := fmt.Sprintf("%s", traitImpl)
			fmt.Printf("TRAIT_IMPL %s\n", strings.ReplaceAll(s, "\n", "\\n"))
			traitTypeStr := strings.ReplaceAll(traitImpl.TraitType.String(), "\n", "\\n")
			typeNameStr := strings.ReplaceAll(traitImpl.TypeName.String(), "\n", "\\n")
			fmt.Printf("TRAIT_DEF %s %s %s\n", traitTypeStr, typeNameStr, strings.ReplaceAll(s, "\n", "\\n"))
		}
		for _, fn := range unit.Functions {
			s := fmt.Sprintf("%s", fn)
			fmt.Printf("FUNC %s\n", strings.ReplaceAll(s, "\n", "\\n"))
			exportStr := "0"
			if fn.Export {
				exportStr = "1"
			}
			bodyStr := strings.ReplaceAll(fmt.Sprintf("%s", fn.Body), "\n", "\\n")
			retTypeStr := strings.ReplaceAll(fn.RetType.String(), "\n", "\\n")
			fmt.Printf("FUNC_DEF %s %s %s %s\n", exportStr, fn.ID, retTypeStr, bodyStr)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown format %q (expected flat, json, or ir)\n", *formatFlag)
		os.Exit(1)
	}
}
