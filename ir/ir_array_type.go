package ir

import "fmt"

type IrArrayType struct {
	ElemType IrType
	Size     int
}

func (t IrArrayType) String() string {
	return fmt.Sprintf("[%v]", t.ElemType)
}
