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
