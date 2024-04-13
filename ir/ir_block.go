package ir

import "fmt"

type blockType int

const (
	moduleBlock = blockType(iota)
)

func (t blockType) String() string {
	switch t {
	case moduleBlock:
		return "module block"
	default:
		panic(fmt.Errorf("unhandled block type %d", t))
	}
}

type block struct {
	typ blockType
}
