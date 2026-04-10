implements bapel.lex

imports {
  bapel.core
  bapel.stl
}

pub type Rune = i8

pub type Scanner = struct {
  file StringView,
  lineNum i64,
}

pub fn newScanner(file: StringView) -> Scanner {
  struct{file = file, lineNum = 1}
}

pub fn peekRune(scanner: Scanner) -> (std::optional Rune) {
  let value = none [Rune] ();

  if StringView_::empty scanner.file {
    return value
  }

  std::make_optional [Rune] (StringView_::front scanner.file)
}

pub fn peekRunes(scanner: Scanner, n: i64) -> StringView {
  StringView_::substr (scanner.file, 0, n)
}

pub fn readRune(scanner: &Scanner) -> (std::optional Rune) {
  let r = peekRune $scanner;
  if !has_value r {
    return r
  }

  let s = $scanner;
  s <- set s {
    file = StringView_::substr (s.file, 1, StringView_::size s.file)
  };

  if get_value r == '\n' {
    s <- set s { lineNum = s.lineNum + 1 }
  }

  ptr::set (scanner, s);

  r
}
