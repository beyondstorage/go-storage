package s3

import (
	"strconv"
)

type objectPageStatus struct {
	delimiter string
	maxKeys   int64
	prefix    string

	// Only used for object
	continuationToken string

	// Only used for part object
	keyMarker      string
	uploadIdMarker string

	expectedBucketOwner string
}

// getServiceContinuationToken equals aws.String, but return nil while empty.
//
// NOTES:
//   aws will return "InvalidArgument: The continuation token provided is incorrect" if
//   input's ContinuationToken is set to "".
func (i objectPageStatus) getServiceContinuationToken() *string {
	if i.continuationToken == "" {
		return nil
	}
	return &i.continuationToken
}

func (i *objectPageStatus) ContinuationToken() string {
	if i.uploadIdMarker != "" {
		return i.continuationToken + "/" + i.uploadIdMarker
	}
	return i.continuationToken
}

type storagePageStatus struct {
	limit    int
	offset   int
	location string
}

func (i *storagePageStatus) ContinuationToken() string {
	return strconv.FormatInt(int64(i.offset), 10)
}

type partPageStatus struct {
	key              string
	maxParts         int64
	partNumberMarker string
	uploadId         string

	expectedBucketOwner string
}

func (i *partPageStatus) ContinuationToken() string {
	return i.partNumberMarker
}
