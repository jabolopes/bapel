package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/jabolopes/bapel/bin2txt"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/lexer"
	"github.com/jabolopes/bapel/query"
	"github.com/pkg/profile"
)

func closeFile(filename string, file **os.File) {
	if *file == nil {
		return
	}

	if err := (*file).Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close %q: %v\n", filename, err)
	}
}

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
	var reader io.Reader
	switch len(args) {
	case 0:
		reader = os.Stdin
	case 1:
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	module, err := bplparser2.ParseFile("stdin", reader)
	if err != nil {
		return err
	}

	fmt.Println(module.Imports)
	fmt.Println(module.Exports)
	fmt.Println(module.Impls)
	for _, source := range module.Body {
		fmt.Println(source)
	}

	return nil
}

func cmdCc(outputFilename string, args []string) error {
	inputFilename := "stdin"
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
		inputFilename = file.Name()
		input = file
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	if len(outputFilename) == 0 {
		outputFilename = bplparser2.ReplaceExtension(inputFilename, ".cc")
	}

	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer closeFile(outputFilename, &outputFile)

	if err := comp.CompileModule(inputFilename, input, outputFile); err != nil {
		return err
	}

	{
		var file *os.File
		file, outputFile = outputFile, nil
		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func cmdBuild(args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		return fmt.Errorf("expected module to build as first argument")
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	builder := comp.NewBuilder()
	return builder.Build(inputFilename)
}

func cmdBin2Txt(inputFilename, outputFilename string, args []string) error {
	inputFile := os.Stdin
	if len(inputFilename) > 0 {
		var err error
		if inputFile, err = os.Open(inputFilename); err != nil {
			return err
		}
		defer closeFile(inputFilename, &inputFile)
	}

	outputFile := os.Stdout
	if len(outputFilename) > 0 {
		var err error
		if outputFile, err = os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			return err
		}
		defer closeFile(inputFilename, &outputFile)
	}

	if err := bin2txt.Disassemble(inputFile, outputFile); err != nil {
		return err
	}

	if inputFile != os.Stdin {
		var file *os.File
		file, inputFile = inputFile, nil
		if err := file.Close(); err != nil {
			return err
		}
	}

	if outputFile != os.Stdout {
		var file *os.File
		file, outputFile = outputFile, nil
		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

// TODO: Delete the stdin option since cmdQuery now depends on the
// filename to determine whether the target is a module, or a base
// module file, or an implementation file, etc.
func cmdQuery(args []string) error {
	var inputFilename string
	switch len(args) {
	case 0:
		inputFilename = "stdin"
	case 1:
		inputFilename = args[0]
	default:
		return fmt.Errorf("too many arguments %q", strings.Join(args, " "))
	}

	if len(path.Ext(inputFilename)) > 0 {
		// Query the module file without recursing into the `impls` section.
		var input *os.File
		if inputFilename == "stdin" {
			input = os.Stdin
		} else {
			file, err := os.Open(inputFilename)
			if err != nil {
				return err
			}
			defer file.Close()

			input = file
		}

		decls, err := query.QueryDecls(inputFilename, input)
		if err != nil {
			return err
		}

		for _, decl := range decls {
			fmt.Printf("%s\n", decl)
		}
	} else {
		// Query the module, recursing into the `impls` section.
		decls, err := query.QueryModuleDecls(inputFilename)
		if err != nil {
			return err
		}

		for _, decl := range decls {
			fmt.Printf("%s\n", decl)
		}
	}

	return nil
}

func run() error {
	defer profile.Start().Stop()

	lexCmd := flag.NewFlagSet("lex", flag.ExitOnError)
	parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)

	ccCmd := flag.NewFlagSet("cc", flag.ExitOnError)
	ccOutputFilename := ccCmd.String("o", "", "File to write the C++ output to.")

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)

	b2tCmd := flag.NewFlagSet("bin2txt", flag.ExitOnError)
	b2tInputFilename := b2tCmd.String("input", "", "File to read binary assemble file from. If empty, reads from standard input.")
	b2tOutputFilename := b2tCmd.String("output", "", "File to write disassembled file to. If empty, writes to standard output.")

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
	case "bin2txt":
		b2tCmd.Parse(args[1:])
		return cmdBin2Txt(*b2tInputFilename, *b2tOutputFilename, b2tCmd.Args())
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
