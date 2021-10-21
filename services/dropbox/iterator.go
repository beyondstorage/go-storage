package dropbox

type objectPageStatus struct {
	limit     uint32
	path      string
	recursive bool

	cursor string
}

func (i *objectPageStatus) ContinuationToken() string {
	return i.cursor
}
