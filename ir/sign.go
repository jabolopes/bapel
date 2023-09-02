package ir

import "fmt"

type Sign byte

const (
	Unsigned = Sign(iota)
	Signed
)

func (s Sign) String() string {
	switch s {
	case Unsigned:
		return "unsigned"
	case Signed:
		return "signed"
	default:
		panic(fmt.Errorf("unhandled Sign %d", s))
	}
}
