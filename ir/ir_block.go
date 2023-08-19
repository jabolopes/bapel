package ir

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
