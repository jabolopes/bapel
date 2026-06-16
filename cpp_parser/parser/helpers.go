package parser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"golang.org/x/exp/constraints"
)

type Token struct {
	Pos  ir.Pos
	Text string
}

func parseNumber[T constraints.Integer](arg string) (T, error) {
	var value T

	if strings.HasPrefix(arg, "0x") {
		// Hexadecimal
		_, err := fmt.Sscanf(arg, "0x%x", &value)

		return value, err
	}

	// Decimal.
	_, err := fmt.Sscanf(arg, "%d", &value)
	return value, err
}

func parseFloat[T constraints.Integer](arg string) (T, T, error) {
	splits := strings.Split(arg, ".")
	if len(splits) != 2 {
		return 0, 0, fmt.Errorf("invalid floating point number %q", arg)
	}

	integer, err := parseNumber[T](splits[0])
	if err != nil {
		return 0, 0, err
	}

	decimal, err := parseNumber[T](splits[1])
	if err != nil {
		return 0, 0, err
	}

	return integer, decimal, nil
}

func makePos(pos1, pos2 ir.Pos) ir.Pos {
	return ir.NewRangePos(pos1.Filename, pos1.BeginLineNum, pos2.EndLineNum)
}

func makePos2(args []any) ir.Pos {
	return makePos(args[0].(Token).Pos, args[len(args)-1].(Token).Pos)
}

func newUnaryOpExpr(id ast.Expr, typeArgs []ir.IrType, expr ast.Expr) (r ast.Expr) {
	pos := makePos(id.Pos, expr.Pos)

	if id.Is(ast.VarExpr) && id.Var.ID == "-" {
		// 0 - $expr
		return ast.Call(pos, id, typeArgs, ast.NewConstExpr(ir.NewIntLiteral(pos, 0)), expr)
	}

	return ast.Call(pos, id, typeArgs, expr)
}

func newBinOpExpr(id ast.Expr, typeArgs []ir.IrType, t1, t2 ast.Expr) ast.Expr {
	return ast.Call(makePos(t1.Pos, t2.Pos), id, typeArgs, t1, t2)
}

func newAliasDecl(id ast.ID, kind ir.IrKind, typ ir.IrType, export bool) ir.IrDecl {
	decl := ir.NewAliasDecl(id.Value, kind, typ, export)
	decl.Pos = makePos(id.Pos, typ.Pos)
	return decl
}

func newTermDecl(id ast.ID, typ ir.IrType, export bool) ir.IrDecl {
	decl := ir.NewTermDecl(id.Value, typ, export)
	decl.Pos = makePos(id.Pos, typ.Pos)
	return decl
}

func newNameDecl(id ast.ID, kind ir.IrKind, export bool) ir.IrDecl {
	decl := ir.NewNameDecl(id.Value, kind, export)
	decl.Pos = id.Pos
	return decl
}

func newQuantifiedType(typ ir.IrType) ir.IrType {
	quantified := ir.QuantifyType(typ)
	quantified.Pos = typ.Pos
	return quantified
}

func newForallType(pos ir.Pos, tvars []ir.VarKind, subType ir.IrType) ir.IrType {
	forall := ir.ForallVars(tvars, subType)
	forall.Pos = pos
	return forall
}

func newFunctionType(arg, ret ir.IrType) ir.IrType {
	typ := ir.NewFunctionType(arg, ret)
	typ.Pos = makePos(arg.Pos, ret.Pos)
	return typ
}

func newAppType(arg, ret ir.IrType) ir.IrType {
	typ := ir.NewAppType(arg, ret)
	typ.Pos = makePos(arg.Pos, ret.Pos)
	return typ
}

func newVarType(id ast.ID) ir.IrType {
	typ := ir.NewVarType(id.Value)
	typ.Pos = id.Pos
	return typ
}

func newNameType(id ast.ID) ir.IrType {
	typ := ir.NewNameType(id.Value)
	typ.Pos = id.Pos
	return typ
}

