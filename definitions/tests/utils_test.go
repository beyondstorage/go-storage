package main

func (s *Service) formatError(op string, err error, args ...string) error {
	panic("not implemented")
}

func (s *Service) String() string {
	return ""
}

func (s *Storage) formatError(op string, err error, args ...string) error {
	return nil
}

func (s *Storage) String() string {
	return ""
}

func (f *Factory) newService() (*Service, error) { return nil, nil }
func (f *Factory) newStorage() (*Storage, error) { return nil, nil }
