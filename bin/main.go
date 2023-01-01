package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jabolopes/bapel/asm"
	"github.com/jabolopes/bapel/vm"
)

var (
	outputFile = flag.String("output", "a.bpl.asm", "File to write the assembly output to.")
)

func run(command, outputFile string) error {
	if command == "" {
		command = "run"
	}

	program, err := asm.AssembleFile(os.Stdin)
	if err != nil {
		return err
	}

	if command == "asm" {
		if len(outputFile) <= 0 {
			// Write assembly to stdout.
			_, err := os.Stdout.Write(program.Data)
			return err
		}

		// Write assemble to file.
		if err := os.WriteFile(outputFile, program.Data, 0644); err != nil {
			return err
		}

		fmt.Printf("Output %s\n", outputFile)
		return nil
	}

	machine := vm.New(program)
	if err := machine.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
		os.Args = append(os.Args[0:1], os.Args[2:]...)
	}

	flag.Parse()

	if err := run(command, *outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
