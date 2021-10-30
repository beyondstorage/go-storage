package definitions

import "sort"

type Operation struct {
	Name        string
	Namespace   string
	Local       bool
	Params      []Field
	Results     []Field
	Pairs       []Pair
	Description string
}

func SortOperations(ops []Operation) []Operation {
	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].Name < ops[j].Name
	})
	return ops
}
