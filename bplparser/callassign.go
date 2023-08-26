package bplparser

import (
	"fmt"
	"log"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

type IsFunction interface {
	IsFunction(string) bool
}

var Compiler IsFunction

// TODO: Merge ParseCall and ParseCall2.
func ParseCall2(args []string) (ir.IrTerm, []string, error) {
	orig := args

	tokens, err := parser.ParseTokens(args)
	if err != nil {
		return ir.IrTerm{}, orig, err
	}

	var id string
	isFunction := false
	isIndexGet := false
	isIndexSet := false
	isWiden := false
	if len(tokens) > 0 && tokens[0].Case == parser.IDToken {
		id = tokens[0].Text

		if Compiler.IsFunction(id) {
			isFunction = true
			tokens = tokens[1:]
		} else if id == "Index.get" {
			isIndexGet = true
			tokens = tokens[1:]
		} else if id == "Index.set" {
			isIndexSet = true
			tokens = tokens[1:]
		} else if id == "widen" {
			isWiden = true
			tokens = tokens[1:]
		}
	}

	terms := make([]ir.IrTerm, len(tokens))
	for i := range tokens {
		terms[i] = ir.NewTokenTerm(tokens[i])
	}

	if isFunction {
		return ir.NewCallTerm(id, terms), nil, nil
	}

	if isIndexGet {
		term, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, orig, err
		}

		index, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, orig, err
		}

		return ir.NewIndexGetTerm(term, index), nil, nil
	}

	if isIndexSet {
		ret, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, orig, err
		}

		index, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, orig, err
		}

		arg, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, orig, err
		}

		return ir.NewIndexSetTerm(ret, index, arg), nil, nil
	}

	if isWiden {
		return ir.NewWidenTerm(ir.NewTupleTerm(terms)), nil, nil
	}

	return ir.NewTupleTerm(terms), nil, nil
}

func ParseCall(args []string) ([]parser.Token, []string, error) {
	orig := args

	tokens, err := parser.ParseTokens(args)
	if err != nil {
		return nil, orig, err
	}

	if Compiler != nil {
		if len(tokens) > 0 && tokens[0].Case == parser.IDToken && Compiler.IsFunction(tokens[0].Text) {
			log.Printf("HERE %q is function", tokens[0].Text)
		}
	}

	return tokens, nil, nil
}

// ParseCallAssign parses call and assignment.
//
// Note that a call is an assignment without the '<-' and without any return
// values.
func ParseCallAssign(args []string) ([]ir.IrTerm, []ir.IrTerm, []string, error) {
	orig := args

	var rets []string
	for ; len(args) > 0; args = args[1:] {
		if args[0] == "<-" {
			args = args[1:]

			if len(args) == 0 {
				return nil, nil, orig, fmt.Errorf("expected at least 1 argument after token '<-'")
			}

			if len(rets) == 0 {
				return nil, nil, orig, fmt.Errorf("expected at least 1 return value before token '<-'")
			}

			argTokens, _, err := ParseCall(args)
			if err != nil {
				return nil, nil, orig, err
			}

			argTerms := make([]ir.IrTerm, len(argTokens))
			for i := range argTokens {
				argTerms[i] = ir.NewTokenTerm(argTokens[i])
			}

			retTokens, err := parser.ParseTokens(rets)
			if err != nil {
				return nil, nil, orig, err
			}

			retTerms := make([]ir.IrTerm, len(retTokens))
			for i := range retTokens {
				retTerms[i] = ir.NewTokenTerm(retTokens[i])
			}

			return argTerms, retTerms, nil, nil
		}

		rets = append(rets, args[0])
	}

	argTokens, _, err := ParseCall(orig)
	if err != nil {
		return nil, nil, orig, err
	}

	argTerms := make([]ir.IrTerm, len(argTokens))
	for i := range argTokens {
		argTerms[i] = ir.NewTokenTerm(argTokens[i])
	}

	return argTerms, nil, nil, nil
}
