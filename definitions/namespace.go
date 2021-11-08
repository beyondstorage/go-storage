package definitions

type Namespace interface {
	Name() string
	Operations() []Operation
	HasFeature(name string) bool
	VirtualFeatures() []Feature
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

func (s Service) VirtualFeatures() []Feature {
	fs := make([]Feature, 0)

	for _, f := range FeaturesArray {
		if f.HasNamespace(NamespaceService) && s.Features.Has(f.Name) {
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

func (s Storage) VirtualFeatures() []Feature {
	fs := make([]Feature, 0)

	for _, f := range FeaturesArray {
		if f.HasNamespace(NamespaceStorage) && s.Features.Has(f.Name) {
			fs = append(fs, f)
		}
	}
	return SortFeatures(fs)
}
