package ir

import (
	"fmt"
	"io"
	"strings"
)

type Indent struct {
	f       io.Writer
	level   int
	spaces  string
	newline bool
}

func (i *Indent) Inc() *Indent {
	i.level++
	i.spaces = strings.Repeat(" ", 2*i.level)
	return i
}

func (i *Indent) Dec() *Indent {
	i.level--
	i.spaces = strings.Repeat(" ", 2*i.level)
	return i
}

func (i *Indent) Println(a ...any) *Indent {
	if i.level > 0 && i.newline {
		fmt.Fprint(i.f, i.spaces)
	}
	fmt.Fprintln(i.f, a...)
	i.newline = true
	return i
}

func (i *Indent) Print(a ...any) *Indent {
	if i.level > 0 && i.newline {
		fmt.Fprint(i.f, i.spaces)
	}
	fmt.Fprint(i.f, a...)
	i.newline = false
	return i
}

func (i *Indent) Printf(format string, a ...any) *Indent {
	if i.level > 0 && i.newline {
		fmt.Fprint(i.f, i.spaces)
	}
	fmt.Fprintf(i.f, format, a...)
	i.newline = strings.HasSuffix(format, "\n")
	return i
}

func NewIndent(f fmt.State) *Indent {
	width := 0
	if w, ok := f.Width(); ok {
		width = w
	}
	return &Indent{f, width, strings.Repeat(" ", 2*width), true /* newline */}
}
