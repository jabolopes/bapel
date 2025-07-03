package query

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser2"
)

const (
	bplWorkspaceFilename = "workspace.bpl"
)

func addSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}

func isPrefix(moduleID ast.ModuleID) bool {
	return strings.Contains(moduleID.Name, "*")
}

func toPrefix(moduleID ast.ModuleID) string {
	return addSlash(path.Join("/", strings.Replace(moduleID.Name, ".*", "/", -1)))
}

type moduleFinder struct {
	modulesByName   map[string]string
	modulesByPrefix map[string]string
}

func (q moduleFinder) moduleBaseFilename(moduleID ast.ModuleID) string {
	var packageName string

	if filename, ok := q.modulesByName[moduleID.Name]; ok {
		// Lookup module by exact name.
		packageName = filename
	} else {
		// Lookup module by prefix.
		moduleName := path.Join("/", strings.Replace(moduleID.Name, ".", "/", -1))

		for {
			moduleName = path.Dir(moduleName)
			if moduleName == "/" {
				break
			}

			if filename, ok := q.modulesByPrefix[addSlash(moduleName)]; ok {
				packageName = filename
				break
			}
		}
	}

	if len(packageName) > 0 {
		glog.V(1).Infof("Module %q is in package %q", moduleID, packageName)
	}

	moduleFilename := strings.Replace(moduleID.Name, ".", "/", -1)
	return fmt.Sprintf("%s.bpl", path.Join(packageName, moduleFilename))
}

func (q moduleFinder) moduleImplFilename(baseFilename string, implID ast.ID) string {
	return path.Join(path.Dir(baseFilename), implID.Value)
}

func newModuleFinder() (moduleFinder, error) {
	var workspace ast.Workspace
	switch _, err := os.Stat(bplWorkspaceFilename); {
	case err == nil:
		workspace, err = bplparser2.ParseWorkspace(bplWorkspaceFilename)
		if err != nil {
			return moduleFinder{}, err
		}
	case errors.Is(err, os.ErrNotExist):
		break
	default:
		return moduleFinder{}, err
	}

	modulesByName := map[string]string{}
	modulesByPrefix := map[string]string{}
	for _, pkg := range workspace.Packages.Packages {
		if isPrefix(pkg.ModuleID) {
			modulesByPrefix[toPrefix(pkg.ModuleID)] = pkg.Filename.Value
		} else {
			modulesByName[pkg.ModuleID.Name] = pkg.Filename.Value
		}
	}

	for name, filename := range modulesByName {
		glog.V(1).Infof("module %q in %q", name, filename)
	}

	for name, filename := range modulesByPrefix {
		glog.V(1).Infof("module prefix %q in %q", name, filename)
	}

	return moduleFinder{modulesByName, modulesByPrefix}, nil
}
