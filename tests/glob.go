package tests

import (
	"errors"
	"path/filepath"
)

func Glob(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, errors.New("found no tests")
	}

	return matches, nil
}
