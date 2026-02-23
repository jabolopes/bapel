package stlc

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ir"
)

func sameTags(left, right ir.IrType) bool {
	if len(left.Tags()) != len(right.Tags()) {
		return false
	}

	for i := range left.Tags() {
		if left.Tags()[i].ID != right.Tags()[i].ID {
			return false
		}
	}

	return true
}

func sameFields(left, right ir.IrType) bool {
	if len(left.Fields()) != len(right.Fields()) {
		return false
	}

	for i := range left.Fields() {
		if left.Fields()[i].ID != right.Fields()[i].ID {
			return false
		}
	}

	return true
}

type existVar struct {
	solution *ir.IrType
}

func (t *Inferencer) isExistVarUnassigned(tvar ir.IrType) bool {
	if !tvar.Is(ir.ExistVarType) {
		return false
	}
	existVar, ok := t.existVars[tvar.ExistVar]
	return ok && existVar.solution == nil
}

func (t *Inferencer) isExistVarAssigned(tvar ir.IrType) bool {
	if !tvar.Is(ir.ExistVarType) {
		return false
	}
	existVar, ok := t.existVars[tvar.ExistVar]
	return ok && existVar.solution != nil
}

func (t *Inferencer) existVarSolution(evar ir.IrType) ir.IrType {
	if !evar.Is(ir.ExistVarType) {
		panic(fmt.Errorf("expected existential variable; got %v", evar))
	}

	existVar := t.existVars[evar.ExistVar]
	solution := existVar.solution
	if t.isExistVarAssigned(*solution) {
		typ := t.existVarSolution(*solution)
		existVar.solution = &typ
		t.existVars[evar.ExistVar] = existVar
		return typ
	}
	return *solution
}

func (t *Inferencer) canAssign(evar, typ ir.IrType) bool {
	if typ.Is(ir.ForallType) {
		return false
	}

	if typ.Is(ir.ExistVarType) {
		// ExistVar generation is monotonically increasing, as such,
		// variables can only be assigned to earlier variables (to avoid
		// cycles).
		return typ.ExistVar < evar.ExistVar
	}

	return true
}

func (t *Inferencer) solveExistVar(evar, typ ir.IrType) {
	if !evar.Is(ir.ExistVarType) {
		panic(fmt.Errorf("expected existential variable; got %v", evar))
	}

	existVar := t.existVars[evar.ExistVar]
	if existVar.solution != nil {
		panic(fmt.Errorf("%v is already solved with solution %s", evar, existVar.solution))
	}

	existVar.solution = &typ
	t.existVars[evar.ExistVar] = existVar
}

func (t *Inferencer) unifyImpl(left, right ir.IrType) error {
	switch {
	case t.isExistVarAssigned(left):
		left = t.existVarSolution(left)
		return t.unify(left, right)

	case t.isExistVarAssigned(right):
		right = t.existVarSolution(right)
		return t.unify(left, right)

	case t.isExistVarUnassigned(left) && t.canAssign(left, right):
		t.solveExistVar(left, right)
		return nil

	case t.isExistVarUnassigned(right) && t.canAssign(right, left):
		t.solveExistVar(right, left)
		return nil

	case left.Is(ir.AppType) && right.Is(ir.AppType):
		if err := t.unify(left.App.Fun, right.App.Fun); err != nil {
			return fmt.Errorf("mismatch in function types: %v", err)
		}
		if err := t.unify(left.App.Arg, right.App.Arg); err != nil {
			return fmt.Errorf("mismatch in argument types: %v", err)
		}
		return nil

	case left.Is(ir.ArrayType) && right.Is(ir.ArrayType):
		if err := t.unify(left.Array.ElemType, right.Array.ElemType); err != nil {
			return fmt.Errorf("mismatch in array element types: %v", err)
		}
		if left.Array.Size != right.Array.Size {
			return fmt.Errorf("expected array with %d elements; got %d elements", left.Array.Size, right.Array.Size)
		}
		return nil

	case right.Is(ir.ForallType):
		var err error
		t.context, _, right, err = t.context.AddFreshType(right)
		if err != nil {
			return err
		}

		return t.unify(left, right)

	case left.Is(ir.ForallType):
		c := left.Forall

		evar := t.newEvar()
		newBody := ir.SubstituteType(c.Type, ir.NewVarType(c.Var), evar)
		return t.unify(left, newBody)

	case left.Is(ir.FunType) && right.Is(ir.FunType):
		if err := t.unify(left.Fun.Arg, right.Fun.Arg); err != nil {
			return err
		}

		return t.unify(left.Fun.Ret, right.Fun.Ret)

	case left.Is(ir.NameType) && t.context.containsConstBind(left.Name) &&
		right.Is(ir.NameType) && t.context.containsConstBind(right.Name) &&
		left.Name == right.Name:
		return nil

	case left.Is(ir.NameType) && t.context.containsAliasBind(left.Name):
		bind, err := t.context.getAliasBind(left.Name)
		if err != nil {
			return err
		}
		return t.unify(bind.Alias.Type, right)

	case right.Is(ir.NameType) && t.context.containsAliasBind(right.Name):
		bind, err := t.context.getAliasBind(right.Name)
		if err != nil {
			return err
		}
		return t.unify(left, bind.Alias.Type)

	case left.Is(ir.StructType) && right.Is(ir.StructType) && sameFields(left, right):
		for i := range left.Fields() {
			if err := t.unify(left.Fields()[i].Type, right.Fields()[i].Type); err != nil {
				return err
			}
		}
		return nil

	case left.Is(ir.TupleType) && right.Is(ir.TupleType) && len(left.Tuple.Elems) == len(right.Tuple.Elems):
		for i := range left.Tuple.Elems {
			if err := t.unify(left.Tuple.Elems[i], right.Tuple.Elems[i]); err != nil {
				return err
			}
		}
		return nil

	case left.Is(ir.VariantType) && right.Is(ir.VariantType) && sameTags(left, right):
		for i := range left.Tags() {
			if err := t.unify(left.Tags()[i].Type, right.Tags()[i].Type); err != nil {
				return err
			}
		}
		return nil

	case left.Is(ir.VarType) && right.Is(ir.VarType) && left.Var == right.Var:
		return nil

	default:
		return fmt.Errorf("expected type %s (%s); got %s (%s)\n context: %s", left.Case, left, right.Case, right, t.context)
	}
}

func (t *Inferencer) unify(left, right ir.IrType) error {
	if err := isWellformedType(t.context, left); err != nil {
		return err
	}

	if err := isWellformedType(t.context, right); err != nil {
		return err
	}

	if err := t.unifyImpl(left, right); err != nil {
		return fmt.Errorf("%s\n  unifying %s and %s", err, left, right)
	}

	if glog.V(1) {
		keys := slices.Collect(maps.Keys(t.existVars))
		slices.Sort(keys)

		var b strings.Builder
		for _, name := range keys {
			existVar := t.existVars[name]
			if existVar.solution == nil {
				b.WriteString(fmt.Sprintf("\n    %s unsolved", ir.NewExistVarType(name)))
			} else {
				b.WriteString(fmt.Sprintf("\n    %s = %s", ir.NewExistVarType(name), *existVar.solution))
			}
		}

		glog.Infof("unify: %s\n  |- unify(%s, %s)%s", t.context, left, right, b.String())
	}

	return nil
}
