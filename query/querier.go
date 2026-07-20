package query

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parse"
)

type SourceFileQuery struct {
	Imports []ir.ModuleID
	Impls   []ir.Filename
	Flags   []ir.Filename
	Decls   []ir.IrDecl
	TraitImpls []ir.IrTraitImpl
}

type ModuleQuery struct {
	Imports []ir.ModuleID
	Impls   []ir.Filename
	Flags   []ir.Filename
	Decls   []ir.IrDecl
	TraitImpls []ir.IrTraitImpl
}

type Querier struct {
	finder moduleFinder
}

func (q Querier) BaseSourceFilename(moduleID ir.ModuleID) ir.Filename {
	return q.finder.baseSourceFilename(moduleID)
}

func (q Querier) ImplSourceFilename(baseFilename ir.Filename, relativeImplFilename ir.Filename) ir.Filename {
	return q.finder.implSourceFilename(baseFilename, relativeImplFilename)
}

func findWorkspaceRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("failed to find workspace root (go.mod)")
}

func parseBplQueryOutput(output string, moduleIDName string) (ModuleQuery, error) {
	var moduleQuery ModuleQuery
	scanner := bufio.NewScanner(strings.NewReader(output))
	section := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, "{") {
			section = strings.TrimSpace(strings.TrimSuffix(line, "{"))
			continue
		}

		if line == "}" {
			section = ""
			continue
		}

		switch section {
		case "imports":
			id := ir.NewModuleID(line, ir.Pos{})
			moduleQuery.Imports = append(moduleQuery.Imports, id)
		case "impls":
			val := strings.Trim(line, "\"")
			moduleQuery.Impls = append(moduleQuery.Impls, ir.NewFilename(val, ir.Pos{}))
		case "flags":
			val := strings.Trim(line, "\"")
			moduleQuery.Flags = append(moduleQuery.Flags, ir.NewFilename(val, ir.Pos{}))
		case "decls":
			if strings.HasPrefix(line, "export ") {
				line = "pub " + strings.TrimPrefix(line, "export ")
			}
			line = strings.Replace(line, ":: ∗ -> ∗ -> ∗", "['a, 'b]", -1)
			line = strings.Replace(line, ":: ∗ -> ∗", "['a]", -1)
			line = strings.Replace(line, ":: * -> * -> *", "['a, 'b]", -1)
			line = strings.Replace(line, ":: * -> *", "['a]", -1)
			decl, err := parse.ParseSymbol[ir.IrDecl]("Decl", moduleIDName, strings.NewReader(line))
			if err != nil {
				return ModuleQuery{}, fmt.Errorf("failed to parse decl %q: %v", line, err)
			}
			moduleQuery.Decls = append(moduleQuery.Decls, decl)
		case "trait impls":
			if strings.HasPrefix(line, "export ") {
				line = "pub " + strings.TrimPrefix(line, "export ")
			}
			line = strings.Replace(line, ":: ∗ -> ∗ -> ∗", "['a, 'b]", -1)
			line = strings.Replace(line, ":: ∗ -> ∗", "['a]", -1)
			line = strings.Replace(line, ":: * -> * -> *", "['a, 'b]", -1)
			line = strings.Replace(line, ":: * -> *", "['a]", -1)
			impl, err := parse.ParseSymbol[ir.IrTraitImpl]("TraitImpl", moduleIDName, strings.NewReader(line))
			if err != nil {
				return ModuleQuery{}, fmt.Errorf("failed to parse trait impl %q: %v", line, err)
			}
			moduleQuery.TraitImpls = append(moduleQuery.TraitImpls, impl)
		}
	}

	if err := scanner.Err(); err != nil {
		return ModuleQuery{}, err
	}

	return moduleQuery, nil
}

func QuerySourceFile(inputFilename string) (SourceFileQuery, error) {
	workspaceRoot, err := findWorkspaceRoot()
	if err != nil {
		return SourceFileQuery{}, err
	}
	bplPath := filepath.Join(workspaceRoot, "bootstrap/bpl")
	if _, err := os.Stat(bplPath); err != nil {
		return SourceFileQuery{}, fmt.Errorf("bpl binary not found at %s; run 'make bootstrap' first", bplPath)
	}

	absInput, err := filepath.Abs(inputFilename)
	if err != nil {
		return SourceFileQuery{}, err
	}

	cmd := exec.Command(bplPath, "query", absInput)
	cmd.Dir = workspaceRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return SourceFileQuery{}, fmt.Errorf("failed to query file %q via %s: %s (%v)", absInput, bplPath, out, err)
	}

	mq, err := parseBplQueryOutput(string(out), inputFilename)
	if err != nil {
		return SourceFileQuery{}, err
	}
	return SourceFileQuery(mq), nil
}

func (q Querier) QueryModule(moduleID ir.ModuleID) (ModuleQuery, error) {
	workspaceRoot, err := findWorkspaceRoot()
	if err != nil {
		return ModuleQuery{}, err
	}
	bplPath := filepath.Join(workspaceRoot, "bootstrap/bpl")
	if _, err := os.Stat(bplPath); err != nil {
		return ModuleQuery{}, fmt.Errorf("bpl binary not found at %s; run 'make bootstrap' first", bplPath)
	}

	cmd := exec.Command(bplPath, "query", moduleID.Name)
	cmd.Dir = workspaceRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ModuleQuery{}, fmt.Errorf("failed to query module %q via %s: %s (%v)", moduleID, bplPath, out, err)
	}

	moduleQuery, err := parseBplQueryOutput(string(out), moduleID.Name)
	if err != nil {
		return ModuleQuery{}, err
	}

	slices.SortFunc(moduleQuery.Imports, ir.CompareModuleID)
	moduleQuery.Imports = slices.CompactFunc(moduleQuery.Imports, ir.EqualsModuleID)

	slices.SortFunc(moduleQuery.Flags, ir.CompareFilename)
	moduleQuery.Flags = slices.Compact(moduleQuery.Flags)

	return moduleQuery, nil
}


func (q Querier) QueryModuleExports(moduleID ir.ModuleID) (ModuleQuery, error) {
	moduleQuery, err := q.QueryModule(moduleID)
	if err != nil {
		return ModuleQuery{}, err
	}

	moduleQuery.Decls = slices.DeleteFunc(moduleQuery.Decls, func(decl ir.IrDecl) bool { return !decl.Export })
	return moduleQuery, nil
}

func New() (Querier, error) {
	finder, err := newModuleFinder(nil)
	if err != nil {
		return Querier{}, err
	}

	return Querier{finder}, nil
}

func NewWithWorkspace(workspace ast.Workspace) (Querier, error) {
	finder, err := newModuleFinder(&workspace)
	if err != nil {
		return Querier{}, err
	}

	return Querier{finder}, nil
}

