package ir

import (
	"fmt"
)

type IrEntity struct {
	ID     string
	Length int
}

func (e IrEntity) String() string {
	return fmt.Sprintf("entity %s { %d } ", e.ID, e.Length)
}

func NewEntity(id string, length int) IrEntity {
	return IrEntity{id, length}
}
