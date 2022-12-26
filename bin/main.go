package main

import (
	"fmt"
	"os"

	"github.com/jabolopes/bapel/ssa"
	"github.com/jabolopes/bapel/vm"
)

func run() error {
	program, err := ssa.AssembleFile(os.Stdin)
	if err != nil {
		return err
	}

	machine := vm.New(program)
	if err := machine.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
