package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type typePredicator struct {
	context Context
	tvars   []ir.VarKind
}

func (t *typePredicator) predicateImpl(typ ir.IrType) (ir.IrType, error) {
	switch typ.Case {
	case ir.AppType:
		c := typ.App

		fun, err := t.predicate(c.Fun)
		if err != nil {
			return ir.IrType{}, err
		}

		arg, err := t.predicate(c.Arg)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewAppType(fun, arg), nil

	case ir.ArrayType:
		c := typ.Array

		elem, err := t.predicate(c.ElemType)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewArrayType(elem, c.Size), nil

	case ir.ExistVarType:
		return typ, nil

	case ir.ForallType:
		c := typ.Forall

		var tvar ir.IrType
		var bodyType ir.IrType
		var err error
		t.context, tvar, bodyType, err = t.context.AddFreshType(typ)
		if err != nil {
			return ir.IrType{}, err
		}

		// TODO: AddFreshType should return the ir.VarKind.
		t.tvars = append(t.tvars, ir.VarKind{tvar.Var, c.Kind})

		return t.predicate(bodyType)

	case ir.FunType:
		c := typ.Fun

		arg, err := t.predicate(c.Arg)
		if err != nil {
			return ir.IrType{}, err
		}

		ret, err := t.predicate(c.Ret)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewFunctionType(arg, ret), nil

	case ir.LambdaType:
		c := typ.Lambda

		var tvar ir.IrType
		var bodyType ir.IrType
		var err error
		t.context, tvar, bodyType, err = t.context.AddFreshType(typ)
		if err != nil {
			return ir.IrType{}, err
		}

		// TODO: AddFreshType should return the ir.VarKind.
		t.tvars = append(t.tvars, ir.VarKind{tvar.Var, c.Kind})

		return t.predicate(bodyType)

	case ir.NameType:
		return typ, nil

	case ir.StructType:
		fields := make([]ir.StructField, 0, len(typ.Fields()))

		for _, field := range typ.Fields() {
			var err error
			field.Type, err = t.predicate(field.Type)
			if err != nil {
				return ir.IrType{}, err
			}

			fields = append(fields, field)
		}

		return ir.NewStructType(fields), nil

	case ir.TupleType:
		elems := make([]ir.IrType, 0, len(typ.Elems()))

		for _, elem := range typ.Elems() {
			var err error
			elem, err = t.predicate(elem)
			if err != nil {
				return ir.IrType{}, err
			}

			elems = append(elems, elem)
		}

		return ir.NewTupleType(elems), nil

	case ir.VariantType:
		tags := make([]ir.VariantTag, 0, len(typ.Fields()))

		for _, tag := range typ.Tags() {
			var err error
			tag.Type, err = t.predicate(tag.Type)
			if err != nil {
				return ir.IrType{}, err
			}

			tags = append(tags, tag)
		}

		return ir.NewVariantType(tags), nil

	case ir.VarType:
		return typ, nil

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (t *typePredicator) predicate(typ ir.IrType) (ir.IrType, error) {
	newType, err := t.predicateImpl(typ)
	if err != nil {
		return ir.IrType{}, err
	}

	newType.Pos = typ.Pos
	return newType, nil
}
