package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

	sources, err := bplparser2.ParseFile("stdin", reader)
	if err != nil {
		return err
	}

	for _, source := range sources {
		fmt.Println(source)
	}

	return nil
}

func cmdCpp(outputFilename, module string, args []string) error {
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
		outputFilename = comp.ReplaceExtension(inputFilename, ".cpp")
	}

	outputFile := os.Stdout
	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer closeFile(outputFilename, &outputFile)

	if len(module) == 0 {
		if err := comp.CompileModuleFile(inputFilename, input, outputFile); err != nil {
			return err
		}
	} else {
		if err := comp.CompileImplFile(inputFilename, module, input, outputFile); err != nil {
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

func cmdQuery(args []string) error {
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

	decls, err := query.QueryExports(inputFilename, input)
	if err != nil {
		return err
	}

	for _, decl := range decls {
		fmt.Printf("%s\n", decl)
	}

	return nil
}

func run() error {
	defer profile.Start().Stop()

	lexCmd := flag.NewFlagSet("lex", flag.ExitOnError)
	parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)

	cppCmd := flag.NewFlagSet("cpp", flag.ExitOnError)
	cppModule := cppCmd.String("m", "", "Module this implementation source file belongs to, e.g., 'program'")
	cppOutputFilename := cppCmd.String("o", "", "File to write the C++ output to.")

	b2tCmd := flag.NewFlagSet("bin2txt", flag.ExitOnError)
	b2tInputFilename := b2tCmd.String("input", "", "File to read binary assemble file from. If empty, reads from standard input.")
	b2tOutputFilename := b2tCmd.String("output", "", "File to write disassembled file to. If empty, writes to standard output.")

	queryCmd := flag.NewFlagSet("query", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand, e.g., 'lex', 'parse', 'cpp', etc")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "lex":
		lexCmd.Parse(os.Args[2:])
		return cmdLex(lexCmd.Args())
	case "parse":
		parseCmd.Parse(os.Args[2:])
		return cmdParse(parseCmd.Args())
	case "cpp":
		cppCmd.Parse(os.Args[2:])
		return cmdCpp(*cppOutputFilename, *cppModule, cppCmd.Args())
	case "bin2txt":
		b2tCmd.Parse(os.Args[2:])
		return cmdBin2Txt(*b2tInputFilename, *b2tOutputFilename, b2tCmd.Args())
	case "query":
		queryCmd.Parse(os.Args[2:])
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
