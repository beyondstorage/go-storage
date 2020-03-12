package azblob

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Service is the azblob config.
type Service struct {
	service azblob.ServiceURL

	loose bool
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer azblob")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list", err, "")
	}()

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return err
	}

	marker := azblob.Marker{}
	var output *azblob.ListContainersSegmentResponse
	for {
		output, err = s.service.ListContainersSegment(opt.Context,
			marker, azblob.ListContainersSegmentOptions{})
		if err != nil {
			return err
		}

		for _, v := range output.ContainerItems {
			store, err := s.newStorage(ps.WithName(v.Name))
			if err != nil {
				return err
			}
			opt.StoragerFunc(store)
		}

		marker = output.NextMarker
		if !marker.NotDone() {
			break
		}
	}
	return
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError("get", err, name)
	}()

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}

	return store, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError("create", err, name)
	}()

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, err
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	_, err = store.bucket.Create(opt.Context, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, name)
	}()

	opt, err := parseServicePairDelete(pairs...)
	if err != nil {
		return err
	}

	bucket := s.service.NewContainerURL(name)
	_, err = bucket.Delete(opt.Context, azblob.ContainerAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (st *Storage, err error) {
	defer func() {
		err = s.formatError("new storage", err, "")
	}()

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	bucket := s.service.NewContainerURL(opt.Name)

	c := &Storage{
		bucket: bucket,

		name:    opt.Name,
		workDir: opt.WorkDir,
		loose:   opt.Loose || s.loose,
	}
	return c, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	if s.loose && errors.Is(err, services.ErrCapabilityInsufficient) {
		return nil
	}

	return &services.ServiceError{
		Op:       op,
		Err:      formatError(err),
		Servicer: s,
		Name:     name,
	}
}
