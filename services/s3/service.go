package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

// Service is the s3 service config.
type Service struct {
	s3 s3iface.S3API
}

// String implements Service.String
func (s *Service) String() string {
	panic("implement me")
}

// Init implements Servicer.Init
func (s *Service) Init(pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) ([]storage.Storager, error) {
	panic("implement me")
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	panic("implement me")
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	panic("implement me")
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}
