package ir

import (
	"fmt"
)

type IrEntity struct {
	ID     string
	Length int
}

func (d IrEntity) String() string {
	return fmt.Sprintf("entity %s : ", d.ID)
}

func NewEntity(id string, length int) IrEntity {
	return IrEntity{id, length}
}
