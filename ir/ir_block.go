package ir

import "fmt"

type blockType int

const (
	moduleBlock = blockType(iota)
	functionBlock
)

func (t blockType) String() string {
	switch t {
	case moduleBlock:
		return "module block"
	case functionBlock:
		return "function block"
	default:
		panic(fmt.Errorf("unhandled block type %d", t))
	}
}

type block struct {
	typ      blockType
	function *struct {
		id     string
		retIDs []string
	}
}

func newFunctionBlock(id string, retIDs []string) block {
	b := block{}
	b.typ = functionBlock
	b.function = &struct {
		id     string
		retIDs []string
	}{id, retIDs}
	return b
}
