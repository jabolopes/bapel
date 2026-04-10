implements bapel.lex

imports {
  bapel.core
	bapel.stl
}

type FSM = struct {
  scanner ref::Ref Scanner,
	tokens Deque Token,
	lineNum i64,
	read std::string,
}
