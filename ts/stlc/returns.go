package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

// allBlocksImpl returns all control blocks that can contain a
// `return` statement.
//
// For example, an `if` term can contain a `return` statement:
//
//	if ... { return ... }
//
// Also, a let term of a match term cannot contain a `return`
// statement.
//
// Also, a lambda term can contain a `return` statement but the lambda
// term defines a separate block of its own, so the return inside the
// lambda does not return from the outer function.
func allBlocksImpl(term ir.IrTerm, blocks *[]ir.IrTerm) {
	switch term.Case {
	case ir.AppTermTerm, ir.AppTypeTerm, ir.AssignTerm, ir.ConstTerm, ir.InjectionTerm,
		ir.LambdaTerm, ir.LetTerm, ir.MatchTerm, ir.ProjectionTerm, ir.ReturnTerm, ir.SetTerm,
		ir.StructTerm, ir.TupleTerm, ir.TypeAbsTerm, ir.VarTerm:
		break

	case ir.BlockTerm:
		c := term.Block

		*blocks = append(*blocks, term)
		for _, t := range c.Terms {
			allBlocksImpl(t, blocks)
		}

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func allReturns(term ir.IrTerm) []ir.IrTerm {
	var blocks []ir.IrTerm
	allBlocksImpl(term, &blocks)

	var returns []ir.IrTerm
	for _, block := range blocks {
		if !block.Is(ir.BlockTerm) {
			panic(fmt.Errorf("expected block term; got %s", block))
		}

		for _, t := range block.Block.Terms {
			if t.Is(ir.ReturnTerm) {
				returns = append(returns, t)
			}
		}
	}

	return returns
}

func lastTermsImpl(term *ir.IrTerm, last *[]*ir.IrTerm) {
	switch term.Case {
	case ir.AppTermTerm, ir.AppTypeTerm, ir.AssignTerm, ir.ConstTerm, ir.InjectionTerm,
		ir.LambdaTerm, ir.LetTerm, ir.MatchTerm, ir.ProjectionTerm, ir.ReturnTerm, ir.SetTerm,
		ir.StructTerm, ir.TupleTerm, ir.TypeAbsTerm, ir.VarTerm:
		*last = append(*last, term)

	case ir.BlockTerm:
		c := term.Block

		if len(c.Terms) > 0 {
			lastTermsImpl(&c.Terms[len(c.Terms)-1], last)
		}

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func lastTerms(term *ir.IrTerm) []*ir.IrTerm {
	var last []*ir.IrTerm
	lastTermsImpl(term, &last)
	return last
}
