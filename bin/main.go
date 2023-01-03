package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/jabolopes/bapel/asm"
	"github.com/jabolopes/bapel/bin2txt"
	"github.com/jabolopes/bapel/vm"
)

func closeFile(filename string, file **os.File) {
	if *file == nil {
		return
	}

	if err := (*file).Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close %q: %v\n", filename, err)
	}
}

func cmdAsm() error {
	var outputFilename string
	flag.StringVar(&outputFilename, "o", "a.bpl.asm", "File to write the assembly output to.")

	flag.Parse()

	outputFile := os.Stdout
	if len(outputFilename) > 0 {
		var err error
		if outputFile, err = os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			return err
		}
		defer closeFile(outputFilename, &outputFile)
	}

	program, err := asm.AssembleFile(os.Stdin)
	if err != nil {
		return err
	}

	if err := gob.NewEncoder(outputFile).Encode(program); err != nil {
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

func cmdRun() error {
	program, err := asm.AssembleFile(os.Stdin)
	if err != nil {
		return err
	}

	machine := vm.New(program)
	return machine.Run()
}

func run(command string) error {
	if command == "" {
		command = "run"
	}

	switch command {
	case "asm":
		return cmdAsm()
	case "bin2txt":
		return cmdBin2Txt()
	case "run":
		return cmdRun()
	default:
		return fmt.Errorf("Unknown command %q", command)
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