func newArrayType(pos ir.Pos, elemType ir.IrType, length int) ir.IrType {
	typ := ir.NewArrayType(elemType, length)
	typ.Pos = pos
	return typ
}

func newStructType(pos ir.Pos, fields []ir.StructField) ir.IrType {
	typ := ir.NewStructType(fields)
	typ.Pos = pos
	return typ
}

func newTupleType(pos ir.Pos, values []ir.IrType) ir.IrType {
	typ := ir.NewTupleType(values)
	typ.Pos = pos
	return typ
}

func newVariantType(pos ir.Pos, fields []ir.VariantTag) ir.IrType {
	typ := ir.NewVariantType(fields)
	typ.Pos = pos
	return typ
}

func newAssignExpr(arg, ret ast.Expr) ast.Expr {
	return ast.NewAssignExpr(makePos(arg.Pos, ret.Pos), arg, ret)
}

func newIfExpr(pos ir.Pos, condition, then ast.Expr, elseExpr *ast.Expr) ast.Expr {
	if elseExpr == nil {
		return ast.NewAppTermExpr(
			pos,
			ast.NewVarExpr(ast.NewID("ifthen", pos)),
			ast.NewTupleExpr(pos, []ast.Expr{condition, then}))
	}

	return ast.NewAppTermExpr(
		pos,
		ast.NewVarExpr(ast.NewID("ifelse", pos)),
		ast.NewTupleExpr(pos, []ast.Expr{condition, then, *elseExpr}))
}


func newMatchArm(tag, arg ast.ID, body ast.Expr) ast.MatchArm {
	return ast.NewMatchArm(tag.Value, arg.Value, body)
}

func newLetExpr(varName ast.ID, varType *ir.IrType, value ast.Expr) ast.Expr {
	return ast.NewLetExpr(makePos(varName.Pos, value.Pos), varName.Value, varType, value)
}

func newAppTypeExpr(expr ast.Expr, types []ir.IrType) ast.Expr {
	pos := expr.Pos
	for _, typ := range types {
		expr = ast.NewAppTypeExpr(pos, expr, typ)
	}
	return expr
}

func newModuleID(tokens []Token) ir.ModuleID {
	var b strings.Builder
	ir.Interleave(tokens, func() { b.WriteString(ir.ModuleIDSeparator) }, func(_ int, token Token) {
		b.WriteString(token.Text)
	})
	return ir.NewModuleID(b.String(), makePos(tokens[0].Pos, tokens[len(tokens)-1].Pos))
}

func newNumberLiteral(token Token) ir.IrLiteral {
	if strings.Contains(token.Text, ".") {
		integer, decimal, err := parseFloat[int64](token.Text)
		if err != nil {
			panic(fmt.Errorf("expected floating point literal; got %q", token.Text))
		}

		return ir.NewFloatLiteral(token.Pos, integer, decimal)
	}

	value, err := parseNumber[int64](token.Text)
	if err != nil {
		panic(fmt.Errorf("expected integer; got %q", token.Text))
	}

	return ir.NewIntLiteral(token.Pos, value)
}

func newRuneLiteral(token Token) ir.IrLiteral {
	text := token.Text

	if !strings.HasPrefix(text, `'`) || !strings.HasSuffix(text, `'`) {
		panic(fmt.Errorf(`expected rune delimited with '\''; got %q`, token.Text))
	}

	text = strings.TrimPrefix(text, `'`)
	text = strings.TrimSuffix(text, `'`)
	return ir.NewRuneLiteral(token.Pos, text)
}

func newStringLiteral(token Token) ir.IrLiteral {
	text := token.Text

	if !strings.HasPrefix(text, `"`) || !strings.HasSuffix(text, `"`) {
		panic(fmt.Errorf(`expected string delimited with '"'; got %q`, token.Text))
	}

	text = strings.TrimPrefix(text, `"`)
	text = strings.TrimSuffix(text, `"`)
	return ir.NewStrLiteral(token.Pos, text)
}
