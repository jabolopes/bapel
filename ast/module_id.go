package ast

import (
	"cmp"
	"fmt"
	"regexp"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

var (
	identifierRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]+$")
)

type ModuleID struct {
	// Module name, e.g., 'main', 'bapel.core', etc.
	Name string
	// File information (if any).
	Pos ir.Pos
}

func (s ModuleID) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	isQuoted := verb == 'q'
	if isQuoted {
		fmt.Fprint(f, `"`)
	}
	fmt.Fprint(f, s.Name)
	if isQuoted {
		fmt.Fprint(f, `"`)
	}
}

func NewModuleID(name string, pos ir.Pos) ModuleID {
	return ModuleID{name, pos}
}

func NewModuleIDFromFilename(filename string) ModuleID {
	filename = strings.TrimSuffix(filename, ".bpl")
	filename = strings.Replace(filename, "/", ".", -1)
	return NewModuleID(filename, ir.Pos{})
}

func ValidateModuleID(moduleID ModuleID) error {
	splits := strings.Split(moduleID.Name, ".")
	if len(splits) <= 0 {
		return fmt.Errorf("invalid module name in module ID %q. Valid module names are, e.g., 'main', 'bapel.core', etc", moduleID)
	}
	for _, split := range splits {
		if !identifierRegex.MatchString(split) {
			return fmt.Errorf("invalid module name in module ID %q. Valid module names are, e.g., 'main', 'bapel.core', etc", moduleID)
		}
	}
	return nil
}

func CompareModuleID(id1, id2 ModuleID) int {
	return cmp.Compare(id1.Name, id2.Name)
}
