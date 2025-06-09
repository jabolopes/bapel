package stlc

import (
	"fmt"
	"log"

	"github.com/jabolopes/bapel/ir"
)

type typeReducer struct {
	*log.Logger
	context Context
}

func (t *typeReducer) reduceImpl(typ ir.IrType) (ir.IrType, error) {
	switch {
	case typ.Is(ir.AppType):
		c := typ.App

		fun, err := t.reduce(c.Fun)
		if err != nil {
			return ir.IrType{}, err
		}

		arg, err := t.reduce(c.Arg)
		if err != nil {
			return ir.IrType{}, err
		}

		if fun.Is(ir.LambdaType) {
			return ir.SubstituteType(fun.Lambda.Type, ir.NewVarType(fun.Lambda.Var), c.Arg), nil
		}

		return ir.NewAppType(fun, arg), nil

	case typ.Is(ir.ArrayType):
		c := typ.Array

		elemType, err := t.reduce(c.ElemType)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewArrayType(elemType, c.Size), nil

	case typ.Is(ir.ForallType):
		c := typ.Forall

		innerType, err := t.reduce(c.Type)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewForallType(c.Var, c.Kind, innerType), nil

	case typ.Is(ir.FunType):
		c := typ.Fun

		argType, err := t.reduce(c.Arg)
		if err != nil {
			return ir.IrType{}, err
		}

		retType, err := t.reduce(c.Ret)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewFunctionType(argType, retType), nil

	case typ.Is(ir.LambdaType):
		c := typ.Lambda

		bodyType, err := t.reduce(c.Type)
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewLambdaType(c.Var, c.Kind, bodyType), nil

	case typ.Is(ir.NameType) && t.context.containsAliasBind(typ.Name):
		bind, err := t.context.getAliasBind(typ.Name)
		if err != nil {
			panic(err)
		}
		return t.reduce(bind.Alias.Type)

	case typ.Is(ir.NameType) && t.context.containsConstBind(typ.Name):
		return typ, nil

	case typ.Is(ir.StructType):
		fields := make([]ir.StructField, 0, len(typ.Fields()))
		for _, field := range typ.Fields() {
			fieldType, err := t.reduce(field.Type)
			if err != nil {
				return ir.IrType{}, err
			}

			fields = append(fields, ir.StructField{field.ID, fieldType})
		}
		return ir.NewStructType(fields), nil

	case typ.Is(ir.TupleType):
		elemTypes := make([]ir.IrType, 0, len(typ.Tuple.Elems))
		for _, elem := range typ.Tuple.Elems {
			elemType, err := t.reduce(elem)
			if err != nil {
				return ir.IrType{}, err
			}

			elemTypes = append(elemTypes, elemType)
		}
		return ir.NewTupleType(elemTypes), nil

	case typ.Is(ir.VariantType):
		tags := make([]ir.VariantTag, 0, len(typ.Tags()))
		for _, tag := range typ.Tags() {
			tagType, err := t.reduce(tag.Type)
			if err != nil {
				return ir.IrType{}, err
			}

			tags = append(tags, ir.VariantTag{tag.ID, tagType})
		}
		return ir.NewVariantType(tags), nil

	case typ.Is(ir.VarType):
		return typ, nil

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (t *typeReducer) reduce(typ ir.IrType) (ir.IrType, error) {
	reduced, err := t.reduceImpl(typ)
	if err != nil {
		return ir.IrType{}, fmt.Errorf("%v\n  reducing %s", err, typ)
	}

	reduced.Pos = typ.Pos

	t.Printf("reduce: %s |- %s => %s", t.context.StringNoImports(), typ, reduced)
	return reduced, nil
}
