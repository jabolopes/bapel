package ir

import "fmt"

type Pos struct {
	Filename     string
	BeginLineNum int
	EndLineNum   int
}

func (p Pos) String() string {
	if p.BeginLineNum == p.EndLineNum {
		return fmt.Sprintf("in %q in line %d", p.Filename, p.BeginLineNum)
	}

	return fmt.Sprintf("in %q in lines %d-%d", p.Filename, p.BeginLineNum, p.EndLineNum)
}

func (p Pos) Format(f fmt.State, verb rune) {
	if commentify := f.Flag('+'); commentify {
		if p.BeginLineNum == p.EndLineNum {
			fmt.Fprintf(f, "/* %s:%d */", p.Filename, p.BeginLineNum)
		} else {
			fmt.Fprintf(f, "/* %s:%d-%d */", p.Filename, p.BeginLineNum, p.EndLineNum)
		}
	} else {
		fmt.Fprint(f, p.String())
	}
}

func NewLinePos(filename string, lineNum int) Pos {
	if lineNum == 0 {
		panic("Invalid line number 0")
	}

	return Pos{filename, lineNum, lineNum}
}

func NewRangePos(filename string, beginLineNum, endLineNum int) Pos {
	if beginLineNum == 0 {
		panic("Invalid begin line number 0")
	}

	if endLineNum == 0 {
		panic("Invalid end line number 0")
	}

	return Pos{filename, beginLineNum, endLineNum}
}
