package bplparser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) parseCall() (ir.IrTerm, error) {
	if len(p.words) == 0 {
		return ir.IrTerm{}, fmt.Errorf("unexpected end of line")
	}

	tokens, err := parser.ParseTokens(p.words)
	if err != nil {
		return ir.IrTerm{}, err
	}

	var id string
	isSingle := false
	isFunction := false
	isIndexGet := false
	isIndexSet := false
	isOpUnary := false
	isWiden := false
	if len(tokens) > 0 && tokens[0].Case == parser.IDToken {
		id = tokens[0].Text
		isSingle = true

		if p.compiler.IsFunction(id) {
			isFunction = true
			tokens = tokens[1:]
		} else if id == "Index.get" {
			isIndexGet = true
			tokens = tokens[1:]
		} else if id == "Index.set" {
			isIndexSet = true
			tokens = tokens[1:]
		} else if strings.ContainsAny(id, "-") {
			isOpUnary = true
			tokens = tokens[1:]
		} else if id == "widen" {
			isWiden = true
			tokens = tokens[1:]
		} else {
			isSingle = false
		}
	}

	isOpBinary := false
	if !isSingle && len(tokens) > 1 && tokens[1].Case == parser.IDToken {
		id = tokens[1].Text

		if strings.ContainsAny(id, "+-*/") {
			isOpBinary = true
			tokens = append(tokens[0:1], tokens[2:]...)
		}
	}

	terms := make([]ir.IrTerm, len(tokens))
	for i := range tokens {
		terms[i] = ir.NewTokenTerm(tokens[i])
	}

	if isFunction {
		return ir.NewCallTerm(id, terms), nil
	}

	if isIndexGet {
		term, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		index, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		return ir.NewIndexGetTerm(term, index), nil
	}

	if isIndexSet {
		ret, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		index, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		arg, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		return ir.NewIndexSetTerm(ret, index, arg), nil
	}

	if isOpUnary {
		term, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		if err := parser.EOL(terms); err != nil {
			return ir.IrTerm{}, err
		}

		return ir.NewOpUnaryTerm(id, term), nil
	}

	if isOpBinary {
		left, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		right, terms, err := parser.ShiftID(terms)
		if err != nil {
			return ir.IrTerm{}, err
		}

		if err := parser.EOL(terms); err != nil {
			return ir.IrTerm{}, err
		}

		return ir.NewOpBinaryTerm(id, left, right), nil
	}

	if isWiden {
		return ir.NewWidenTerm(ir.NewTupleTerm(terms)), nil
	}

	return ir.NewTupleTerm(terms), nil
}

func (p *Parser) ParseCall() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseCall()
		return err
	})
	return
}
