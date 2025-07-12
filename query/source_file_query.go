package query

import (
	"fmt"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
)

type SourceFileQuery struct {
	Imports []ast.ModuleID
	Impls   []ast.Filename
	Flags   []ast.Filename
	Decls   []ir.IrDecl
}

func (q SourceFileQuery) Format(f fmt.State, verb rune) {
	empty := true
	newline := func() {
		if empty {
			empty = false
		} else {
			fmt.Fprintln(f)
		}
	}

	if len(q.Imports) > 0 {
		newline()

		fmt.Fprintln(f, "imports {")
		for _, moduleID := range q.Imports {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), moduleID)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}")
	}

	if len(q.Impls) > 0 {
		newline()

		fmt.Fprintln(f, "impls {")
		for _, impl := range q.Impls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 'q'), impl)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}")
	}

	if len(q.Flags) > 0 {
		newline()

		fmt.Fprintln(f, "flags {")
		for _, flag := range q.Flags {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 'q'), flag)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}")
	}

	if len(q.Decls) > 0 {
		newline()

		fmt.Fprintln(f, "decls {")
		for _, decl := range q.Decls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), decl)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}")
	}
}
