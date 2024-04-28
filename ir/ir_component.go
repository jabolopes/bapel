package ir

import (
	"fmt"
)

type IrComponent struct {
	ElemType IrType
	Length   int
}

func (c IrComponent) String() string {
	return fmt.Sprintf("component [%s, %d]", c.ElemType, c.Length)
}

func NewComponent(elemType IrType, length int) IrComponent {
	return IrComponent{elemType, length}
}
