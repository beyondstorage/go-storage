package definitions

type Namespace interface {
	Name() string
	Operations() []Operation
	HasFeature(name string) bool
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

func (s Storage) Name() string {
	return NamespaceStorage
}

func (s Storage) Operations() []Operation {
	return SortOperations(OperationsStorage)
}

func (s Storage) HasFeature(name string) bool {
	return s.Features.Has(name)
}
