package definitions

type Metadata struct {
	Name       string // Service name.
	Pairs      []Pair // System pairs.
	Infos      []Info // System infos.
	Namespaces map[string]Namespace
}

type Namespace struct {
	Required  []Pair
	Optional  []Pair
	Features  []Feature
	Functions map[string][]Pair
}
