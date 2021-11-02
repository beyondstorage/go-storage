package definitions

type Metadata struct {
	Name string

	Pairs   []Pair
	Infos   []Info
	Factory []Pair

	Service Namespace
	Storage Namespace
}

func (m Metadata) Normalize() Metadata {
	m.buildDefaultPairs()
	return m
}

func (m Metadata) buildDefaultPairs() {
	dp := make([]Pair, 0)
	for _, v := range m.Pairs {
		if !v.Defaultable {
			continue
		}
		dp = append(dp, Pair{
			Name:        "default_" + v.Name,
			Type:        v.Type,
			Description: "default value for " + v.Name,
		})
	}
	m.Pairs = append(m.Pairs, dp...)
}
