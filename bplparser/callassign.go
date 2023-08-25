package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/parser"
)

func ParseCall(args []string) ([]parser.Token, []string, error) {
	orig := args

	tokens, err := parser.ParseTokens(args)
	if err != nil {
		return nil, orig, err
	}

	return tokens, nil, nil
}

// ParseCallAssign parses call and assignment.
//
// Note that a call is an assignment without the '<-' and without any return
// values.
func ParseCallAssign(args []string) ([]parser.Token, []string, error) {
	orig := args

	var rets []string
	for ; len(args) > 0; args = args[1:] {
		if args[0] == "<-" {
			args = args[1:]

			if len(args) == 0 {
				return nil, orig, fmt.Errorf("expected at least 1 argument after token '<-'")
			}

			if len(rets) == 0 {
				return nil, orig, fmt.Errorf("expected at least 1 return value before token '<-'")
			}

			argTokens, _, err := ParseCall(args)
			if err != nil {
				return nil, orig, err
			}

			return argTokens, rets, nil
		}

		rets = append(rets, args[0])
	}

	return ParseCall(orig)
}
