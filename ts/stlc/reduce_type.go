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

func (t *typeReducer) reduceImpl(typ ir.IrType) ir.IrType {
	switch {
	case typ.Is(ir.AppType):
		c := typ.App

		fun := t.reduce(c.Fun)
		arg := t.reduce(c.Arg)

		if fun.Is(ir.LambdaType) {
			return ir.SubstituteType(fun.Lambda.Type, ir.NewVarType(fun.Lambda.Var), c.Arg)
		}

		return ir.NewAppType(fun, arg)

	case typ.Is(ir.ArrayType):
		c := typ.Array

		elemType := t.reduce(c.ElemType)
		return ir.NewArrayType(elemType, c.Size)

	case typ.Is(ir.ForallType):
		c := typ.Forall

		innerType := t.reduce(c.Type)
		return ir.NewForallType(c.Var, c.Kind, innerType)

	case typ.Is(ir.FunType):
		c := typ.Fun

		argType := t.reduce(c.Arg)
		retType := t.reduce(c.Ret)
		return ir.NewFunctionType(argType, retType)

	case typ.Is(ir.LambdaType):
		c := typ.Lambda

		bodyType := t.reduce(c.Type)
		return ir.NewLambdaType(c.Var, c.Kind, bodyType)

	case typ.Is(ir.NameType) && t.context.containsAliasBind(typ.Name):
		bind, err := t.context.getAliasBind(typ.Name)
		if err != nil {
			panic(err)
		}
		return t.reduce(bind.Alias.Type)

	case typ.Is(ir.NameType) && t.context.containsConstBind(typ.Name):
		return typ

	case typ.Is(ir.StructType):
		fields := make([]ir.StructField, 0, len(typ.Fields()))
		for _, field := range typ.Fields() {
			fieldType := t.reduce(field.Type)
			fields = append(fields, ir.StructField{field.ID, fieldType})
		}
		return ir.NewStructType(fields)

	case typ.Is(ir.TupleType):
		elemTypes := make([]ir.IrType, 0, len(typ.Tuple.Elems))
		for _, elem := range typ.Tuple.Elems {
			elemType := t.reduce(elem)
			elemTypes = append(elemTypes, elemType)
		}
		return ir.NewTupleType(elemTypes)

	case typ.Is(ir.VariantType):
		tags := make([]ir.VariantTag, 0, len(typ.Tags()))
		for _, tag := range typ.Tags() {
			tagType := t.reduce(tag.Type)
			tags = append(tags, ir.VariantTag{tag.ID, tagType})
		}
		return ir.NewVariantType(tags)

	case typ.Is(ir.VarType):
		return typ

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (t *typeReducer) reduce(typ ir.IrType) ir.IrType {
	reduced := t.reduceImpl(typ)
	reduced.Pos = typ.Pos

	t.Printf("reduce: %s |- %s => %s", t.context.StringNoImports(), typ, reduced)
	return reduced
}
