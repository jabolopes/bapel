implements bapel.lex

imports {
  bapel.core
	bapel.stl
}

pub type TokenType = i64

pub type Token = struct {
  LineNum i64,
	Type TokenType,
	Value String,
}
