package bos

type objectPageStatus struct {
	delimiter string
	marker    string
	maxKeys   int
	prefix    string
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.marker
}

type storagePageStatus struct{}

func (i *storagePageStatus) ContinuationToken() string {
	return ""
}
