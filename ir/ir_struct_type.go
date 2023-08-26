package ir

import (
	"fmt"
	"strings"
)

type StructField struct {
	// TODO: Rename to ID.
	Name string
	Type IrType
}

func (f StructField) String() string {
	return fmt.Sprintf("%s %s", f.Name, f.Type)
}

type IrStructType struct {
	Fields []StructField
}

func (t IrStructType) Names() []string {
	names := make([]string, len(t.Fields))
	for i, field := range t.Fields {
		names[i] = field.Name
	}
	return names
}

func (t IrStructType) String() string {
	var b strings.Builder
	b.WriteString("{")
	if len(t.Fields) > 0 {
		b.WriteString(fmt.Sprintf("%s", t.Fields[0]))
		for _, field := range t.Fields[1:] {
			b.WriteString(fmt.Sprintf(", %s", field))
		}
	}
	b.WriteString("}")
	return b.String()
}
