package parser

import (
	"io"
)

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
