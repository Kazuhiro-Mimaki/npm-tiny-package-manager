package resolver

import "npm-tiny-package-manager/types"

type DependencyStackItem struct {
	Name         types.PackageName
	Version      types.Version
	Dependencies types.Dependencies
}

type DependencyStack struct {
	Items []DependencyStackItem
}

func (s *DependencyStack) append(item DependencyStackItem) {
	s.Items = append(s.Items, item)
}

func (s *DependencyStack) pop() DependencyStackItem {
	if len(s.Items) == 0 {
		return DependencyStackItem{}
	}
	item := s.Items[len(s.Items)-1]
	s.Items = s.Items[:len(s.Items)-1]
	return item
}
