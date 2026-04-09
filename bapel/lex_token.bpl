implements bapel.lex

imports {
  bapel.core
}

type TokenType = i64

type Token = struct {
  LineNum i64,
	Type TokenType,
	Value std::string,
}
