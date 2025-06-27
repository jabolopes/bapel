package stlc

import (
	"fmt"
	"log"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ir"
)

type typeReducer struct {
	*log.Logger
}

func (t *typeReducer) reduceImpl(ctx Context, typ ir.IrType) ir.IrType {
	switch {
	case typ.Is(ir.AppType):
		c := typ.App

		fun := t.reduce(ctx, c.Fun)
		arg := t.reduce(ctx, c.Arg)

		if fun.Is(ir.LambdaType) {
			return ir.SubstituteType(fun.Lambda.Type, ir.NewVarType(fun.Lambda.Var), c.Arg)
		}

		return ir.NewAppType(fun, arg)

	case typ.Is(ir.ArrayType):
		c := typ.Array

		elemType := t.reduce(ctx, c.ElemType)
		return ir.NewArrayType(elemType, c.Size)

	case typ.Is(ir.ForallType):
		c := typ.Forall

		var tvar ir.IrType
		var bodyType ir.IrType
		var err error
		ctx, tvar, bodyType, err = ctx.AddFreshType(typ)
		if err != nil {
			panic(err)
		}

		return ir.NewForallType(tvar.Var, c.Kind, t.reduce(ctx, bodyType))

	case typ.Is(ir.FunType):
		c := typ.Fun

		argType := t.reduce(ctx, c.Arg)
		retType := t.reduce(ctx, c.Ret)
		return ir.NewFunctionType(argType, retType)

	case typ.Is(ir.LambdaType):
		c := typ.Lambda

		bodyType := t.reduce(ctx, c.Type)
		return ir.NewLambdaType(c.Var, c.Kind, bodyType)

	case typ.Is(ir.NameType) && ctx.containsAliasBind(typ.Name):
		bind, err := ctx.getAliasBind(typ.Name)
		if err != nil {
			panic(err)
		}
		return t.reduce(ctx, bind.Alias.Type)

	case typ.Is(ir.NameType) && ctx.containsConstBind(typ.Name):
		return typ

	case typ.Is(ir.StructType):
		fields := make([]ir.StructField, 0, len(typ.Fields()))
		for _, field := range typ.Fields() {
			fieldType := t.reduce(ctx, field.Type)
			fields = append(fields, ir.StructField{field.ID, fieldType})
		}
		return ir.NewStructType(fields)

	case typ.Is(ir.TupleType):
		elemTypes := make([]ir.IrType, 0, len(typ.Tuple.Elems))
		for _, elem := range typ.Tuple.Elems {
			elemType := t.reduce(ctx, elem)
			elemTypes = append(elemTypes, elemType)
		}
		return ir.NewTupleType(elemTypes)

	case typ.Is(ir.VariantType):
		tags := make([]ir.VariantTag, 0, len(typ.Tags()))
		for _, tag := range typ.Tags() {
			tagType := t.reduce(ctx, tag.Type)
			tags = append(tags, ir.VariantTag{tag.ID, tagType})
		}
		return ir.NewVariantType(tags)

	case typ.Is(ir.VarType) && ctx.containsTypeVarBind(typ.Var):
		return typ

	default:
		panic(fmt.Errorf("unhandled %T %d: %v", typ.Case, typ.Case, typ))
	}
}

func (t *typeReducer) reduce(ctx Context, typ ir.IrType) ir.IrType {
	if err := isWellformedType(ctx, typ); err != nil {
		panic(fmt.Sprintf("failed to reduce %s because type is not wellformed: %v", typ, err))
	}

	reduced := t.reduceImpl(ctx, typ)
	reduced.Pos = typ.Pos

	glog.V(1).Infof("reduce: %s |- %s => %s", ctx.String(), typ, reduced)
	return reduced
}
