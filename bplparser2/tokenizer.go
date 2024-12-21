package bplparser2

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/constraints"
)

func parseNumber[T constraints.Integer](arg string) (T, error) {
	var value T

	if strings.HasPrefix(arg, "0x") {
		// Hexadecimal
		_, err := fmt.Sscanf(arg, "0x%x", &value)

		return value, err
	}

	// Decimal.
	_, err := fmt.Sscanf(arg, "%d", &value)
	return value, err
}

func parseToken(text string) (parser.Token, error) {
	if value, err := parseNumber[int64](text); err == nil {
		return parser.Token{parser.NumberToken, text, value}, nil
	}
	return parser.NewIDToken(text), nil
}
