package stlc

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ir"
)

func (t *Inferencer) solveTypeImpl(typ ir.IrType) ir.IrType {
	switch {
	case t.isExistVarAssigned(typ):
		solution := t.existVarSolution(typ)
		return t.solveType(solution)

	case t.isExistVarUnassigned(typ):
		return typ

	case typ.Is(ir.AppType):
		c := typ.App

		fun := t.solveType(c.Fun)
		arg := t.solveType(c.Arg)
		return ir.NewAppType(fun, arg)

	case typ.Is(ir.ArrayType):
		c := typ.Array

		elemType := t.solveType(c.ElemType)
		return ir.NewArrayType(elemType, c.Size)

	case typ.Is(ir.ForallType):
		c := typ.Forall

		return ir.NewForallType(c.Var, c.Kind, t.solveType(c.Type))

	case typ.Is(ir.FunType):
		c := typ.Fun

		argType := t.solveType(c.Arg)
		retType := t.solveType(c.Ret)
		return ir.NewFunctionType(argType, retType)

	case typ.Is(ir.LambdaType):
		c := typ.Lambda

		return ir.NewLambdaType(c.Var, c.Kind, t.solveType(c.Type))

	case typ.Is(ir.NameType):
		return typ

	case typ.Is(ir.StructType):
		fields := make([]ir.StructField, 0, len(typ.Fields()))
		for _, field := range typ.Fields() {
			fieldType := t.solveType(field.Type)
			fields = append(fields, ir.StructField{field.ID, fieldType})
		}
		return ir.NewStructType(fields)

	case typ.Is(ir.TupleType):
		elemTypes := make([]ir.IrType, 0, len(typ.Tuple.Elems))
		for _, elem := range typ.Tuple.Elems {
			elemType := t.solveType(elem)
			elemTypes = append(elemTypes, elemType)
		}
		return ir.NewTupleType(elemTypes)

	case typ.Is(ir.VariantType):
		tags := make([]ir.VariantTag, 0, len(typ.Tags()))
		for _, tag := range typ.Tags() {
			tagType := t.solveType(tag.Type)
			tags = append(tags, ir.VariantTag{tag.ID, tagType})
		}
		return ir.NewVariantType(tags)

	case typ.Is(ir.VarType):
		return typ

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (t *Inferencer) solveType(typ ir.IrType) ir.IrType {
	solved := t.solveTypeImpl(typ)
	solved.Pos = typ.Pos

	glog.V(1).Infof("solveType: %s |- %s => %s", t.context.String(), typ, solved)
	return solved
}
