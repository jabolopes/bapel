package parser

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	filename string
	errors   []string
}

func NewCustomErrorListener(filename string) *CustomErrorListener {
	return &CustomErrorListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		filename:             filename,
	}
}

func (c *CustomErrorListener) Errors() []string {
	return c.errors
}

func (c *CustomErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	if offendingToken, ok := offendingSymbol.(antlr.Token); ok {
		tokenType := offendingToken.GetTokenType()
		startLine := offendingToken.GetLine()

		text := offendingToken.GetText()
		newlines := strings.Count(text, "\n")
		endLine := startLine + newlines

		switch tokenType {
		case bapelParserUNTERMINATED_STRING_LITERAL:
			if startLine == endLine {
				msg = fmt.Sprintf("unterminated string (\" ... \") starting at line %d", startLine)
			} else {
				c.errors = append(c.errors, fmt.Sprintf("in %q in lines %d-%d: unterminated string (\" ... \") starting at line %d", c.filename, startLine, endLine, startLine))
				return
			}

		case bapelParserUNTERMINATED_RAW_STRING_LITERAL:
			if startLine == endLine {
				msg = fmt.Sprintf("unterminated raw string (` ... `) starting at line %d", startLine)
			} else {
				c.errors = append(c.errors, fmt.Sprintf("in %q in lines %d-%d: unterminated raw string (` ... `) starting at line %d", c.filename, startLine, endLine, startLine))
				return
			}

		case bapelParserUNTERMINATED_BLOCK_COMMENT:
			if startLine == endLine {
				msg = fmt.Sprintf("unterminated block comment (/* ... */) starting at line %d", startLine)
			} else {
				c.errors = append(c.errors, fmt.Sprintf("in %q in lines %d-%d: unterminated block comment (/* ... */) starting at line %d", c.filename, startLine, endLine, startLine))
				return
			}

		case bapelParserUNTERMINATED_RUNE_LITERAL:
			if startLine == endLine {
				msg = fmt.Sprintf("unterminated rune (' ... ') starting at line %d", startLine)
			} else {
				c.errors = append(c.errors, fmt.Sprintf("in %q in lines %d-%d: unterminated rune (' ... ') starting at line %d", c.filename, startLine, endLine, startLine))
				return
			}
		}
	}

	if strings.HasPrefix(msg, "token recognition error at: ") {
		raw := strings.TrimPrefix(msg, "token recognition error at: ")
		raw = strings.Trim(raw, "'")
		if strings.HasPrefix(raw, "\\x") {
			var val int
			if _, err := fmt.Sscanf(raw, "\\x%x", &val); err == nil {
				msg = fmt.Sprintf("unexpected token '\\x%02x' (%d) at line %d", val, val, line)
			}
		} else if len(raw) > 0 {
			char := raw[0]
			if char < 32 || char > 126 {
				msg = fmt.Sprintf("unexpected token '\\x%02x' (%d) at line %d", char, char, line)
			} else {
				msg = fmt.Sprintf("unexpected token '%s' (%d) at line %d", raw, int(char), line)
			}
		}
	}

	c.errors = append(c.errors, fmt.Sprintf("in %q in line %d: %s", c.filename, line, msg))
}
