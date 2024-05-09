package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func inferKindApply(context Context, fun ir.IrType, arg ir.IrType) (ir.IrKind, error) {
	funKind, err := InferKind(context, fun)
	if err != nil {
		return ir.IrKind{}, err
	}

	argKind, err := InferKind(context, arg)
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
	switch typ.Case {
	case ir.AppType:
		c := typ.App
		return inferKindApply(context, c.Fun, c.Arg)

	case ir.ArrayType:
		return ir.NewTypeKind(), nil

	case ir.ForallType:
		newContext, bodyType, err := context.AddFreshType(typ)
		if err != nil {
			return ir.IrKind{}, err
		}

		kind, err := InferKind(newContext, bodyType)
		if err != nil {
			return ir.IrKind{}, err
		}
		if !ir.EqualsKind(kind, ir.NewTypeKind()) {
			return ir.IrKind{}, fmt.Errorf("expected type %s to have kind %s, but got kind %s", typ, ir.NewTypeKind(), kind)
		}

		return kind, nil

	case ir.FunType:
		return ir.NewTypeKind(), nil

	case ir.LambdaType:
		newContext, bodyType, err := context.AddFreshType(typ)
		if err != nil {
			return ir.IrKind{}, err
		}

		retKind, err := InferKind(newContext, bodyType)
		if err != nil {
			return ir.IrKind{}, err
		}

		return ir.NewArrowKind(typ.Lambda.Kind, retKind), nil

	case ir.NameType:
		if bind, err := context.getAliasBind(typ.Name); err == nil {
			return InferKind(context, bind.Alias.Type)
		}

		if !context.containsNameBind(typ.Name) {
			return ir.IrKind{}, fmt.Errorf("type %q is undefined", typ.Name)
		}

		// TODO: I suspect the kind should come from the name bind, otherwise we're
		// assuming all name binds have type kind, and that can't be true.
		return ir.NewTypeKind(), nil

	case ir.StructType:
		return ir.NewTypeKind(), nil

	case ir.TupleType:
		return ir.NewTypeKind(), nil

	case ir.VarType:
		bind, err := context.getTypeVarBind(typ.Var)
		if err != nil {
			return ir.IrKind{}, err
		}
		return bind.TypeVar.Kind, nil

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func InferKind(context Context, typ ir.IrType) (ir.IrKind, error) {
	kind, err := inferKindImpl(context, typ)
	if err != nil {
		return ir.IrKind{}, fmt.Errorf("%v\n  inferring kind for type %s", err, typ)
	}

	return kind, nil
}
