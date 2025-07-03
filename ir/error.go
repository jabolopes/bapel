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

type Validation struct {
	Errors []Error
}

func (v *Validation) OK() bool {
	return len(v.Errors) == 0
}

func (v *Validation) AddError(err Error) *Validation {
	v.Errors = append(v.Errors, err)
	return v
}

func (v *Validation) AddErr(pos Pos, err error) *Validation {
	return v.AddError(NewError(pos, err.Error()))
}

func (v *Validation) Join(other Validation) *Validation {
	v.Errors = append(v.Errors, other.Errors...)
	return v
}

func (v *Validation) Err() error {
	if v.OK() {
		return nil
	}
	return TopErrors(v.Errors)
}
