package ir

import (
	"fmt"
	"strings"
)

type StructField struct {
	ID   string
	Type IrType
}

func (f StructField) String() string {
	return fmt.Sprintf("%s %s", f.ID, f.Type)
}

type IrStructType struct {
	Fields []StructField
}

func (t IrStructType) FieldByIndex(index int) (StructField, bool) {
	if index >= 0 && index < len(t.Fields) {
		return t.Fields[index], true
	}
	return StructField{}, false
}

func (t IrStructType) FieldByID(id string) (StructField, bool) {
	for _, field := range t.Fields {
		if field.ID == id {
			return field, true
		}
	}
	return StructField{}, false
}

func (t IrStructType) FieldIDs() []string {
	ids := make([]string, len(t.Fields))
	for i, field := range t.Fields {
		ids[i] = field.ID
	}
	return ids
}

func (t IrStructType) String() string {
	var b strings.Builder
	b.WriteString("{")
	if len(t.Fields) > 0 {
		b.WriteString(t.Fields[0].String())
		for _, field := range t.Fields[1:] {
			b.WriteString(fmt.Sprintf(", %s", field))
		}
	}
	b.WriteString("}")
	return b.String()
}
