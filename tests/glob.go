package tests

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	"github.com/google/go-cmp/cmp"
)

var regen bool
var glob string

func init() {
	flag.BoolVar(&regen, "regen", false, "Whether to regenerate test output files.")
	flag.StringVar(&glob, "glob", "", "Override test pattern glob.")
}

func Glob(pattern string) ([]string, error) {
	if len(glob) > 0 {
		pattern = glob
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, errors.New("found no tests")
	}

	return matches, nil
}

func DiffOutRegen(got, wantFile string) (string, error) {
	if regen {
		if err := os.WriteFile(wantFile, []byte(got), 0644); err != nil {
			return "", err
		}
	}

	want, err := os.ReadFile(wantFile)
	if err != nil {
		return "", err
	}

	if diff := cmp.Diff(string(want), got); len(diff) > 0 {
		return diff, nil
	}

	return "", nil
}
