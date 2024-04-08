package parser

import (
	"fmt"
	"io"

	"golang.org/x/exp/constraints"
)

func shiftIf[T comparable](args []T, token T, err error) ([]T, error) {
	if len(args) == 0 || args[0] != token {
		return args, err
	}
	return args[1:], nil
}

func shiftIfEnd[T comparable](args []T, token T, err error) ([]T, error) {
	if len(args) == 0 || args[len(args)-1] != token {
		return args, err
	}
	return args[:len(args)-1], nil
}

func shift[T any](args []T, err error) (T, []T, error) {
	var t T
	if len(args) == 0 {
		return t, args, err
	}
	return args[0], args[1:], nil
}

func ShiftLiteral[T comparable](args []T, token T) ([]T, error) {
	return shiftIf(args, token, fmt.Errorf("expected token '%v'; got %v", token, args))
}

func ShiftLiteralEnd[T comparable](args []T, token T) ([]T, error) {
	return shiftIfEnd(args, token, fmt.Errorf("expected token '%v' at end of line; got %v", token, args))
}

func ShiftID[T any](args []T) (T, []T, error) {
	return shift(args, fmt.Errorf("expected identifier; got %v", args))
}

func ShiftNumber[T constraints.Integer](args []string) (T, []string, error) {
	var t T

	if len(args) == 0 {
		return t, args, io.EOF
	}

	number, err := parseNumber[T](args[0])
	if err != nil {
		return t, args, err
	}

	return number, args[1:], nil
}

func ShiftToken(args []string) (Token, []string, error) {
	if len(args) == 0 {
		return Token{}, args, io.EOF
	}

	token, err := parseToken(args[0])
	if err != nil {
		return Token{}, args, err
	}

	return token, args[1:], nil
}

func EOL[T any](args []T) error {
	if len(args) > 0 {
		return fmt.Errorf("expected end of line; got %v", args)
	}

	return nil
}
