package tests

import typ "github.com/aos-dev/go-storage/v3/types"

func (s *Service) formatError(op string, err error, args ...string) error {
	panic("not implemented")
}

func (s *Storage) formatError(op string, err error, args ...string) error {
	panic("not implemented")
}

func NewServicer(pairs ...typ.Pair) (typ.Servicer, error) {
	return nil, nil
}

func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	return nil, nil
}
