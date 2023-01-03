package bin2txt

import (
	"encoding/gob"
	"os"

	"github.com/jabolopes/bapel/ir"
)

func Disassemble(inputFile *os.File, outputFile *os.File) error {
	var program ir.IrProgram
	if err := gob.NewDecoder(inputFile).Decode(&program); err != nil {
		return err
	}

	disassembler := newDisassembler(NewByteArrayDecoder(program.Data), outputFile)
	disassembler.printf("Symbols:\n")

	symbols := program.Header.Symbols
	{
		disassembler.incIndentation()
		for _, symbol := range symbols {
			disassembler.printf("  %s() offset=%d\n", symbol.Id, symbol.Offset)
		}
		disassembler.decIndentation()
	}

	disassembler.printf("Data:\n")

	table := newBindTable()
	currentFunction := 0
	{
		disassembler.incIndentation()

		for disassembler.dec().Len() > 0 {
			if currentFunction < len(symbols) {
				if symbols[currentFunction].Offset < uint64(disassembler.dec().Offset()) {
					currentFunction++
				}
			}

			if currentFunction < len(symbols) {
				if symbol := symbols[currentFunction]; symbol.Offset == uint64(disassembler.dec().Offset()) {
					disassembler.decIndentation()
					disassembler.printf("%s:\n", symbol.Id)
					disassembler.incIndentation()
					disassembler.printf("enterSize=%d\n", disassembler.dec().GetI16())
					currentFunction++
					continue
				}
			}

			opcode, err := disassembler.dec().GetOpCode()
			if err != nil {
				return err
			}

			if err := table.ops[opcode](disassembler); err != nil {
				return err
			}
		}

		disassembler.decIndentation()
	}

	return nil
}
