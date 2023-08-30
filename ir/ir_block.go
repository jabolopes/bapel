package ir

import "fmt"

type blockType int

const (
	moduleBlock = blockType(iota)
	importsBlock
	exportsBlock
	declsBlock
	functionBlock
	ifThenBlock
	ifElseBlock
	elseBlock
)

type block struct {
	typ      blockType
	function *struct {
		id     string
		retIDs []string
	}
}

func newBlock(typ blockType) block {
	if typ == functionBlock {
		panic(fmt.Errorf("use newFunctionBlock() for function blocks instead"))
	}

	b := block{}
	b.typ = typ
	return b
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
