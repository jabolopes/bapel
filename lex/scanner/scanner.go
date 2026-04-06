package scanner

type Scanner struct {
	file    []rune
	lineNum int
}

func (s *Scanner) LineNum() int { return s.lineNum }

func (s *Scanner) PeekRune() (rune, bool) {
	var r rune

	if len(s.file) == 0 {
		return r, false
	}

	return s.file[0], true
}

func (s *Scanner) PeekRunes(n int) ([]rune, bool) {
	if len(s.file) < n {
		return nil, false
	}

	return s.file[:n], true
}

func (s *Scanner) ReadRune() (rune, bool) {
	r, ok := s.PeekRune()
	if !ok {
		return r, false
	}

	s.file = s.file[1:]

	if r == '\n' {
		s.lineNum++
	}

	return r, true
}

func New(file []rune) *Scanner {
	return &Scanner{file, 1 /* lineNum */}
}
