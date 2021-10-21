package uss

type objectPageStatus struct {
	limit  string // limit passed in header as a string
	prefix string
	iter   string
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.iter
}
