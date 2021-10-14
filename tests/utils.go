package tests

import typ "go.beyondstorage.io/v5/types"

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

func NewServicer(pairs ...typ.Pair) (typ.Servicer, error) {
	return &Service{Pairs: pairs}, nil
}

func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	return &Storage{Pairs: pairs}, nil
}
