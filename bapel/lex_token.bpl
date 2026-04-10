implements bapel.lex

imports {
  bapel.core
	bapel.stl
}

type TokenType = i64

type Token = struct {
  LineNum i64,
	Type TokenType,
	Value String,
}
