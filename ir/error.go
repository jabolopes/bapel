package ir

import (
	"fmt"
)

type Error struct {
	Pos     Pos
	Message string
}

func (e Error) String() string {
	return fmt.Sprintf("%v:\n  %s", e.Pos, e.Message)
}

func NewError(Pos Pos, Message string) Error {
	return Error{Pos, Message}
}
