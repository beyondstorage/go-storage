package onedrive

type objectPageStatus struct {
	// basic info
	limit uint32
	rp    string
	dir   string

	// for continuationToken
	continuationToken string

	// iterator finish status
	done bool
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.continuationToken
}
