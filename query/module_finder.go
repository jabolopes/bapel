package query

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/bplparser"
)

const (
	bplWorkspaceFilename = "workspace.bpl"
)

type moduleFinder struct {
	modulesByName   map[string]string
	modulesByPrefix map[string]string
}

func (q moduleFinder) lookupModuleByName(moduleID ast.ModuleID) (string, bool) {
	filename, ok := q.modulesByName[moduleID.Name]
	return filename, ok
}

func (q moduleFinder) lookupModuleByPrefix(moduleID ast.ModuleID) (string, bool) {
	name := moduleID.Name // e.g., 'bapel.core'

	for {
		index := strings.LastIndex(name, ".")
		if index == -1 {
			name = ""
		} else {
			name = name[:index] // e.g., 'bapel'
		}

		if filename, ok := q.modulesByPrefix[name]; ok {
			return filename, true
		}

		if len(name) == 0 {
			return "", false
		}
	}
}

func (q moduleFinder) moduleBaseFilename(moduleID ast.ModuleID) ast.Filename {
	var packageName string

	if filename, ok := q.lookupModuleByName(moduleID); ok {
		packageName = filename
	} else if filename, ok := q.lookupModuleByPrefix(moduleID); ok {
		packageName = filename
	}

	if len(packageName) > 0 {
		glog.V(1).Infof("Module %q is in package %q", moduleID, packageName)
	}

	moduleFilename := strings.Replace(moduleID.Name, ".", "/", -1)
	return ast.NewFilename(fmt.Sprintf("%s.bpl", path.Join(packageName, moduleFilename)), moduleID.Pos)
}

func (q moduleFinder) moduleImplFilename(baseFilename, relativeImplFilename ast.Filename) ast.Filename {
	implFilename := path.Join(path.Dir(baseFilename.Value), relativeImplFilename.Value)
	return ast.NewFilename(implFilename, relativeImplFilename.Pos)
}

func newModuleFinder(initialWorkspace *ast.Workspace) (moduleFinder, error) {
	var workspace ast.Workspace

	if initialWorkspace != nil {
		workspace = *initialWorkspace
	} else {
		switch _, err := os.Stat(bplWorkspaceFilename); {
		case err == nil:
			workspace, err = bplparser.ParseWorkspace(bplWorkspaceFilename)
			if err != nil {
				return moduleFinder{}, err
			}

		case errors.Is(err, os.ErrNotExist):
			break
		default:
			return moduleFinder{}, err
		}
	}

	modulesByName := map[string]string{}
	modulesByPrefix := map[string]string{}
	for _, pkg := range workspace.Packages.Packages {
		switch {
		case pkg.Is(ast.ModulePackage):
			c := pkg.Module
			modulesByName[c.ModuleID.Name] = pkg.Filename.Value
		case pkg.Is(ast.PrefixPackage):
			c := pkg.Prefix
			modulesByPrefix[c.Prefix.Name] = pkg.Filename.Value
		default:
			panic(fmt.Errorf("unhandled %T %d", pkg.Case, pkg.Case))
		}
	}

	for name, filename := range modulesByName {
		glog.V(1).Infof("Module %q in %q", name, filename)
	}

	for name, filename := range modulesByPrefix {
		glog.V(1).Infof("Prefix %q in %q", name, filename)
	}

	return moduleFinder{modulesByName, modulesByPrefix}, nil
}
