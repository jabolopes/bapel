package parser

import "fmt"

func Shift[T any](args []T, err error) (T, []T, error) {
	var t T
	if len(args) == 0 {
		return t, args, err
	}
	return args[0], args[1:], nil
}

func ShiftIf[T comparable](args []T, token T, err error) ([]T, error) {
	if len(args) == 0 || args[0] != token {
		return args, err
	}
	return args[1:], nil
}

func ShiftIfEnd[T comparable](args []T, token T, err error) ([]T, error) {
	if len(args) == 0 || args[len(args)-1] != token {
		return args, err
	}
	return args[:len(args)-1], nil
}

func ShiftBalancedParens(args []string) ([]string, []string) {
	count := 0
	for i, arg := range args {
		switch arg {
		case "(":
			count++
		case ")":
			count--

			if count <= 0 {
				return args[0 : i+1], args[i+1:]
			}
		default:
			continue
		}
	}

	return args, nil
}

func EOL[T any](args []T) error {
	if len(args) > 0 {
		return fmt.Errorf("expected end of line; got %v", args)
	}

	return nil
}
