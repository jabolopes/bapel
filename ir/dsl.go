package ir

import "github.com/jabolopes/bapel/parser"

func Call(id string, args ...IrTerm) IrTerm {
	return NewAppTermTerm(NewTokenTerm(parser.NewIDToken(id)), NewTupleTerm(args))
}

func CallPF(id string, types []IrType, args ...IrTerm) IrTerm {
	term := NewTokenTerm(parser.NewIDToken(id))
	for _, typ := range types {
		term = NewAppTypeTerm(term, typ)
	}
	return NewAppTermTerm(term, NewTupleTerm(args))
}

func Const(id string) IrType {
	return NewNameType(id)
}

func Forall(tvar string, typ IrType) IrType {
	return NewForallType(tvar, typ)
}

func Tvar(tvar string) IrType {
	return NewVarType(tvar)
}

func Terms(terms ...IrTerm) IrTerm {
	return NewTupleTerm(terms)
}

func Types(types ...IrType) IrType {
	return NewTupleType(types)
}

func TypesA(types ...IrType) []IrType {
	return append([]IrType{}, types...)
}

func ID(id string) IrTerm {
	return NewTokenTerm(parser.NewIDToken(id))
}

func Number(value int64) IrTerm {
	return NewTokenTerm(parser.NewNumberToken(value))
}
