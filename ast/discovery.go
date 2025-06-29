package ast

import (
	"fmt"
	"path"
)

func ModuleBaseFilename(moduleID ModuleID) string {
	// TODO: Generalize.
	var packageFilename string
	if moduleID.PackageID == "bpl" {
		packageFilename = "/home/jose/Projects/bapel"
	}

	return fmt.Sprintf("%s.bpl", path.Join(packageFilename, moduleID.Name))
}

func ModuleImplFilename(baseFilename string, implID ID) string {
	return path.Join(path.Dir(baseFilename), implID.Value)
}
