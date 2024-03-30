package bplparser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/slices"
)

func (p *Parser) parseCallImpl() (ir.IrTerm, error) {
	tokens, err := parser.ParseTokens(p.shiftTillEOL())
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

		if id == "-" {
			isOpUnary = true
			tokens = tokens[1:]
		} else if p.compiler.IsFunction(id) {
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

	// Parse types passing into the call.
	//
	// Example:
	//   f [MyType]
	var types []ir.IrType
	i := slices.IndexFunc(tokens, func(token parser.Token) bool {
		return token.Text == "["
	})
	j := slices.IndexFunc(tokens, func(token parser.Token) bool {
		return token.Text == "]"
	})
	if i < j {
		for _, token := range tokens[i+1 : j] {
			// TODO: Deduplicate with 'Parser.parseTypeImpl()'.

			var typ ir.IrType

			var r rune
			for _, r = range token.Text {
				break
			}

			if r == '\'' {
				typ = ir.NewVarType(strings.TrimPrefix(token.Text, "'"))
			} else if unicode.IsLetter(r) {
				typ = ir.NewNameType(token.Text)
			} else {
				return ir.IrTerm{}, fmt.Errorf("expected type variable or type name; got %q", token)
			}

			types = append(types, typ)
		}

		tokens = slices.Delete(tokens, i, j+1)
	}

	terms := make([]ir.IrTerm, len(tokens))
	for i := range tokens {
		terms[i] = ir.NewTokenTerm(tokens[i])
	}

	if isFunction {
		return ir.NewCallTerm(id, types, ir.NewTupleTerm(terms)), nil
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

		if err := parser.EOL(terms); err != nil {
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

		if err := parser.EOL(terms); err != nil {
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

		return ir.NewCallTerm(id, nil /* types */, ir.NewTupleTerm([]ir.IrTerm{ir.NewTokenTerm(parser.NewNumberToken(0)), term})), nil
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

		return ir.NewCallTerm(id, nil /* types */, ir.NewTupleTerm([]ir.IrTerm{left, right})), nil
	}

	if isWiden {
		return ir.NewWidenTerm(ir.NewTupleTerm(terms)), nil
	}

	return ir.NewTupleTerm(terms), nil
}

func (p *Parser) parseCall() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseCallImpl()
		return err
	})
	return
}
