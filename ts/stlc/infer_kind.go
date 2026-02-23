package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func inferKindApply(context Context, fun ir.IrType, arg ir.IrType) (ir.IrKind, error) {
	funKind, err := inferKind(context, fun)
	if err != nil {
		return ir.IrKind{}, err
	}

	argKind, err := inferKind(context, arg)
	if err != nil {
		return ir.IrKind{}, err
	}

	if !funKind.Is(ir.ArrowKind) {
		return ir.IrKind{}, fmt.Errorf("expected arrow kind (%s) in type application", ir.NewArrowKind(ir.NewTypeKind(), ir.NewTypeKind()))
	}

	if !ir.EqualsKind(funKind.Arrow.Arg, argKind) {
		return ir.IrKind{}, fmt.Errorf("expected argument in type application (%s) to match function argument (%s)", argKind, funKind.Arrow.Arg)
	}

	return funKind.Arrow.Ret, nil
}

func inferKindImpl(context Context, typ ir.IrType) (ir.IrKind, error) {
	switch {
	case typ.Is(ir.AppType):
		c := typ.App
		return inferKindApply(context, c.Fun, c.Arg)

	case typ.Is(ir.ArrayType):
		return ir.NewTypeKind(), nil

	case typ.Is(ir.ExistVarType):
		return ir.NewTypeKind(), nil

	case typ.Is(ir.ForallType):
		newContext, _, bodyType, err := context.AddFreshType(typ)
		if err != nil {
			return ir.IrKind{}, err
		}

		kind, err := inferKind(newContext, bodyType)
		if err != nil {
			return ir.IrKind{}, err
		}
		if !ir.EqualsKind(kind, ir.NewTypeKind()) {
			return ir.IrKind{}, fmt.Errorf("expected type %s to have kind %s, but got kind %s", typ, ir.NewTypeKind(), kind)
		}

		return kind, nil

	case typ.Is(ir.FunType):
		return ir.NewTypeKind(), nil

	case typ.Is(ir.LambdaType):
		newContext, _, bodyType, err := context.AddFreshType(typ)
		if err != nil {
			return ir.IrKind{}, err
		}

		retKind, err := inferKind(newContext, bodyType)
		if err != nil {
			return ir.IrKind{}, err
		}

		return ir.NewArrowKind(typ.Lambda.Kind, retKind), nil

	case typ.Is(ir.NameType) && context.containsAliasBind(typ.Name):
		bind, err := context.getAliasBind(typ.Name)
		if err != nil {
			panic(err)
		}
		return inferKind(context, bind.Alias.Type)

	case typ.Is(ir.NameType) && context.containsConstBind(typ.Name):
		bind, err := context.getConstBind(typ.Name)
		if err != nil {
			panic(err)
		}
		return bind.Const.Kind, nil

	case typ.Is(ir.NameType):
		return ir.IrKind{}, fmt.Errorf("type %q is undefined", typ.Name)

	case typ.Is(ir.StructType):
		return ir.NewTypeKind(), nil

	case typ.Is(ir.TupleType):
		return ir.NewTypeKind(), nil

	case typ.Is(ir.VariantType):
		return ir.NewTypeKind(), nil

	case typ.Is(ir.VarType):
		bind, err := context.getTypeVarBind(typ.Var)
		if err != nil {
			return ir.IrKind{}, err
		}
		return bind.TypeVar.Kind, nil

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func inferKind(context Context, typ ir.IrType) (ir.IrKind, error) {
	kind, err := inferKindImpl(context, typ)
	if err != nil {
		return ir.IrKind{}, fmt.Errorf("%v\n  inferring kind for type %s", err, typ)
	}

	return kind, nil
}
