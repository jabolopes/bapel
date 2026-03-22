package ir

import (
	"cmp"
	"fmt"
)

type IrKindCase int

const (
	// Kind of types.
	TypeKind IrKindCase = iota
	// Kind of type constructors.
	ArrowKind
)

func (c IrKindCase) String() string {
	switch c {
	case TypeKind:
		return "Type"
	case ArrowKind:
		return "Arrow"
	default:
		panic(fmt.Sprintf("unhandled %T %d", c, c))
	}
}

type typeKind struct{}

func (k *typeKind) String() string {
	return "∗"
}

type arrowKind struct {
	Arg IrKind
	Ret IrKind
}

func (k *arrowKind) String() string {
	if k.Arg.Is(ArrowKind) {
		return fmt.Sprintf("(%s) -> %s", k.Arg, k.Ret)
	}
	return fmt.Sprintf("%s -> %s", k.Arg, k.Ret)
}

type IrKind struct {
	Case  IrKindCase
	Type  *typeKind
	Arrow *arrowKind
}

func (k IrKind) String() string {
	if k.Case == 0 && k.Type == nil {
		return ""
	}

	switch k.Case {
	case TypeKind:
		return k.Type.String()
	case ArrowKind:
		return k.Arrow.String()
	default:
		panic(fmt.Sprintf("unhandled %T %d", k.Case, k.Case))
	}
}

func (k IrKind) Is(c IrKindCase) bool {
	return k.Case == c
}

func NewTypeKind() IrKind {
	return IrKind{
		Case: TypeKind,
		Type: &typeKind{},
	}
}

func NewArrowKind(fun, arg IrKind) IrKind {
	return IrKind{
		Case:  ArrowKind,
		Arrow: &arrowKind{fun, arg},
	}
}

func CompareKind(k1, k2 IrKind) int {
	if c := cmp.Compare(k1.Case, k2.Case); c != 0 {
		return c
	}

	switch k1.Case {
	case TypeKind:
		return 0
	case ArrowKind:
		if c1 := CompareKind(k1.Arrow.Arg, k2.Arrow.Arg); c1 != 0 {
			return c1
		}
		return CompareKind(k1.Arrow.Ret, k2.Arrow.Ret)
	default:
		panic(fmt.Errorf("unhandled %T %d", k1.Case, k1.Case))
	}
}

func EqualsKind(k1, k2 IrKind) bool {
	switch {
	case k1.Is(TypeKind) && k2.Is(TypeKind):
		return true
	case k1.Is(ArrowKind) && k2.Is(ArrowKind):
		return EqualsKind(k1.Arrow.Ret, k2.Arrow.Ret) &&
			EqualsKind(k1.Arrow.Arg, k2.Arrow.Arg)
	default:
		return false
	}
}

type VarKind struct {
	// Type variable.
	Var string
	// Kind of type variable.
	Kind IrKind
}

func (t VarKind) String() string {
	return fmt.Sprintf("%s :: %s", t.Var, t.Kind)
}

func CompareVarKind(vk1, vk2 VarKind) int {
	if c := cmp.Compare(vk1.Var, vk2.Var); c != 0 {
		return c
	}
	return CompareKind(vk1.Kind, vk2.Kind)
}
