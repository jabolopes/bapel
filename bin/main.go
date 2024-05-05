package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jabolopes/bapel/bin2txt"
	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/comp"
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

func cmdLex() error {
	parser := bplparser.NewParser()
	parser.Open(os.Stdin)
	for parser.Scan() {
		fmt.Printf("LINE: %q\n", parser.Line())

		words := parser.Words()
		for _, word := range words {
			fmt.Printf("WORD: %q\n", word)
		}
	}

	return nil
}

func cmdParse() error {
	sources, err := bplparser.ParseFile(os.Stdin)
	if err != nil {
		return err
	}

	for _, source := range sources {
		fmt.Println(source)
	}

	return nil
}

func cmdCpp() error {
	var outputFilename string
	flag.StringVar(&outputFilename, "o", "a.bpl.cpp", "File to write the C++ output to.")

	flag.Parse()

	outputFile := os.Stdout
	if len(outputFilename) > 0 {
		var err error
		if outputFile, err = os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return err
		}
		defer closeFile(outputFilename, &outputFile)
	}

	if err := comp.CompileFile(os.Stdin, outputFile); err != nil {
		return err
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

func cmdBin2Txt() error {
	var inputFilename string
	flag.StringVar(&inputFilename, "input", "", "File to read binary assemble file from. If empty, reads from standard input.")

	var outputFilename string
	flag.StringVar(&outputFilename, "output", "", "File to write disassembled file to. If empty, writes to standard output.")

	flag.Parse()

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

func cmdQuery() error {
	var inputFilename string
	flag.StringVar(&inputFilename, "input", "", "Bapel source file to query, e.g., 'myfile.bpl'.")

	flag.Parse()

	inputFile := os.Stdin
	if len(inputFilename) > 0 {
		var err error
		if inputFile, err = os.Open(inputFilename); err != nil {
			return err
		}
		defer closeFile(inputFilename, &inputFile)
	}

	decls, err := query.QueryExports(inputFile)
	if err != nil {
		return err
	}

	for _, decl := range decls {
		fmt.Printf("%s\n", decl)
	}

	return nil
}

func run(command string) error {
	defer profile.Start().Stop()

	if command == "" {
		command = "run"
	}

	switch command {
	case "lex":
		return cmdLex()
	case "parse":
		return cmdParse()
	case "cpp":
		return cmdCpp()
	case "bin2txt":
		return cmdBin2Txt()
	case "query":
		return cmdQuery()
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
		os.Args = append(os.Args[0:1], os.Args[2:]...)
	}

	if err := run(command); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
