package dropbox

import "errors"

var (
	// ErrUnexpectedEntry is the error returned when Dropbox service has returned an unexpected kind of entry.
	ErrUnexpectedEntry = errors.New("unexpected entry")
)
