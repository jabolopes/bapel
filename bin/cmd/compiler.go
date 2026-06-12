package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
)

func main() {
	outputFilename := flag.String("o", "", "File to write the C++ output to.")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: compiler [-o output_file] <input_file>")
		os.Exit(1)
	}
	inputFilename := args[0]

	out := *outputFilename
	if out == "" {
		out = parse.ReplaceExtension(inputFilename, ".h")
	}

	querier, err := query.New()
	if err != nil {
		log.Fatalf("failed to create querier: %v", err)
	}

	err = comp.CompileBPLDirect(querier, inputFilename, out)
	if err != nil {
		log.Fatalf("failed to compile: %v", err)
	}
}
