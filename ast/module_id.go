package ast

import (
	"cmp"
	"fmt"
	"regexp"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

var (
	identifierRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]+$")
)

// Module identifier.
//
// If ModuleID is 'bpl:core', then `PackageID` is 'bpl', and
// `Name` is 'core'.
//
// If ModuleID is 'utils', then `PackageID` is 'main', and `Name` is
// 'utils'.
type ModuleID struct {
	// Package ID, e.g., 'bpl', 'main', etc.
	PackageID string
	// Module name, e.g., 'core', 'utils', etc.
	Name string
	// File information (if any).
	Pos ir.Pos
}

func (s ModuleID) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	if len(s.PackageID) > 0 {
		fmt.Fprint(f, s.PackageID)
		fmt.Fprint(f, ":")
	}

	fmt.Fprint(f, s.Name)
}

func NewModuleID(packageID, name string, pos ir.Pos) ModuleID {
	if packageID == "" {
		packageID = "main"
	}
	return ModuleID{packageID, name, pos}
}

func NewModuleIDFromFilename(filename string) ModuleID {
	filename = strings.TrimSuffix(filename, ".bpl")
	return NewModuleID("", filename, ir.Pos{})
}

func ValidateModuleID(moduleID ModuleID) error {
	if !identifierRegex.MatchString(moduleID.PackageID) {
		return fmt.Errorf("invalid package ID in module ID '%s'; must be an identifier", moduleID)
	}
	if !identifierRegex.MatchString(moduleID.Name) {
		return fmt.Errorf("invalid module name in module ID '%s'; must be an identifier", moduleID)
	}
	return nil
}

func CompareModuleID(id1, id2 ModuleID) int {
	if c := cmp.Compare(id1.PackageID, id2.PackageID); c != 0 {
		return c
	}
	return cmp.Compare(id1.Name, id2.Name)
}
