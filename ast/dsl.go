package ast

import "github.com/jabolopes/bapel/ir"

func Call(pos ir.Pos, id Expr, types []ir.IrType, args ...Expr) Expr {
	expr := id
	for _, typ := range types {
		expr = NewAppTypeExpr(pos, expr, typ)
	}
	if len(args) == 0 {
		return expr
	}
	return NewAppTermExpr(pos, expr, NewTupleExpr(pos, args))
}

func Lambda(pos ir.Pos, tvars []ir.VarKind, args []ir.FunctionArg, body Expr) Expr {
	if len(tvars) > 0 {
		return NewTypeAbsExpr(pos, tvars[0], Lambda(pos, tvars[1:], args, body))
	}

	if len(args) > 0 {
		return NewLambdaExpr(pos, args[0], Lambda(pos, nil, args[1:], body))
	}

	return body
}
