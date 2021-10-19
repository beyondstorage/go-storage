package gdrive

type objectPageStatus struct {
	limit     uint32
	path      string
	pageToken string
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.pageToken
}
