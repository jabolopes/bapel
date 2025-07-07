package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/build"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/lexer"
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

	lexer := lexer.New(input)

	line := 0
	for {
		token, ok := lexer.NextToken()
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

	return lexer.ScanErr()
}

func cmdParse(args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected the module to query as first argument. The module can be a module ID (e.g., 'main') or a module file (e.g., 'main.bpl' or 'main_impl.cc'")
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	if path.Base(inputFilename) == "workspace.bpl" {
		workspace, err := bplparser2.ParseWorkspace(inputFilename)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", workspace)
		return nil
	}

	module, err := bplparser2.ParseModuleFile(inputFilename)
	if err != nil {
		return err
	}

	fmt.Println(module.Imports)
	fmt.Println(module.Impls)
	for _, source := range module.Body {
		fmt.Println(source)
	}

	return nil
}

func cmdCc(outputFilename string, args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected module file to build as first argument")
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	if len(outputFilename) == 0 {
		outputFilename = bplparser2.ReplaceExtension(inputFilename, ".ccm")
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
		return fmt.Errorf("expected the module to build as first argument. The module can be a module ID (e.g., 'main') or a module base file (e.g., 'main.bpl'")
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

	var moduleID ast.ModuleID
	if len(path.Ext(inputFilename)) > 0 {
		moduleID = ast.NewModuleIDFromFilename(inputFilename)
	} else {
		// Query the module, recursing into the `impls` section.
		moduleID = ast.NewModuleID(inputFilename, ir.Pos{})
	}

	return builder.Build(moduleID)
}

func cmdQuery(args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected the module to query as first argument. The module can be a module ID (e.g., 'main') or a module file (e.g., 'main.bpl' or 'main_impl.cc'")
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	if len(path.Ext(inputFilename)) > 0 {
		// Query the module file only, without recursing into the `impls` section.
		decls, err := query.QueryFileDecls(inputFilename)
		if err != nil {
			return err
		}

		for _, decl := range decls {
			fmt.Printf("%s\n", decl)
		}

		return nil
	}

	{
		querier, err := query.New()
		if err != nil {
			return err
		}

		// Query the module, recursing into the `impls` section.
		moduleID := ast.NewModuleID(inputFilename, ir.Pos{})

		{
			module, err := querier.QueryModuleMetadata(moduleID)
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", module)
		}

		{
			decls, err := querier.QueryModuleDecls(moduleID)
			if err != nil {
				return err
			}

			for _, decl := range decls {
				fmt.Printf("%s\n", decl)
			}
		}
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
		return cmdQuery(queryCmd.Args())
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
