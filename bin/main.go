package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/jabolopes/bapel/build"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/lex"
	"github.com/jabolopes/bapel/parse"
	"github.com/jabolopes/bapel/query"
)

func cmdLex(args []string) error {
	var input io.Reader
	switch len(args) {
	case 0:
		input = os.Stdin
	case 1:
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer file.Close()
		input = file
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	lex := lex.New(input)

	line := 0
	for {
		token, ok := lex.NextToken()
		if !ok {
			break
		}

		if line != token.LineNum {
			line = token.LineNum
			fmt.Printf("LINE %d:\n", line)
		}

		if len(token.Value) > 0 {
			fmt.Printf("TOKEN: %s\n", token.Value)
		} else {
			fmt.Printf("TOKEN: %v\n", token)
		}
	}

	return lex.ScanErr()
}

func cmdParse(args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected the input to parse as first argument. The input can be a module ID (e.g., 'main') or a source file (e.g., 'main.bpl' or 'main_impl.cc'")
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	if path.Base(inputFilename) == "workspace.bpl" {
		workspace, err := parse.ParseWorkspace(inputFilename)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", workspace)
		return nil
	}

	sourceFile, err := parse.ParseSourceFile(inputFilename)
	if err != nil {
		return err
	}

	fmt.Println(sourceFile.Imports)
	fmt.Println(sourceFile.Impls)
	for _, source := range sourceFile.Body {
		fmt.Println(source)
	}

	return nil
}

func cmdCc(outputFilename string, args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected source file to build as first argument")
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	if len(outputFilename) == 0 {
		outputFilename = parse.ReplaceExtension(inputFilename, ".ccm")
	}

	querier, err := query.New()
	if err != nil {
		return err
	}

	return comp.CompileBPLToCCM(querier, inputFilename, outputFilename)
}

func cmdBuild(args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected the input to build as first argument. The input can be a module ID (e.g., 'main') or a base source file (e.g., 'main.bpl'")
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	querier, err := query.New()
	if err != nil {
		return err
	}
	builder := build.NewBuilder(querier)

	var moduleID ir.ModuleID
	if len(path.Ext(inputFilename)) > 0 {
		moduleID = ir.NewModuleIDFromFilename(inputFilename)
	} else {
		// Query the module, recursing into the `impls` section.
		moduleID = ir.NewModuleID(inputFilename, ir.Pos{})
	}

	return builder.Build(moduleID)
}

func cmdQuery(queryStr string, args []string) error {
	var input string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected the input to query as first argument. The input can be a module ID (e.g., 'main') or a source file (e.g., 'main.bpl' or 'main_impl.cc'")
	case 1:
		input = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	querier, err := query.New()
	if err != nil {
		return err
	}

	isSourceFilename := strings.HasPrefix(input, "./")
	switch {
	case queryStr == "resolve" && isSourceFilename:
		sourceFile, err := parse.ParseSourceFile(input)
		if err != nil {
			return err
		}

		unit, err := comp.ResolveSourceFile(querier, sourceFile)
		if err != nil {
			return err
		}

		fmt.Printf("%v\n", unit)

	case len(queryStr) == 0 && isSourceFilename:
		// Query the source file only, without recursing into the `impls` section.
		result, err := query.QuerySourceFile(input)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", result)

	case len(queryStr) == 0:
		querier, err := query.New()
		if err != nil {
			return err
		}

		// Query the module, recursing into the `impls` section.
		moduleID := ir.NewModuleID(input, ir.Pos{})

		result, err := querier.QueryModule(moduleID)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", result)

	default:
		return fmt.Errorf("unknown combination of query %q and input %q", queryStr, input)
	}

	return nil
}

func run() error {
	// Uncomment to run the profiler.
	//
	// defer profile.Start().Stop()

	lexCmd := flag.NewFlagSet("lex", flag.ExitOnError)
	parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)

	ccCmd := flag.NewFlagSet("cc", flag.ExitOnError)
	ccOutputFilename := ccCmd.String("o", "", "File to write the C++ output to.")

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)

	queryCmd := flag.NewFlagSet("query", flag.ExitOnError)
	queryStr := queryCmd.String("q", "", "Query to ask")

	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("expected subcommand, e.g., 'lex', 'parse', 'cc', etc")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "lex":
		lexCmd.Parse(args[1:])
		return cmdLex(lexCmd.Args())
	case "parse":
		parseCmd.Parse(args[1:])
		return cmdParse(parseCmd.Args())
	case "cc":
		ccCmd.Parse(args[1:])
		return cmdCc(*ccOutputFilename, ccCmd.Args())
	case "build":
		buildCmd.Parse(args[1:])
		return cmdBuild(buildCmd.Args())
	case "query":
		queryCmd.Parse(args[1:])
		return cmdQuery(*queryStr, queryCmd.Args())
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
