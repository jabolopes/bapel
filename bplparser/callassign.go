package bplparser

import "fmt"

// ParseCallAssign parses call and assignment.
//
// Note that a call is an assignment without the '<-' and without any return
// values.
func ParseCallAssign(args []string) ([]string, []string, error) {
	orig := args

	var rets []string
	for ; len(args) > 0; args = args[1:] {
		if args[0] == "<-" {
			args = args[1:]

			if len(args) == 0 {
				return nil, nil, fmt.Errorf("expected at least 1 argument after token '<-'")
			}

			if len(rets) == 0 {
				return nil, nil, fmt.Errorf("expected at least 1 return value before token '<-'")
			}

			return args, rets, nil
		}

		rets = append(rets, args[0])
	}

	return orig, nil, nil
}
