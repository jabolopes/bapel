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

func (t blockType) String() string {
	switch t {
	case moduleBlock:
		return "module block"
	case importsBlock:
		return "imports block"
	case exportsBlock:
		return "exports block"
	case declsBlock:
		return "decls block"
	case functionBlock:
		return "function block"
	case ifThenBlock:
		return "if block"
	case ifElseBlock:
		return "if block"
	case elseBlock:
		return "if block"
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
