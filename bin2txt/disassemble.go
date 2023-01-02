package bin2txt

import (
	"io"
	"os"
)

func Disassemble(inputFile *os.File, outputFile *os.File) error {
	data, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	table := newBindTable()

	for disassembler := newDisassembler(NewByteArrayDecoder(data), outputFile); disassembler.dec().Len() > 0; {
		opcode, err := disassembler.dec().GetOpCode()
		if err != nil {
			return err
		}

		if err := table.ops[opcode](disassembler); err != nil {
			return err
		}
	}

	return nil
}
