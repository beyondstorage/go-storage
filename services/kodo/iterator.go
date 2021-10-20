package kodo

type objectPageStatus struct {
	delimiter string
	limit     int
	prefix    string
	marker    string
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.marker
}

type storagePageStatus struct {
	marker string
	limit  int
}

func (i *storagePageStatus) ContinuationToken() string {
	return i.marker
}
