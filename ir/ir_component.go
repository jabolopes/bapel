package ir

import (
	"fmt"
)

type IrComponent struct {
	ID     string
	TypeID string
	Length int
}

func (c IrComponent) String() string {
	return fmt.Sprintf("component %s { %s %d } ", c.ID, c.TypeID, c.Length)
}

func NewComponent(id string, typeID string, length int) IrComponent {
	return IrComponent{id, typeID, length}
}
