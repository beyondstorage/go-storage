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
	m.buildFeaturePairs()
	return m
}

func (m *Metadata) buildDefaultPairs() {
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
	m.Factory = append(m.Factory, dp...)
}

func (m *Metadata) buildFeaturePairs() {
	dp := make(map[string]bool)
	for _, v := range []Namespace{m.Service, m.Storage} {
		for _, f := range v.ListFeatures(FeatureTypeVirtual) {
			dp["enable_"+f.Name] = true
		}
	}

	for name := range dp {
		m.Factory = append(m.Factory, PairMap[name])
	}
}
