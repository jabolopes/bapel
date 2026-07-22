package query

import (
	"bufio"
	"bytes"
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
	Imports    []ir.ModuleID
	Impls      []ir.Filename
	Flags      []ir.Filename
	Decls      []ir.IrDecl
	TraitImpls []ir.IrTraitImpl
}

type ModuleQuery struct {
	Imports    []ir.ModuleID
	Impls      []ir.Filename
	Flags      []ir.Filename
	Decls      []ir.IrDecl
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
			declText := line
			if strings.HasPrefix(declText, "export ") {
				declText = "pub " + strings.TrimPrefix(declText, "export ")
			}
			declText = strings.ReplaceAll(declText, "∗", "*")
			declText = strings.ReplaceAll(declText, ":: * -> * -> *", "['a, 'b]")
			declText = strings.ReplaceAll(declText, ":: * -> *", "['a]")
			for strings.Contains(declText, "= fun (") {
				idx := strings.Index(declText, "= fun (")
				closeIdx := strings.Index(declText[idx:], ")")
				if closeIdx != -1 {
					declText = declText[:idx+1] + declText[idx+closeIdx+1:]
				} else {
					break
				}
			}
			decl, err := parse.ParseSymbol[ir.IrDecl]("Decl", moduleIDName, strings.NewReader(declText))
			if err != nil {
				return ModuleQuery{}, fmt.Errorf("failed to parse decl %q: %v", line, err)
			}
			moduleQuery.Decls = append(moduleQuery.Decls, decl)
		}
	}

	if err := scanner.Err(); err != nil {
		return ModuleQuery{}, err
	}

	return moduleQuery, nil
}

func QuerySourceFile(inputFilename string) (SourceFileQuery, error) {
	if strings.HasSuffix(inputFilename, ".h") {
		content, err := os.ReadFile(inputFilename)
		if err != nil {
			return SourceFileQuery{}, err
		}
		var mq ModuleQuery
		scanner := bufio.NewScanner(bytes.NewReader(content))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			idx := strings.Index(line, "@bpl:")
			if idx == -1 {
				continue
			}
			declText := strings.TrimSpace(line[idx+len("@bpl:"):])
			if strings.HasPrefix(declText, "export ") {
				declText = "pub " + strings.TrimPrefix(declText, "export ")
			}
			declText = strings.ReplaceAll(declText, "∗", "*")
			declText = strings.ReplaceAll(declText, ":: * -> * -> *", "['a, 'b]")
			declText = strings.ReplaceAll(declText, ":: * -> *", "['a]")
			decl, err := parse.ParseSymbol[ir.IrDecl]("Decl", inputFilename, strings.NewReader(declText))
			if err != nil {
				return SourceFileQuery{}, fmt.Errorf("failed to parse decl %q in %s: %v", declText, inputFilename, err)
			}
			mq.Decls = append(mq.Decls, decl)
		}
		return SourceFileQuery(mq), nil
	}

	sourceFile, err := parse.ParseSourceFile(inputFilename)
	if err != nil {
		return SourceFileQuery{}, err
	}

	var mq ModuleQuery
	mq.Imports = append(mq.Imports, sourceFile.Imports.IDs...)
	mq.Impls = append(mq.Impls, sourceFile.Impls.Filenames...)
	mq.Flags = append(mq.Flags, sourceFile.Flags.Filenames...)

	for _, source := range sourceFile.Body {
		switch {
		case source.Is(ast.DeclSource):
			mq.Decls = append(mq.Decls, source.Decl.Decl)
		case source.Is(ast.FunctionSource):
			mq.Decls = append(mq.Decls, source.Function.Decl())
		case source.Is(ast.TraitSource):
			mq.Decls = append(mq.Decls, source.Trait.Decl())
		case source.Is(ast.ImplSource):
			irImpl, err := source.Impl.Impl.ToIr()
			if err != nil {
				return SourceFileQuery{}, err
			}
			mq.TraitImpls = append(mq.TraitImpls, irImpl)
		}
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

	baseFilename := q.finder.baseSourceFilename(moduleID)
	if strings.HasSuffix(baseFilename.Value, ".bpl") {
		if sf, err := parse.ParseSourceFile(baseFilename.Value); err == nil {
			for _, source := range sf.Body {
				if source.Is(ast.ImplSource) {
					if irImpl, err := source.Impl.Impl.ToIr(); err == nil {
						moduleQuery.TraitImpls = append(moduleQuery.TraitImpls, irImpl)
					}
				}
			}
		}
	}

	for _, relativeImplFilename := range moduleQuery.Impls {
		if !strings.HasSuffix(relativeImplFilename.Value, ".bpl") {
			continue
		}
		implFilename := q.finder.implSourceFilename(baseFilename, relativeImplFilename)
		if sf, err := parse.ParseSourceFile(implFilename.Value); err == nil {
			for _, source := range sf.Body {
				if source.Is(ast.ImplSource) {
					if irImpl, err := source.Impl.Impl.ToIr(); err == nil {
						moduleQuery.TraitImpls = append(moduleQuery.TraitImpls, irImpl)
					}
				}
			}
		}
	}

	for i := range moduleQuery.Decls {
		moduleQuery.Decls[i].Pos.Filename = moduleID.Name
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
