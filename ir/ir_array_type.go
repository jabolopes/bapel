package ir

import "fmt"

type IrArrayType struct {
	ElemType IrType
	Size     int
}

func (t IrArrayType) String() string {
	return fmt.Sprintf("[%v]", t.ElemType)
}

func MatchesArrayType(formal, actual IrArrayType, widen bool) error {
	if err := MatchesType(formal.ElemType, actual.ElemType, widen); err != nil {
		return fmt.Errorf("mismatch in array element types: %v", err)
	}

	if formal.Size != actual.Size {
		return fmt.Errorf("expected array with %d elements; got %d elements", formal.Size, actual.Size)
	}

	return nil
}

func SizeOfArrayType(typ IrArrayType) int {
	return SizeOfType(typ.ElemType) * typ.Size
}
