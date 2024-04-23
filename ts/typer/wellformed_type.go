package typer

import "fmt"

func WellformedType(context Context, typ Type) error {
	switch typ.Case {
	case ExistVarType:
		if !context.ContainsExistVarType(*typ.ExistVar) {
			return fmt.Errorf("type %s is not wellformed", typ)
		}
		return nil

	case ForallType:
		context = context.AddType(NewVarType(typ.Forall.Var))
		return WellformedType(context, typ.Forall.Type)

	case FunType:
		if err := WellformedType(context, typ.Fun.Arg); err != nil {
			return err
		}
		return WellformedType(context, typ.Fun.Ret)

	case NameType:
		if !context.ContainsNameType(*typ.Name) {
			return fmt.Errorf("type %s is not wellformed", typ)
		}
		return nil

	case VarType:
		if !context.ContainsVarType(*typ.Var) {
			return fmt.Errorf("type %s is not wellformed", typ)
		}
		return nil

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}
