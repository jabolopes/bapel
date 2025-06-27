package stlc

import (
	"fmt"

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

type unifier struct {
	context   Context
	existVars map[string]existVar
}

func (t *unifier) isExistVarUnassigned(tvar ir.IrType) bool {
	if !tvar.Is(ir.VarType) {
		return false
	}
	existVar, ok := t.existVars[tvar.Var]
	return ok && existVar.solution == nil
}

func (t *unifier) isExistVarAssigned(tvar ir.IrType) bool {
	if !tvar.Is(ir.VarType) {
		return false
	}
	existVar, ok := t.existVars[tvar.Var]
	return ok && existVar.solution != nil
}

func (t *unifier) existVarSolution(tvar string) ir.IrType {
	existVar := t.existVars[tvar]
	solution := existVar.solution
	if solution.Is(ir.VarType) && t.isExistVarAssigned(*solution) {
		return t.existVarSolution(solution.Var)
	}
	return *solution
}

func (t *unifier) canAssign(tvar, typ ir.IrType) bool {
	ok, err := t.context.WellformedUnderTvar(tvar, typ)
	if err != nil {
		panic(err)
	}
	return ok
}

func (t *unifier) solveExistVar(tvar string, typ ir.IrType) {
	existVar := t.existVars[tvar]
	if existVar.solution != nil {
		panic(fmt.Errorf("existential variable %q is already solved with solution %s", tvar, existVar.solution))
	}

	existVar.solution = &typ
	t.existVars[tvar] = existVar
}

func (t *unifier) unifyImpl(left, right ir.IrType) error {
	switch {
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
		var tvar ir.IrType
		var err error
		t.context, tvar, right, err = t.context.AddFreshType(right)
		if err != nil {
			return err
		}

		t.existVars[tvar.Var] = existVar{}

		return t.unify(left, right)

	case left.Is(ir.ForallType) && right.Is(ir.ForallType):
		return fmt.Errorf("unhandled unification of %s and %s", left, right)

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

	case t.isExistVarUnassigned(left) && !right.Is(ir.ForallType) && t.canAssign(left, right):
		t.solveExistVar(left.Var, right)
		return nil

	case t.isExistVarAssigned(left):
		left = t.existVarSolution(left.Var)
		return t.unify(left, right)

	case t.isExistVarUnassigned(right) && !left.Is(ir.ForallType) && t.canAssign(right, left):
		t.solveExistVar(right.Var, left)
		return nil

	case right.Is(ir.VarType) && t.isExistVarAssigned(right):
		right = t.existVarSolution(right.Var)
		return t.unify(left, right)

	default:
		return fmt.Errorf("expected type %s (%s); got %s (%s)\n context: %s", left.Case, left, right.Case, right, t.context.StringNoDecls())
	}
}

func (t *unifier) unify(left, right ir.IrType) error {
	if err := isWellformedType(t.context, left); err != nil {
		return err
	}

	if err := isWellformedType(t.context, right); err != nil {
		return err
	}

	if err := t.unifyImpl(left, right); err != nil {
		return fmt.Errorf("%s\n  unifying %s and %s", err, left, right)
	}

	glog.V(1).Infof("unify: %s |- unify(%s, %s)", t.context, left, right)
	return nil
}

type applicativeUnification struct {
	// Whether all variables are solved.
	solved bool
	// Solutions for the type variables in the `forallType` passed to
	// the `unifyApplicativeSpine`. Solutions are in the same order as
	// the forall type variables in forallType.
	forallTypeEvars []existVar
	// Solution to the type variable in the `retType`. If `retType` was
	// already a type (non-nil), then this variable is automatically
	// solved.
	retTypeEvar existVar
}

func (t *unifier) unifyApplicativeSpine(forallType, argType ir.IrType, retType *ir.IrType) (applicativeUnification, error) {
	spine := applicativeUnification{solved: true}

	// Existential variables for the `forallType` solution.
	var forallTypeEvars []string

	// Existential variable for the `retType` solution (if any).
	var retTypeEvar *string

	var leftType ir.IrType
	{
		typ := forallType
		for typ.Is(ir.ForallType) {
			var tvar ir.IrType
			var err error
			t.context, tvar, typ, err = t.context.AddFreshType(typ)
			if err != nil {
				return applicativeUnification{}, err
			}

			t.existVars[tvar.Var] = existVar{}
			forallTypeEvars = append(forallTypeEvars, tvar.Var)
		}

		if !typ.Is(ir.FunType) {
			return applicativeUnification{}, nil
		}

		leftType = typ
	}

	var rightType ir.IrType
	{
		if retType == nil {
			evar := t.context.GenFreshVarType()

			var err error
			t.context, err = t.context.AddBind(NewTypeVarBind(evar.Var, ir.NewTypeKind()))
			if err != nil {
				return applicativeUnification{}, err
			}

			t.existVars[evar.Var] = existVar{}

			rightType = ir.NewFunctionType(argType, evar)
			retTypeEvar = &evar.Var
		} else {
			rightType = ir.NewFunctionType(argType, *retType)
			spine.retTypeEvar = existVar{retType}
		}
	}

	if err := t.unify(leftType, rightType); err != nil {
		return applicativeUnification{}, err
	}

	spine.forallTypeEvars = make([]existVar, 0, len(forallTypeEvars))
	for _, evar := range forallTypeEvars {
		existVar, ok := t.existVars[evar]
		if !ok {
			return applicativeUnification{}, fmt.Errorf("missing existential variable %q", evar)
		}

		spine.forallTypeEvars = append(spine.forallTypeEvars, existVar)

		if existVar.solution == nil {
			spine.solved = false
			glog.V(1).Infof("unifier: forall existential variable %s unsolved", evar)
		} else {
			glog.V(1).Infof("unifier: forall existential variable %s = %s", evar, *existVar.solution)
		}
	}

	if retTypeEvar != nil {
		existVar, ok := t.existVars[*retTypeEvar]
		if !ok {
			return applicativeUnification{}, fmt.Errorf("missing existential variable %q", *retTypeEvar)
		}

		spine.retTypeEvar = existVar

		if existVar.solution == nil {
			spine.solved = false
			glog.V(1).Infof("unifier: return type existential variable %s unsolved", *retTypeEvar)
		} else {
			glog.V(1).Infof("unifier: return type existential variable %s = %s", *retTypeEvar, *existVar.solution)
		}
	}

	return spine, nil
}

func newUnifier(context Context) *unifier {
	return &unifier{context, map[string]existVar{}}
}
