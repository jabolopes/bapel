package ir

import "fmt"

type Pos struct {
	Filename     string
	BeginLineNum int
	EndLineNum   int
	Line         string
}

func (p Pos) String() string {
	if p.BeginLineNum == p.EndLineNum {
		return fmt.Sprintf("In %q in line %d: %s", p.Filename, p.BeginLineNum, p.Line)
	}

	return fmt.Sprintf("In %q in lines %d-%d: %s", p.Filename, p.BeginLineNum, p.EndLineNum, p.Line)
}

// TODO: Remove space after colon to put filename and linenum together.
func (p Pos) Format(f fmt.State, verb rune) {
	if commentify := f.Flag('+'); commentify {
		if p.BeginLineNum == p.EndLineNum {
			fmt.Fprintf(f, "/* %s: %d */", p.Filename, p.BeginLineNum)
		} else {
			fmt.Fprintf(f, "/* %s: %d-%d */", p.Filename, p.BeginLineNum, p.EndLineNum)
		}
	} else {
		fmt.Fprint(f, p.String())
	}
}
