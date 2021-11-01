package definitions

import (
	"fmt"
	"sort"
)

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

func GetOperations(ns string) []Operation {
	switch ns {
	case NamespaceService:
		return OperationsService
	case NamespaceStorage:
		return OperationsStorage
	default:
		panic(fmt.Errorf("invalid namespace: %s", ns))
	}
}

func GetOperation(ns, name string) Operation {
	switch ns {
	case NamespaceService:
		return OperationsServiceMap[name]
	case NamespaceStorage:
		return OperationsStorageMap[name]
	default:
		panic(fmt.Errorf("invalid namespace: %s", ns))
	}
}
