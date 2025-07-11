package query

import (
	"fmt"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
)

type SourceFileQuery struct {
	Imports []ast.ModuleID
	Impls   []ast.Filename
	Decls   []ir.IrDecl
}

func (q SourceFileQuery) Format(f fmt.State, verb rune) {
	if len(q.Imports) > 0 {
		fmt.Fprintln(f, "imports {")
		for _, moduleID := range q.Imports {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), moduleID)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}")
	}

	if len(q.Impls) > 0 {
		if len(q.Imports) > 0 {
			fmt.Fprintln(f)
		}

		fmt.Fprintln(f, "impls {")
		for _, moduleID := range q.Impls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), moduleID)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}")
	}

	if len(q.Decls) > 0 {
		if len(q.Impls) > 0 {
			fmt.Fprintln(f)
		}

		fmt.Fprintln(f, "decls {")
		for _, decl := range q.Decls {
			fmt.Fprint(f, "  ")
			fmt.Fprintf(f, fmt.FormatString(f, 's'), decl)
			fmt.Fprintln(f)
		}
		fmt.Fprint(f, "}")
	}
}
