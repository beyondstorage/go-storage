package gcs

type objectPageStatus struct {
	delimiter string
	prefix    string
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.prefix
}
