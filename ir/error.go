package ir

import (
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	Pos     Pos
	Message string
}

func (e Error) String() string {
	return fmt.Sprintf("%v:\n  %s", e.Pos, e.Message)
}

func NewError(pos Pos, message string) Error {
	return Error{pos, message}
}

func TopErrors(errs []Error) error {
	var b strings.Builder

	firstErrors := errs[:min(10, len(errs))]

	interleave(firstErrors, func() { b.WriteString("\n\n") }, func(_ int, err Error) {
		b.WriteString(err.String())
	})

	if len(errs) > len(firstErrors) {
		b.WriteString("\n\nToo many errors.")
	}

	return errors.New(b.String())
}
