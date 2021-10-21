package minio

import (
	"strconv"

	"github.com/minio/minio-go/v7"
)

type storagePageStatus struct {
	buckets []minio.BucketInfo
}

func (i *storagePageStatus) ContinuationToken() string {
	return ""
}

type objectPageStatus struct {
	bufferSize int
	counter    int
	options    minio.ListObjectsOptions

	objChan <-chan minio.ObjectInfo
}

func (i *objectPageStatus) ContinuationToken() string {
	return strconv.Itoa(i.counter)
}
