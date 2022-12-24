package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jabolopes/bapel/asm"
	"github.com/jabolopes/bapel/ssa"
	"github.com/jabolopes/bapel/vm"
)

var (
	useSsa = flag.Bool("ssa", false, "Whether to use SSA syntax.")
)

func run(useSsa bool) error {
	var program vm.OpProgram
	var err error
	if useSsa {
		program, err = ssa.AssembleFile(os.Stdin)
	} else {
		program, err = asm.AssembleFile(os.Stdin)
	}
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
	flag.Parse()

	if err := run(*useSsa); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
