package bin2txt

import (
	"fmt"
	"os"
	"strings"
)

type disassembler struct {
	decoder     *ByteArrayDecoder
	outputFile  *os.File
	indentation int
}

func (d *disassembler) dec() *ByteArrayDecoder { return d.decoder }
func (d *disassembler) out() *os.File          { return d.outputFile }

func (d *disassembler) incIndentation() *disassembler {
	d.indentation++
	return d
}

func (d *disassembler) decIndentation() *disassembler {
	d.indentation--
	return d
}

func (d *disassembler) printf(format string, a ...any) (n int, err error) {
	if d.indentation > 0 {
		format = strings.Repeat(" ", d.indentation*2) + format
	}
	return fmt.Fprintf(d.outputFile, format, a...)
}

func newDisassembler(decoder *ByteArrayDecoder, outputFile *os.File) *disassembler {
	return &disassembler{decoder, outputFile, 0 /* indentation */}
}
