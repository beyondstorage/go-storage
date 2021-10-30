package definitions

import (
	"github.com/Xuanwo/templateutils"
	"sort"
)

type Info struct {
	Namespace   string
	Category    string
	Name        string
	Type        Type
	Export      bool
	Description string

	// Only infos that declared inside definitions can set global as true.
	global bool
}

func (i Info) Global() bool {
	return i.global
}

func (i Info) NameForStructField() string {
	if i.Export {
		return templateutils.ToPascal(i.Name)
	} else {
		return templateutils.ToCamel(i.Name)
	}
}

func (i Info) NameForFunctionName() string {
	return templateutils.ToPascal(i.Name)
}

func SortInfos(is []Info) []Info {
	sort.Slice(is, func(i, j int) bool {
		x, y := is[i], is[j]

		if x.Namespace != y.Namespace {
			return x.Namespace < y.Namespace
		}
		if x.Category != y.Category {
			return x.Category < y.Category
		}
		return x.Name < y.Name
	})
	return is
}
