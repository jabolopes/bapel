package bin2txt

import (
	"fmt"
	"os"
)

type disassembler struct {
	decoder    *ByteArrayDecoder
	outputFile *os.File
}

func (m *disassembler) dec() *ByteArrayDecoder { return m.decoder }
func (m *disassembler) out() *os.File          { return m.outputFile }

func (m *disassembler) printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(m.outputFile, format, a...)
}

func newDisassembler(decoder *ByteArrayDecoder, outputFile *os.File) *disassembler {
	return &disassembler{decoder, outputFile}
}
