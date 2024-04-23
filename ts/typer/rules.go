package typer

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func Reduce(context Context) (retContext Context, retError error) {
	var bind Bind
	bind, retContext = context.Pop()
	retError = nil

	// Rule 1: Γ, a −> Γ
	if bind.Is(TypeBind) && bind.Type.Type.Is(VarType) {
		return
	}

	// Rule 2: Γ, â -> Γ
	if bind.Is(TypeBind) && bind.Type.Type.Is(ExistVarType) {
		return
	}

	// Rule 3: Γ, x : A -> Γ
	if bind.Is(TermBind) {
		return
	}

	if bind.Is(JudgeBind) && bind.Judge.Judge.Is(SubtypeJudge) {
		j := *bind.Judge.Judge.Subtype

		// Rule 4: Γ ⊩ 1 ≤ 1 -> Γ
		if j.Left.Is(NameType) && j.Right.Is(NameType) && *j.Left.Name == *j.Right.Name {
			return
		}

		// Rule 5: Γ ⊩ a ≤ a −> Γ
		if j.Left.Is(VarType) && j.Right.Is(VarType) && *j.Left.Var == *j.Right.Var {
			return
		}

		// Rule 6: Γ ⊩ â ≤ â −> Γ
		if j.Left.Is(ExistVarType) && j.Right.Is(ExistVarType) && *j.Left.ExistVar == *j.Right.ExistVar {
			return
		}

		// Rule 7: Γ ⊩ A1 → A2 ≤ B1 → B2 -> Γ ⊩ A2 ≤ B2 ⊩ B1 ≤ A1
		if j.Left.Is(FunType) && j.Right.Is(FunType) {
			c1 := *j.Left.Fun
			c2 := *j.Right.Fun
			A1 := c1.Arg
			A2 := c1.Ret
			B1 := c2.Arg
			B2 := c2.Ret

			retContext = retContext.AddJudge(NewSubtypeJudge(A2, B2))
			retContext = retContext.AddJudge(NewSubtypeJudge(B1, A1))
			return
		}

		// Rule 8: Γ ⊩ ∀a. A ≤ B -> Γ, â ⊩ [â/a]A ≤ B		when B != ∀a. B′
		if j.Left.Is(ForallType) && !j.Right.Is(ForallType) {
			a := j.Left.Forall.Var
			A := j.Left.Forall.Type
			B := j.Right

			var â Type
			retContext, â = retContext.AddFreshExistVarType()

			retContext.AddJudge(NewSubtypeJudge(substituteVar(A, a, â), B))
			return
		}

		// Rule 9: Γ ⊩ A ≤ ∀b. B -> Γ, b ⊩ A ≤ B
		if j.Right.Is(ForallType) {
			A := j.Left
			b := j.Right.Forall.Var
			B := j.Right.Forall.Type

			// Substitution ensures type variable freshness.
			var tvar Type
			retContext, tvar = retContext.AddFreshVarType()
			B = substituteVar(B, b, tvar)

			retContext = retContext.AddJudge(NewSubtypeJudge(A, B))
			return
		}

		// Rule 10: Γ[â] ⊩ â ≤ A → B -> [â1 → â/â](Γ[â1, â2] ⊩ â1 → â2 ≤ A → B)	when â not-member FV (A) ∪ FV (B
		if j.Left.Is(ExistVarType) && j.Right.Is(FunType) &&
			retContext.ContainsExistVarType(*j.Left.ExistVar) &&
			!ContainsFreeVar(j.Right.Fun.Arg, j.Left.ExistVar.Var) &&
			!ContainsFreeVar(j.Right.Fun.Ret, j.Left.ExistVar.Var) {
			â := *j.Left.ExistVar

			var â1, â2 Type
			retContext, â1 = retContext.GenFreshExistVarType()
			retContext, â2 = retContext.GenFreshExistVarType()
			retContext = retContext.ReplaceExistVar(â, []Type{â1, â2})
			funType := NewFunType(â1, â2)
			retContext = retContext.AddJudge(NewSubtypeJudge(funType, j.Right))
			retContext = substituteVarInContext(retContext, â.Var, funType)
			return
		}

		// Rule 11: Γ[â] ⊩ A → B ≤ â -> [â1 → â2/â](Γ[â1, â2] ⊩ A → B ≤ â1 → â2 )	when â not-member FV (A) ∪ FV (B)
		if j.Left.Is(FunType) && j.Right.Is(ExistVarType) &&
			retContext.ContainsExistVarType(*j.Right.ExistVar) &&
			!ContainsFreeVar(j.Left.Fun.Arg, j.Right.ExistVar.Var) &&
			!ContainsFreeVar(j.Left.Fun.Ret, j.Right.ExistVar.Var) {
			â := *j.Right.ExistVar

			var â1, â2 Type
			retContext, â1 = retContext.GenFreshExistVarType()
			retContext, â2 = retContext.GenFreshExistVarType()
			retContext = retContext.ReplaceExistVar(â, []Type{â1, â2})
			funType := NewFunType(â1, â2)
			retContext = retContext.AddJudge(NewSubtypeJudge(j.Left, funType))
			retContext = substituteVarInContext(retContext, â.Var, funType)
			return
		}

		// Rule 12: Γ[â][β^] ⊩ â ≤ β^ -> [â/β^](Γ[â][])
		if j.Left.Is(ExistVarType) && j.Right.Is(ExistVarType) &&
			retContext.ContainsExistVarTypesInOrder(*j.Left.ExistVar, *j.Right.ExistVar) {
			â := j.Left
			bhat := *j.Right.ExistVar

			retContext = retContext.ReplaceExistVar(bhat, nil /* replacements */)
			retContext = substituteVarInContext(retContext, bhat.Var, â)
			return
		}

		// Rule 13: Γ[â][β^] ⊩ β^ ≤ â -> [â/β^](Γ[â][])
		if j.Left.Is(ExistVarType) && j.Right.Is(ExistVarType) &&
			retContext.ContainsExistVarTypesInOrder(*j.Right.ExistVar, *j.Left.ExistVar) {
			bhat := *j.Left.ExistVar
			â := j.Right

			retContext = retContext.ReplaceExistVar(bhat, nil /* replacements */)
			retContext = substituteVarInContext(retContext, bhat.Var, â)
			return
		}

		// Rule 14: Γ[a][β^] ⊩ a ≤ β^ -> [a/β^](Γ[a][])
		if j.Left.Is(VarType) && j.Right.Is(ExistVarType) &&
			retContext.ContainsVarsInOrder(*j.Left.Var, *j.Right.ExistVar) {
			a := j.Left
			bhat := *j.Right.ExistVar

			retContext = retContext.ReplaceExistVar(bhat, nil /* replacements */)
			retContext = substituteVarInContext(retContext, bhat.Var, a)
			return
		}

		// Rule 15: Γ[a][β^] ⊩ β^ ≤ a -> [a/β^](Γ[a][])
		if j.Left.Is(ExistVarType) && j.Right.Is(VarType) &&
			retContext.ContainsVarsInOrder(*j.Right.Var, *j.Left.ExistVar) {
			bhat := *j.Left.ExistVar
			a := j.Right

			retContext = retContext.ReplaceExistVar(bhat, nil /* replacements */)
			retContext = substituteVarInContext(retContext, bhat.Var, a)
			return
		}

		// Rule 16: Γ[β^] ⊩ 1 ≤ β^ -> [1/β^](Γ[])
		if j.Left.Is(NameType) && j.Right.Is(ExistVarType) &&
			retContext.ContainsExistVarType(*j.Right.ExistVar) {
			name := j.Left
			bhat := *j.Right.ExistVar

			retContext = retContext.ReplaceExistVar(bhat, nil /* replacements */)
			retContext = substituteVarInContext(retContext, bhat.Var, name)
			return
		}

		// Rule 17: Γ[β^] ⊩ β^ ≤ 1 -> [1/β^](Γ[])
		if j.Left.Is(ExistVarType) && j.Right.Is(NameType) &&
			retContext.ContainsExistVarType(*j.Left.ExistVar) {
			bhat := *j.Left.ExistVar
			name := j.Right

			retContext = retContext.ReplaceExistVar(bhat, nil /* replacements */)
			retContext = substituteVarInContext(retContext, bhat.Var, name)
			return
		}
	}

	if bind.Is(JudgeBind) && bind.Judge.Judge.Is(CheckJudge) {
		j := *bind.Judge.Judge.Check

		// Rule 18: Γ ⊩ e ⇐ B -> Γ ⊩ e ⇒a a ≤ B	when e != λx . e′ and B != ∀a. B
		//
		// TODO: Eventually, add !j.Term.Is(LambdaTerm).
		if !j.Type.Is(ForallType) {
			var tvar string
			retContext, tvar = retContext.GenFreshID()
			retContext = retContext.AddJudge(
				NewInferenceJudge(
					j.Term,
					tvar,
					NewSubtypeJudge(NewVarType(tvar), j.Type)))
			return
		}

		// Rule 19: Γ ⊩ e ⇐ ∀a. A -> Γ, a ⊩ e ⇐ A
		if j.Type.Is(ForallType) {
			var tvar Type
			retContext, tvar = retContext.AddFreshVarType()
			a := substituteVar(j.Type.Forall.Type, j.Type.Forall.Var, tvar)

			retContext = retContext.AddJudge(NewCheckJudge(j.Term, a))
			return
		}

		// Rule 20: Γ ⊩ λx . e ⇐ A → B -> Γ, x : A ⊩ e ⇐ B
		//
		// TODO: Implement.

		// Rule 21: Γ[â] ⊩ λx . e ⇐ â -> [â1 → â2/â](Γ[â1, â2], x : â1 ⊩ e ⇐ â2 )
		//
		// TODO: Implement.
	}

	if bind.Is(JudgeBind) && bind.Judge.Judge.Is(InferenceJudge) {
		j := *bind.Judge.Judge.Inference

		// Rule 22: Γ ⊩ x ⇒a ω -> Γ ⊩ [A/a]ω	when (x : A) ∈ Γ
		if j.Term.Is(ir.TokenTerm) && j.Term.Token.Is(parser.IDToken) &&
			context.ContainsTermBind(j.Term.Token.Text) {
			a := context.GetTermType(j.Term.Token.Text)
			retContext = retContext.AddJudge(substituteVarInJudge(j.Judge, j.Var, a))
			return
		}

		// Rule 23: Γ ⊩ (e : A) ⇒a ω -> Γ ⊩ [A/a]ω ⊩ e ⇐ A
		//
		// TODO: Implement if type annotations are needed.

		// Rule 24: Γ ⊩ () ⇒a ω -> Γ ⊩ [1/a]ω
		//
		// TODO: Figure out how to implement this for bapel's syntax.

		// Rule 25: Γ ⊩ λx . e ⇒a ω −→25 Γ, â, β^ ⊩ [â → β^/a]ω, x : â ⊩ e ⇐ β^
		//
		// TODO: Implement if lambda abstractions are supported.

		// Rule 26: Γ ⊩ e1 e2 ⇒a ω -> Γ ⊩ e1 ⇒b (b • e2 ⇒⇒a ω)
		if j.Term.Is(ir.CallTerm) {
			c := j.Term.Call
			e1 := ir.NewTokenTerm(parser.NewIDToken(c.ID))
			e2 := c.Arg

			a := j.Var

			var b Type
			retContext, b = retContext.GenFreshVarType()

			retContext = retContext.AddJudge(
				NewInferenceJudge(
					e1,
					b.Var.Var,
					NewApplicationInferenceJudge(b, e2, a, j.Judge)))
			return
		}
	}

	if bind.Is(JudgeBind) && bind.Judge.Judge.Is(ApplicationInferenceJudge) {
		j := *bind.Judge.Judge.ApplicationInference

		// Rule 27: Γ ⊩ ∀a. A • e ⇒⇒a ω -> Γ, â ⊩ [â/a]A • e ⇒⇒a ω
		if j.Type.Is(ForallType) {
			c := *j.Type.Forall
			e := j.Term
			a := j.Var
			w := j.Judge

			var â Type
			retContext, â = retContext.AddFreshExistVarType()

			retContext = retContext.AddJudge(
				NewApplicationInferenceJudge(
					substituteVar(c.Type, a, â),
					e,
					a,
					w))
			return
		}

		// Rule 28: Γ ⊩ A → C • e ⇒⇒a ω -> Γ ⊩ [C/a]ω ⊩ e ⇐ A
		if j.Type.Is(FunType) {
			c := *j.Type.Fun
			A := c.Arg
			C := c.Ret
			e := j.Term
			a := j.Var
			w := j.Judge

			retContext = retContext.AddJudge(substituteVarInJudge(w, a, C))
			retContext = retContext.AddJudge(NewCheckJudge(e, A))
			return
		}

		// Rule 29: Γ[â] ⊩ â • e ⇒⇒a ω -> [â1 → â2/Dα](Γ[â1, â2] ⊩ â1 → â2 • e ⇒⇒a ω)
		if j.Type.Is(ExistVarType) && retContext.ContainsExistVarType(*j.Type.ExistVar) {
			â := *j.Type.ExistVar
			e := j.Term
			a := j.Var
			w := j.Judge

			var â1, â2 Type
			retContext, â1 = retContext.GenFreshExistVarType()
			retContext, â2 = retContext.GenFreshExistVarType()
			retContext = retContext.ReplaceExistVar(â, []Type{â1, â2})
			funType := NewFunType(â1, â2)
			retContext = retContext.AddJudge(NewApplicationInferenceJudge(funType, e, a, w))
			retContext = substituteVarInContext(retContext, â.Var, funType)
		}
	}

	retError = fmt.Errorf("failed to reduce context")
	return
}
