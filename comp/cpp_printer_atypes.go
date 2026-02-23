package comp

import (
	"crypto/sha1"
	"fmt"
	"maps"
	"slices"

	"github.com/jabolopes/bapel/ir"
)

type anonymousType struct {
	// Generated name type for the originally anonymous type.
	nameType ir.IrType
	// Generated alias decl for the newly named type.
	decl ir.IrDecl
}

func (p *CppPrinter) hashType(typ ir.IrType) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(typ.String())))
}

func (p *CppPrinter) genNameType(typ ir.IrType) ir.IrType {
	hash := p.hashType(typ)
	name := fmt.Sprintf("__anonym_%s", hash)
	p.anonymousTypes[hash] = anonymousType{ir.NewNameType(name), ir.NewAliasDecl(name, ir.NewTypeKind(), typ, false /* export */)}
	return ir.NewNameType(name)
}

func (p *CppPrinter) recordAnonymousTypes(typ ir.IrType) ir.IrType {
	switch typ.Case {
	case ir.AppType:
		c := typ.App
		return ir.NewAppType(p.recordAnonymousTypes(c.Fun), p.recordAnonymousTypes(c.Arg))

	case ir.ArrayType:
		c := typ.Array
		return ir.NewArrayType(p.recordAnonymousTypes(c.ElemType), c.Size)

	case ir.ExistVarType:
		return typ

	case ir.ForallType:
		c := typ.Forall
		return ir.NewForallType(c.Var, c.Kind, p.recordAnonymousTypes(c.Type))

	case ir.FunType:
		c := typ.Fun
		return ir.NewFunctionType(p.recordAnonymousTypes(c.Arg), p.recordAnonymousTypes(c.Ret))

	case ir.LambdaType:
		c := typ.Lambda
		return ir.NewLambdaType(c.Var, c.Kind, p.recordAnonymousTypes(c.Type))

	case ir.NameType:
		return typ

	case ir.StructType:
		return p.genNameType(typ)

	case ir.TupleType:
		c := typ.Tuple
		elems := make([]ir.IrType, len(c.Elems))
		for i := range c.Elems {
			elems[i] = p.recordAnonymousTypes(c.Elems[i])
		}
		return ir.NewTupleType(elems)

	case ir.VariantType:
		c := typ.Variant
		tags := make([]ir.VariantTag, len(c.Tags))
		for i := range c.Tags {
			tags[i] = c.Tags[i]
			tags[i].Type = p.recordAnonymousTypes(tags[i].Type)
		}
		return ir.NewVariantType(tags)

	case ir.VarType:
		return typ

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (p *CppPrinter) recordAnonymousTypesFromUnit(unit *ir.IrUnit) error {
	for i := range unit.Decls {
		decl := &unit.Decls[i]
		if decl.Is(ir.TermDecl) {
			decl.Term.Type = p.recordAnonymousTypes(decl.Term.Type)
		}
	}

	for i := range unit.Functions {
		fun := &unit.Functions[i]
		// for i := range fun.Args {
		// 	arg := &fun.Args[i]
		// 	arg.Type = p.recordAnonymousTypes(arg.Type)
		// }
		fun.RetType = p.recordAnonymousTypes(fun.RetType)
	}

	{
		anoDecls := make([]ir.IrDecl, 0, len(p.anonymousTypes))
		for _, key := range slices.Sorted(maps.Keys(p.anonymousTypes)) {
			anoDecls = append(anoDecls, p.anonymousTypes[key].decl)
		}

		unit.Decls = append(unit.Decls, anoDecls...)

		var err error
		unit.Decls, err = ir.TopoSortDecls(unit.Decls)
		if err != nil {
			return err
		}
	}

	return nil
}
