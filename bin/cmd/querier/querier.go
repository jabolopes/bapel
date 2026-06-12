package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
)

func main() {
	queryStr := flag.String("q", "", "Query to ask")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: querier [-q query] <input>")
		os.Exit(1)
	}
	input := args[0]

	querier, err := query.New()
	if err != nil {
		log.Fatalf("failed to create querier: %v", err)
	}

	isSourceFilename := strings.HasPrefix(input, "./") || strings.HasSuffix(input, ".bpl")
	// Note: modified isSourceFilename check to also check for .bpl suffix to be more robust.

	switch {
	case *queryStr == "resolve" && isSourceFilename:
		sourceFile, err := parse.ParseSourceFile(input)
		if err != nil {
			log.Fatalf("failed to parse source file: %v", err)
		}

		unit, err := comp.ResolveSourceFile(querier, sourceFile)
		if err != nil {
			log.Fatalf("failed to resolve: %v", err)
		}

		fmt.Printf("%v\n", unit)

	case (*queryStr == "infer" || *queryStr == "typecheck") && isSourceFilename:
		options := comp.TypecheckOptions{}
		if *queryStr == "infer" {
			options.SkipTermTypechecker = true
		}

		unit, err := comp.TypecheckSourceFile(querier, options, input)
		if err != nil {
			log.Fatalf("failed to typecheck: %v", err)
		}

		fmt.Printf("%v\n", unit)

	case len(*queryStr) == 0 && isSourceFilename:
		result, err := query.QuerySourceFile(input)
		if err != nil {
			log.Fatalf("failed to query source file: %v", err)
		}

		fmt.Printf("%s\n", result)

	case len(*queryStr) == 0:
		moduleID := ir.NewModuleID(input, ir.Pos{})

		result, err := querier.QueryModule(moduleID)
		if err != nil {
			log.Fatalf("failed to query module: %v", err)
		}

		fmt.Printf("%s\n", result)

	default:
		log.Fatalf("unknown combination of query %q and input %q", *queryStr, input)
	}
}
