package ir

import "fmt"

type Pos struct {
	Filename string
	LineNum  int
	Line     string
}

func (p Pos) String() string {
	return fmt.Sprintf("In %q in line %d: %s", p.Filename, p.LineNum, p.Line)
}
