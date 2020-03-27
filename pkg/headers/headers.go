package headers

// We treat all headers as HTTP/2 headers
const (
	// ContentLength entity-header field indicates the size of the
	// entity-body, in decimal number of OCTETs, sent to the recipient or,
	// in the case of the HEAD method, the size of the entity-body that
	// would have been sent had the request been a GET.
	ContentLength = "content-length"
	// ContentType entity-header field indicates the media type of the
	// entity-body sent to the recipient or, in the case of the HEAD method,
	// the media type that would have been sent had the request been a GET.
	ContentType = "content-type"
	// ETag response-header field provides the current value of the
	// entity tag for the requested variant. The entity tag
	// MAY be used for comparison with other entities from the same resource
	ETag = "etag"
	// LastModified entity-header field indicates the date and time at
	// which the origin server believes the variant was last modified.
	LastModified = "last-modified"
	// Location response-header field is used to redirect the recipient
	// to a location other than the Request-URI for completion of the
	// request or identification of a new resource. For 201 (Created)
	// responses, the Location is that of the new resource which was created
	// by the request. For 3xx responses, the location SHOULD indicate the
	// server's preferred URI for automatic redirection to the resource. The
	// field value consists of a single absolute URI.
	Location = "location"
)
