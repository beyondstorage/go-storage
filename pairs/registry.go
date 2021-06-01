package pairs

import "reflect"

// PairMap is a map of the types of pairs.
type PairMap map[string]reflect.Type

var (
	globalPairMap   PairMap
	servicePairMaps map[string]PairMap
)

// RegisterServicePairMap will register a service's pair map.
func RegisterServicePairMap(ty string, m PairMap) {
	for pair := range m {
		if _, ok := globalPairMap[pair]; ok {
			panic("service pair name is duplicate with global pair: " + pair)
		}
	}
	servicePairMaps[ty] = m
}

func init() {
	globalPairMap = globalPairs()
	servicePairMaps = make(map[string]PairMap)
}
