package parser

import (
	"fmt"
	"io"

	"golang.org/x/exp/constraints"
)

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

func ShiftToken[T comparable](args []T, token T) ([]T, error) {
	return ShiftIf(args, token, fmt.Errorf("expected token '%v'; got %v", token, args))
}

func ShiftID[T any](args []T) (T, []T, error) {
	return Shift(args, fmt.Errorf("expected identified; got %v", args))
}

func ShiftIfEnd[T comparable](args []T, token T, err error) ([]T, error) {
	if len(args) == 0 || args[len(args)-1] != token {
		return args, err
	}
	return args[:len(args)-1], nil
}

func ShiftNumber[T constraints.Integer](args []string) (T, []string, error) {
	var t T

	if len(args) == 0 {
		return t, args, io.EOF
	}

	number, err := ParseNumber[T](args[0])
	if err != nil {
		return t, args, err
	}

	return number, args[1:], nil
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
