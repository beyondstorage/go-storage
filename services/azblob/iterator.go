package azblob

import "github.com/Azure/azure-storage-blob-go/azblob"

type objectPageStatus struct {
	delimiter  string
	maxResults int32
	prefix     string
	marker     azblob.Marker
}

func (i *objectPageStatus) ContinuationToken() string {
	if i.marker.NotDone() {
		return *i.marker.Val
	}
	return ""
}

type storagePageStatus struct {
	marker     azblob.Marker
	maxResults int32
}

func (i *storagePageStatus) ContinuationToken() string {
	if i.marker.NotDone() {
		return *i.marker.Val
	}
	return ""
}
