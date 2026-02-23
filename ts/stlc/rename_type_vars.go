package stlc

import (
	"errors"
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/list"
)

type substitution struct {
	source ir.IrType
	target ir.IrType
}

type typeVarRenamer struct {
	context       Context
	err           error
	substitutions list.List[substitution]
}

func (t *typeVarRenamer) lookupSubstitution(sourceTvar string) (ir.IrType, bool) {
	for it := t.substitutions.Iterate(); ; {
		_, substitution, ok := it.Next()
		if !ok {
			break
		}

		if substitution.source.Var == sourceTvar {
			return substitution.target, true
		}
	}

	return ir.IrType{}, false
}

func (t *typeVarRenamer) createSubstitution(sourceTvar string, sourceKind ir.IrKind) (string, error) {
	targetTvar := t.context.GenFreshVarType()

	var err error
	t.context, err = t.context.AddBind(NewTypeVarBind(targetTvar.Var, sourceKind))
	if err != nil {
		return "", err
	}

	t.substitutions = t.substitutions.Add(substitution{ir.NewVarType(sourceTvar), targetTvar})
	return targetTvar.Var, nil
}

func (t *typeVarRenamer) renameImpl(typ ir.IrType) ir.IrType {
	switch typ.Case {
	case ir.AppType:
		c := typ.App
		return ir.NewAppType(t.rename(c.Fun), t.rename(c.Arg))

	case ir.ArrayType:
		c := typ.Array
		return ir.NewArrayType(t.rename(c.ElemType), c.Size)

	case ir.ExistVarType:
		return typ

	case ir.ForallType:
		c := typ.Forall

		origContext := t.context
		origSubstitutions := t.substitutions
		defer func() {
			t.context = origContext
			t.substitutions = origSubstitutions
		}()

		targetTvar, err := t.createSubstitution(c.Var, c.Kind)
		if err != nil {
			t.err = errors.Join(t.err, err)
			return ir.IrType{}
		}

		return ir.NewForallType(targetTvar, c.Kind, t.rename(c.Type))

	case ir.FunType:
		c := typ.Fun
		return ir.NewFunctionType(t.rename(c.Arg), t.rename(c.Ret))

	case ir.LambdaType:
		c := typ.Lambda

		origContext := t.context
		origSubstitutions := t.substitutions
		defer func() {
			t.context = origContext
			t.substitutions = origSubstitutions
		}()

		targetTvar, err := t.createSubstitution(c.Var, c.Kind)
		if err != nil {
			t.err = errors.Join(t.err, err)
			return ir.IrType{}
		}

		return ir.NewLambdaType(targetTvar, c.Kind, t.rename(c.Type))

	case ir.NameType:
		return typ

	case ir.StructType:
		fields := make([]ir.StructField, len(typ.Struct.Fields))
		for i := range typ.Struct.Fields {
			fields[i] = typ.Struct.Fields[i]
			fields[i].Type = t.rename(fields[i].Type)
		}
		return ir.NewStructType(fields)

	case ir.TupleType:
		elems := make([]ir.IrType, len(typ.Tuple.Elems))
		for i := range typ.Tuple.Elems {
			elems[i] = t.rename(typ.Tuple.Elems[i])
		}
		return ir.NewTupleType(elems)

	case ir.VariantType:
		tags := make([]ir.VariantTag, len(typ.Variant.Tags))
		for i := range typ.Variant.Tags {
			tags[i] = typ.Variant.Tags[i]
			tags[i].Type = t.rename(tags[i].Type)
		}
		return ir.NewVariantType(tags)

	case ir.VarType:
		if targetTvar, ok := t.lookupSubstitution(typ.Var); ok {
			return targetTvar
		}
		return typ

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (t *typeVarRenamer) rename(typ ir.IrType) ir.IrType {
	if t.err != nil {
		return ir.IrType{}
	}

	return t.renameImpl(typ)
}

func renameTypeVarsWithSubstitutions(context Context, typ ir.IrType, substitutions []substitution) (Context, ir.IrType, error) {
	origContext := context

	subs := list.New[substitution]()
	for _, substitution := range substitutions {
		subs = subs.Add(substitution)
	}

	renamer := &typeVarRenamer{context, nil /* err */, subs}
	typ = renamer.rename(typ)
	if renamer.err != nil {
		return origContext, ir.IrType{}, renamer.err
	}

	return renamer.context, typ, nil
}

func renameTypeVars(context Context, typ ir.IrType) (Context, ir.IrType, error) {
	return renameTypeVarsWithSubstitutions(context, typ, nil /* substitutions */)
}
