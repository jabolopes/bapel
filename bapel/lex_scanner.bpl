implements bapel.lex

imports {
  bapel.core
}

pub type Rune = i8

pub type Scanner = struct {
	file std::string_view,
	lineNum i64,
}

pub fn newScanner(file: std::string_view) -> Scanner {
  struct{file = file, lineNum = 1}
}

pub fn peekRune(scanner: Scanner) -> (std::optional Rune) {
	let value = none [Rune] ();

	if string_view::empty scanner.file {
		return value
	}

	std::make_optional [Rune] (string_view::front scanner.file)
}

pub fn peekRunes(scanner: Scanner, n: i64) -> std::string_view {
	string_view::substr (scanner.file, 0, n)
}
