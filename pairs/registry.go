package pairs

// PairMap is a map of the types of pairs.
type PairMap map[string]PairInfo

type PairInfo struct {
	// New returns an interface{} representing a pointer to a struct.
	// The struct has one field Value representing new zero value for the pair's value type.
	New func() interface{}
}

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
