package ast

import (
	"fmt"
	"path"
	"strings"
)

func ModuleBaseFilename(moduleID ModuleID) string {
	moduleFilename := strings.Replace(moduleID.Name, ".", "/", -1)
	return fmt.Sprintf("%s.bpl", moduleFilename)
}

func ModuleImplFilename(baseFilename string, implID ID) string {
	return path.Join(path.Dir(baseFilename), implID.Value)
}
