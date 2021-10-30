package definitions

import "go.beyondstorage.io/v5/services"

type Metadata struct {
	Name       string
	Pairs      []Pair
	Infos      []Info
	Namespaces map[string]Namespace

	Factory services.Factory
}

type Namespace struct {
	Required []Pair
	Optional []Pair
	Features []Feature
	Function map[string]Function
}

type Function struct {
	Optional []Pair
}
