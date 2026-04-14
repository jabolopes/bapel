implements bapel.lex

imports {
  bapel.core
  bapel.stl
}

pub type Rune = i8

pub type Scanner = struct {
  file StringView,
  left i64,
  right i64,
  leftLineNum i64,
  rightLineNum i64,
}

pub fn newScanner(file: StringView) -> Scanner {
  struct{
    file = file,
    left = 0,
    right = 0,
    leftLineNum = 1,
    rightLineNum = 1,
  }
}

pub fn peekRune(scanner: Scanner) -> (std::optional Rune) {
  if scanner.right < StringView_::size scanner.file {
    let r = StringView_::at (scanner.file, scanner.right);
    return std::make_optional [Rune] r
  }

  none [Rune] ()
}

pub fn peekString(scanner: Scanner, str: StringView) -> bool {
  let peek = StringView_::substr (scanner.file, scanner.right, StringView_::size str);
  str == peek
}

pub fn readRune(scanner: &Scanner) -> (std::optional Rune) {
  let r = peekRune $scanner;
  if !has_value r {
    return r
  }

  let s = $scanner;
  s <- set s { right = s.right + 1 };

  if get_value r == '\n' {
    s <- set s { rightLineNum = s.rightLineNum + 1 }
  }

  ptr::set (scanner, s);

  r
}

pub fn readString(scanner: &Scanner, str: StringView) -> bool {
  if peekString ($scanner, str) {
    return true
  }
  false
}

pub fn current(scanner: Scanner) -> StringView {
  StringView_::substr (scanner.file, scanner.left, scanner.right)
}

pub fn ignore(scanner: &Scanner) -> () {
  let s = $scanner;
  s <- set s {
    left = s.right,
    leftLineNum = s.rightLineNum,
  };
  ptr::set (scanner, s)
}
