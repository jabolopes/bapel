package scanner

import (
	"bufio"
	"io"

	"github.com/emirpasic/gods/v2/lists"
	"github.com/emirpasic/gods/v2/lists/arraylist"
)

type Scanner struct {
	reader  *bufio.Reader
	err     error
	lineNum int
	read    lists.List[rune]
}

func (s *Scanner) Err() error   { return s.err }
func (s *Scanner) LineNum() int { return s.lineNum }

func (s *Scanner) PeekRune() (rune, bool) {
	if r, ok := s.read.Get(0); ok {
		return r, true
	}

	if s.err != nil {
		return 0, false
	}

	r, _, err := s.reader.ReadRune()
	if err != nil {
		s.err = err
		return r, false
	}

	s.read.Add(r)

	return r, true
}

func (s *Scanner) PeekRunes(n int) ([]rune, bool) {
	rs := make([]rune, 0, n)
	for i := 0; i < n; i++ {
		if r, ok := s.read.Get(i); ok {
			rs = append(rs, r)
			continue
		}

		if s.err != nil {
			return nil, false
		}

		r, _, err := s.reader.ReadRune()
		if err != nil {
			s.err = err
			return nil, false
		}

		s.read.Add(r)
		rs = append(rs, r)
	}

	return rs, true
}

func (s *Scanner) ReadRune() (rune, error) {
	if r, ok := s.read.Get(0); ok {
		s.read.Remove(0)

		if r == '\n' {
			s.lineNum++
		}

		return r, nil
	}

	if s.err != nil {
		return 0, s.err
	}

	r, _, err := s.reader.ReadRune()
	if err != nil {
		s.err = err
		return r, err
	}

	if r == '\n' {
		s.lineNum++
	}

	return r, nil
}

func New(reader io.Reader) *Scanner {
	return &Scanner{
		bufio.NewReader(reader),
		nil,
		1, /* lineNum */
		arraylist.New[rune](),
	}
}
