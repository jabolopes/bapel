package stlc

import "github.com/jabolopes/bapel/ir"

func allBlocksImpl(term ir.IrTerm, blocks *[]ir.IrTerm) {
	switch term.Case {
	case ir.AppTermTerm, ir.AppTypeTerm, ir.AssignTerm, ir.ConstTerm, ir.InjectionTerm, ir.IndexGetTerm, ir.IndexSetTerm, ir.LetTerm, ir.ReturnTerm, ir.TupleTerm, ir.VarTerm:
		break

	case ir.BlockTerm:
		c := term.Block

		*blocks = append(*blocks, term)
		for _, t := range c.Terms {
			allBlocksImpl(t, blocks)
		}

	case ir.IfTerm:
		c := term.If

		allBlocksImpl(c.Then, blocks)
		if c.Else != nil {
			allBlocksImpl(*c.Else, blocks)
		}
	}
}

func allReturns(term ir.IrTerm) []ir.IrTerm {
	var blocks []ir.IrTerm
	allBlocksImpl(term, &blocks)

	var returns []ir.IrTerm
	for _, block := range blocks {
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
	case ir.AppTermTerm, ir.AppTypeTerm, ir.AssignTerm, ir.ConstTerm, ir.InjectionTerm, ir.IndexGetTerm, ir.IndexSetTerm, ir.LetTerm, ir.ReturnTerm, ir.TupleTerm, ir.VarTerm:
		*last = append(*last, term)

	case ir.BlockTerm:
		c := term.Block

		if len(c.Terms) > 0 {
			lastTermsImpl(&c.Terms[len(c.Terms)-1], last)
		}

	case ir.IfTerm:
		c := term.If

		lastTermsImpl(&c.Then, last)
		if c.Else != nil {
			lastTermsImpl(c.Else, last)
		}
	}
}

func lastTerms(term *ir.IrTerm) []*ir.IrTerm {
	var last []*ir.IrTerm
	lastTermsImpl(term, &last)
	return last
}
