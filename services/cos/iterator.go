package cos

type objectPageStatus struct {
	delimiter      string
	maxKeys        int
	prefix         string
	keyMarker      string
	uploadIdMarker string
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.uploadIdMarker
}

type storagePageStatus struct {
	marker  string
	maxKeys int
}

func (i *storagePageStatus) ContinuationToken() string {
	return i.marker
}

type partPageStatus struct {
	key              string
	uploadId         string
	maxParts         string
	partNumberMarker string
}

func (i *partPageStatus) ContinuationToken() string {
	return i.partNumberMarker
}
