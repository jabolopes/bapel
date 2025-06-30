package ast

import (
	"fmt"
	"path"
	"strings"
)

func ModuleBaseFilename(moduleID ModuleID) string {
	// TODO: Generalize.
	var packageFilename string
	if moduleID.PackageID == "bpl" {
		packageFilename = "/home/jose/Projects/bapel"
	}

	moduleFilename := strings.Replace(moduleID.Name, ".", "/", -1)

	return fmt.Sprintf("%s.bpl", path.Join(packageFilename, moduleFilename))
}

func ModuleImplFilename(baseFilename string, implID ID) string {
	return path.Join(path.Dir(baseFilename), implID.Value)
}
