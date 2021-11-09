package definitions

type Namespace interface {
	Name() string
	Operations() []Operation
	HasFeature(name string) bool
	ListFeatures(ty ...string) []Feature
	ListPairs(name string) []Pair
}

func (s Service) Name() string {
	return NamespaceService
}

func (s Service) Operations() []Operation {
	return SortOperations(OperationsService)
}

func (s Service) HasFeature(name string) bool {
	return s.Features.Has(name)
}

func (s Service) ListFeatures(ty ...string) []Feature {
	fs := make([]Feature, 0)

	m := make(map[string]bool)
	for _, v := range ty {
		m[v] = true
	}

	for _, f := range FeaturesService {
		if s.Features.Has(f.Name) && m[f.Type] {
			fs = append(fs, f)
		}
	}
	return SortFeatures(fs)
}

func (s Storage) Name() string {
	return NamespaceStorage
}

func (s Storage) Operations() []Operation {
	return SortOperations(OperationsStorage)
}

func (s Storage) HasFeature(name string) bool {
	return s.Features.Has(name)
}

func (s Storage) ListFeatures(ty ...string) []Feature {
	fs := make([]Feature, 0)

	m := make(map[string]bool)
	for _, v := range ty {
		m[v] = true
	}

	for _, f := range FeaturesStorage {
		if s.Features.Has(f.Name) && m[f.Type] {
			fs = append(fs, f)
		}
	}
	return SortFeatures(fs)
}
