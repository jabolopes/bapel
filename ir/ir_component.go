package ir

import (
	"fmt"
)

type IrComponent struct {
	ID       string
	ElemType IrType
	Length   int
}

func (c IrComponent) String() string {
	return fmt.Sprintf("component %s { %s %d } ", c.ID, c.ElemType, c.Length)
}

func NewComponent(id string, elemType IrType, length int) IrComponent {
	return IrComponent{id, elemType, length}
}
