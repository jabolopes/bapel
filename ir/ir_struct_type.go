package ir

import (
	"fmt"
	"sort"
	"strings"
)

type StructField struct {
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

func MatchesStructType(formal, actual IrStructType, widen bool) error {
	if len(formal.Fields) != len(actual.Fields) {
		return fmt.Errorf("expected %d fields; got %d", len(formal.Fields), len(actual.Fields))
	}

	formalFields := append([]StructField{}, formal.Fields...)
	actualFields := append([]StructField{}, actual.Fields...)

	sort.Slice(formalFields, func(i, j int) bool {
		return formalFields[i].Name < formalFields[j].Name
	})
	sort.Slice(actualFields, func(i, j int) bool {
		return actualFields[i].Name < actualFields[j].Name
	})

	for i := range formalFields {
		if formalFields[i].Name != actualFields[i].Name {
			return fmt.Errorf("expected field names %v; got %v", formal.Names(), actual.Names())
		}

		if err := MatchesType(formalFields[i].Type, actualFields[i].Type, widen); err != nil {
			return err
		}
	}

	return nil
}
