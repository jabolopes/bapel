implements bapel.lex

imports {
  bapel.core
  bapel.stl
}

pub type FSM = struct {
  scanner ref::Ref Scanner,
  tokens Deque Token,
  lineNum i64,
  read String,
}
