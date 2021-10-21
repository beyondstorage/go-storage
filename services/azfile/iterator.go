package azfile

import "github.com/Azure/azure-storage-file-go/azfile"

type objectPageStatus struct {
	maxResults int32
	prefix     string
	marker     azfile.Marker
}

func (i *objectPageStatus) ContinuationToken() string {
	if i.marker.NotDone() {
		return *i.marker.Val
	}
	return ""
}
