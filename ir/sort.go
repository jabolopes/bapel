package ir

import (
	"fmt"

	"github.com/ninedraft/tsort"
	"golang.org/x/exp/slices"
)

func groupByDecl(decls []IrDecl) ([]IrDecl, []IrDecl, []IrDecl) {
	var termDecls []IrDecl
	var aliasDecls []IrDecl
	var nameDecls []IrDecl
	for _, decl := range decls {
		switch decl.Case {
		case TermDecl:
			termDecls = append(termDecls, decl)
		case AliasDecl:
			aliasDecls = append(aliasDecls, decl)
		case NameDecl:
			nameDecls = append(nameDecls, decl)
		default:
			panic(fmt.Errorf("unhandled %T %d", decl.Case, decl.Case))
		}
	}

	return termDecls, aliasDecls, nameDecls
}

func topoSortAliasDecls(decls []IrDecl) ([]IrDecl, error) {
	nodes := make([]int, 0, len(decls))
	nodesByTypeName := map[string]int{}

	edges := map[int][]int{}

	for nodeID, decl := range decls {
		{
			nodes = append(nodes, nodeID)
			nodesByTypeName[decl.Alias.ID] = nodeID
		}

		{
			freeVars := getFreeVarsFromType(decls[nodeID].Alias.Type)

			var toIDs []int
			for _, fvar := range freeVars {
				if !fvar.Is(NameType) {
					continue
				}

				nodeID, ok := nodesByTypeName[fvar.Name]
				if !ok {
					continue
				}

				toIDs = append(toIDs, nodeID)
			}

			edges[nodeID] = toIDs
		}
	}

	getEdges := func(nodeID int) []int {
		if toIDs, ok := edges[nodeID]; ok {
			return toIDs
		}
		return nil
	}

	sorted, hasCycle := tsort.Sort(nodes, getEdges)
	if hasCycle {
		return nil, fmt.Errorf("cyclic dependency between declarations")
	}

	sortedDecls := make([]IrDecl, 0, len(decls))
	for _, nodeID := range sorted {
		sortedDecls = append(sortedDecls, decls[nodeID])
	}
	slices.Reverse(sortedDecls)

	return sortedDecls, nil
}

func TopoSortDecls(decls []IrDecl) ([]IrDecl, error) {
	termDecls, aliasDecls, nameDecls := groupByDecl(decls)

	slices.SortFunc(nameDecls, CompareDecl)
	slices.SortFunc(termDecls, CompareDecl)

	aliasDecls, err := topoSortAliasDecls(aliasDecls)
	if err != nil {
		return nil, err
	}

	return append(nameDecls, append(aliasDecls, termDecls...)...), nil
}
