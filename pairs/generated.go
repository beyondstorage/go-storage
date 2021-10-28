// Code generated by go generate cmd/definitions; DO NOT EDIT.
package pairs

import (
	"context"

	"go.beyondstorage.io/v5/pkg/httpclient"
	"go.beyondstorage.io/v5/types"
)

// WithContentDisposition will apply content_disposition value to Options.
//
// ContentDisposition
func WithContentDisposition(v string) (p types.Pair) {
	return types.Pair{Key: "content_disposition", Value: v}
}

// WithContentMd5 will apply content_md5 value to Options.
//
// ContentMd5
func WithContentMd5(v string) (p types.Pair) {
	return types.Pair{Key: "content_md5", Value: v}
}

// WithContentType will apply content_type value to Options.
//
// ContentType
func WithContentType(v string) (p types.Pair) {
	return types.Pair{Key: "content_type", Value: v}
}

// WithContext will apply context value to Options.
//
// Context
func WithContext(v context.Context) (p types.Pair) {
	return types.Pair{Key: "context", Value: v}
}

// WithContinuationToken will apply continuation_token value to Options.
//
// ContinuationToken specify the continuation token for list
func WithContinuationToken(v string) (p types.Pair) {
	return types.Pair{Key: "continuation_token", Value: v}
}

// WithCredential will apply credential value to Options.
//
// Credential specify how to provide credential for service or storage
func WithCredential(v string) (p types.Pair) {
	return types.Pair{Key: "credential", Value: v}
}

// WithDefaultContentType will apply default_content_type value to Options.
//
// DefaultContentType
func WithDefaultContentType(v string) (p types.Pair) {
	return types.Pair{Key: "default_content_type", Value: v}
}

// WithDefaultIoCallback will apply default_io_callback value to Options.
//
// DefaultIoCallback specify what todo every time we read data from source
func WithDefaultIoCallback(v func([]byte)) (p types.Pair) {
	return types.Pair{Key: "default_io_callback", Value: v}
}

// WithEndpoint will apply endpoint value to Options.
//
// Endpoint specify how to provide endpoint for service or storage
func WithEndpoint(v string) (p types.Pair) {
	return types.Pair{Key: "endpoint", Value: v}
}

// WithHTTPClientOptions will apply http_client_options value to Options.
//
// HTTPClientOptions
func WithHTTPClientOptions(v *httpclient.Options) (p types.Pair) {
	return types.Pair{Key: "http_client_options", Value: v}
}

// WithInterceptor will apply interceptor value to Options.
//
// Interceptor
func WithInterceptor(v types.Interceptor) (p types.Pair) {
	return types.Pair{Key: "interceptor", Value: v}
}

// WithIoCallback will apply io_callback value to Options.
//
// IoCallback specify what todo every time we read data from source
func WithIoCallback(v func([]byte)) (p types.Pair) {
	return types.Pair{Key: "io_callback", Value: v}
}

// WithListMode will apply list_mode value to Options.
//
// ListMode
func WithListMode(v types.ListMode) (p types.Pair) {
	return types.Pair{Key: "list_mode", Value: v}
}

// WithLocation will apply location value to Options.
//
// Location specify the location for service or storage
func WithLocation(v string) (p types.Pair) {
	return types.Pair{Key: "location", Value: v}
}

// WithMultipartID will apply multipart_id value to Options.
//
// MultipartID
func WithMultipartID(v string) (p types.Pair) {
	return types.Pair{Key: "multipart_id", Value: v}
}

// WithName will apply name value to Options.
//
// Name specify the storage name
func WithName(v string) (p types.Pair) {
	return types.Pair{Key: "name", Value: v}
}

// WithObjectMode will apply object_mode value to Options.
//
// ObjectMode ObjectMode hint
func WithObjectMode(v types.ObjectMode) (p types.Pair) {
	return types.Pair{Key: "object_mode", Value: v}
}

// WithOffset will apply offset value to Options.
//
// Offset specify offset for this request, storage will seek to this offset before read
func WithOffset(v int64) (p types.Pair) {
	return types.Pair{Key: "offset", Value: v}
}

// WithSize will apply size value to Options.
//
// Size specify size for this request, storage will only read limited content data
func WithSize(v int64) (p types.Pair) {
	return types.Pair{Key: "size", Value: v}
}

// WithWorkDir will apply work_dir value to Options.
//
// WorkDir specify the work dir for service or storage, every operation will be relative to this dir.
// work_dir SHOULD be an absolute path.
// work_dir will be default to / if not set.
// work_dir SHOULD be Unix style for object storage services.
// For fs storage service on windows platform, the behavior is defined separately.
func WithWorkDir(v string) (p types.Pair) {
	return types.Pair{Key: "work_dir", Value: v}
}
