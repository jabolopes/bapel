package typer

import (
	"fmt"
)

// TypeCase is a type case. Denoted by A, B, C, ...
type TypeCase int

const (
	// Existential variable.
	//
	// â
	ExistVarType TypeCase = iota
	// Forall type.
	//
	// forall a. A.
	ForallType
	// Function type.
	//
	// A -> B
	FunType
	// Name type (or unit type).
	//
	// 1
	NameType
	// Type variable.
	//
	// a
	VarType
)

type existVar struct {
	Var string
}

type forall struct {
	Var  string
	Type Type
}

type function struct {
	Arg Type
	Ret Type
}

type name struct {
	ID string
}

type typeVar struct {
	Var string
}

type Type struct {
	Case     TypeCase
	ExistVar *existVar
	Forall   *forall
	Fun      *function
	Name     *name
	Var      *typeVar
}

func (t Type) String() string {
	{
		var d Type
		if t == d {
			return ""
		}
	}

	switch t.Case {
	case ExistVarType:
		return fmt.Sprintf("^%s", t.ExistVar.Var)
	case ForallType:
		return fmt.Sprintf("forall [%s] %s", t.Forall.Var, t.Forall.Type)
	case FunType:
		return fmt.Sprintf("%s -> %s", t.Fun.Arg, t.Fun.Ret)
	case NameType:
		return t.Name.ID
	case VarType:
		return fmt.Sprintf("'%s", t.Var)
	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}

func (t Type) Is(c TypeCase) bool {
	return t.Case == c
}

func NewExistVarType(evar string) Type {
	return Type{
		Case:     ExistVarType,
		ExistVar: &existVar{evar},
	}
}

func NewForallType(tvar string, typ Type) Type {
	return Type{
		Case:   ForallType,
		Forall: &forall{tvar, typ},
	}
}

func NewFunType(arg, ret Type) Type {
	return Type{
		Case: FunType,
		Fun:  &function{arg, ret},
	}
}

func NewVarType(tvar string) Type {
	return Type{
		Case: VarType,
		Var:  &typeVar{tvar},
	}
}
