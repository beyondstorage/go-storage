package qingstor

import (
	"strconv"
)

type objectPageStatus struct {
	delimiter    string
	limit        int
	marker       string
	prefix       string
	partIdMarker string
}

func (i *objectPageStatus) ContinuationToken() string {
	if i.partIdMarker != "" {
		return i.marker + "/" + i.partIdMarker
	}
	return i.marker
}

type storagePageStatus struct {
	limit    int
	offset   int
	location string
}

func (i *storagePageStatus) ContinuationToken() string {
	return strconv.FormatInt(int64(i.offset), 10)
}

type partPageStatus struct {
	prefix           string
	limit            int
	partNumberMarker int
	uploadID         string
}

func (i *partPageStatus) ContinuationToken() string {
	return strconv.FormatInt(int64(i.partNumberMarker), 10)
}
