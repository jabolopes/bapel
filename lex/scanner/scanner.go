package scanner

import "unicode/utf8"

type Scanner struct {
	file         string
	left         int
	right        int
	leftLineNum  int
	rightLineNum int
}

func (s *Scanner) restFile() string {
	if s.right <= len(s.file) {
		return s.file[s.right:]
	}
	return ""
}

func (s *Scanner) LineNum() int { return s.leftLineNum }

func (s *Scanner) PeekRune() (rune, bool) {
	file := s.restFile()

	r, _ := utf8.DecodeRuneInString(file)
	return r, r != utf8.RuneError
}

func (s *Scanner) PeekString(str string) bool {
	file := s.restFile()

	return len(str) <= len(file) && str == file[:len(str)]
}

func (s *Scanner) ReadRune() (rune, bool) {
	file := s.restFile()

	r, size := utf8.DecodeRuneInString(file)
	if r == utf8.RuneError {
		return r, false
	}

	s.right += size

	if r == '\n' {
		s.rightLineNum++
	}

	return r, true
}

func (s *Scanner) ReadString(str string) bool {
	if s.PeekString(str) {
		for range str {
			s.ReadRune()
		}
		return true
	}
	return false
}

func (s *Scanner) Current() string {
	if s.left <= len(s.file) {
		return s.file[s.left:s.right]
	}
	return ""
}

func (s *Scanner) Ignore() {
	s.left = s.right
	s.leftLineNum = s.rightLineNum
}

func New(file string) *Scanner {
	return &Scanner{file, 0, 0, 1 /* leftLineNum */, 1 /* rightLineNum */}
}
