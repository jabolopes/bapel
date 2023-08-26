package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

// ParseCallAssign parses call and assignment.
//
// Note that a call is an assignment without the '<-' and without any return
// values.
func ParseCallAssign(args []string) (ir.IrTerm, []string, error) {
	orig := args

	var rets []string
	for ; len(args) > 0; args = args[1:] {
		if args[0] == "<-" {
			args = args[1:]

			if len(args) == 0 {
				return ir.IrTerm{}, orig, fmt.Errorf("expected at least 1 argument after token '<-'")
			}

			if len(rets) == 0 {
				return ir.IrTerm{}, orig, fmt.Errorf("expected at least 1 return value before token '<-'")
			}

			callTerm, _, err := ParseCall(args)
			if err != nil {
				return ir.IrTerm{}, orig, err
			}

			retTokens, err := parser.ParseTokens(rets)
			if err != nil {
				return ir.IrTerm{}, orig, err
			}

			retTerms := make([]ir.IrTerm, len(retTokens))
			for i := range retTokens {
				retTerms[i] = ir.NewTokenTerm(retTokens[i])
			}

			return ir.NewAssignTerm(callTerm, ir.NewTupleTerm(retTerms)), nil, nil
		}

		rets = append(rets, args[0])
	}

	return ParseCall(orig)
}
