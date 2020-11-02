// Code generated by go generate internal/cmd; DO NOT EDIT.
package pairs

import (
	"context"

	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/pkg/endpoint"
	"github.com/aos-dev/go-storage/v2/pkg/httpclient"
	"github.com/aos-dev/go-storage/v2/types"
)

// All available pairs.
const (
	// ContentMd5 will // ContentMd5
	ContentMd5 = "content_md5"
	// ContentType will // ContentType
	ContentType = "content_type"
	// Context will // Context context in all request
	Context = "context"
	// ContinuationToken will // ContinuationToken specify the continuation token for list_dir or list_prefix.
	ContinuationToken = "continuation_token"
	// Credential will // Credential specify how to provide credential for service or storage
	Credential = "credential"
	// Endpoint will // Endpoint specify how to provide endpoint for service or storage
	Endpoint = "endpoint"
	// Expire will // Expire specify when the url returned by reach will expire
	Expire = "expire"
	// HTTPClientOptions will // HTTPClientOptions sepcify the options for the http client
	HTTPClientOptions = "http_client_options"
	// Interceptor will // Interceptor
	Interceptor = "interceptor"
	// Location will // Location specify the location for service or storage
	Location = "location"
	// Name will // Name specify the storage name
	Name = "name"
	// Offset will // Offset specify offset for this request, storage will seek to this offset before read
	Offset = "offset"
	// PairPolicy will // PairPolicy
	PairPolicy = "pair_policy"
	// ReadCallbackFunc will // ReadCallbackFunc specify what todo every time we read data from source
	ReadCallbackFunc = "read_callback_func"
	// Size will // Size specify size for this request, storage will only read limited content data
	Size = "size"
	// StorageClass will // StorageClass
	StorageClass = "storage_class"
	// WorkDir will // WorkDir specify the work dir for service or storage, every operation will be relative to this dir. work_dir MUST start with / for every storage services. work_dir will be default to / if not set.
	//  For fs storage service on windows platform, the behavior is undefined.
	WorkDir = "work_dir"
)

// WithContentMd5 will apply content_md5 value to Options
// This pair is used to // ContentMd5
func WithContentMd5(v string) types.Pair {
	return types.Pair{
		Key:   ContentMd5,
		Value: v,
	}
}

// WithContentType will apply content_type value to Options
// This pair is used to // ContentType
func WithContentType(v string) types.Pair {
	return types.Pair{
		Key:   ContentType,
		Value: v,
	}
}

// WithContext will apply context value to Options
// This pair is used to // Context context in all request
func WithContext(v context.Context) types.Pair {
	return types.Pair{
		Key:   Context,
		Value: v,
	}
}

// WithContinuationToken will apply continuation_token value to Options
// This pair is used to // ContinuationToken specify the continuation token for list_dir or list_prefix.
func WithContinuationToken(v string) types.Pair {
	return types.Pair{
		Key:   ContinuationToken,
		Value: v,
	}
}

// WithCredential will apply credential value to Options
// This pair is used to // Credential specify how to provide credential for service or storage
func WithCredential(v *credential.Provider) types.Pair {
	return types.Pair{
		Key:   Credential,
		Value: v,
	}
}

// WithEndpoint will apply endpoint value to Options
// This pair is used to // Endpoint specify how to provide endpoint for service or storage
func WithEndpoint(v endpoint.Provider) types.Pair {
	return types.Pair{
		Key:   Endpoint,
		Value: v,
	}
}

// WithExpire will apply expire value to Options
// This pair is used to // Expire specify when the url returned by reach will expire
func WithExpire(v int) types.Pair {
	return types.Pair{
		Key:   Expire,
		Value: v,
	}
}

// WithHTTPClientOptions will apply http_client_options value to Options
// This pair is used to // HTTPClientOptions sepcify the options for the http client
func WithHTTPClientOptions(v *httpclient.Options) types.Pair {
	return types.Pair{
		Key:   HTTPClientOptions,
		Value: v,
	}
}

// WithInterceptor will apply interceptor value to Options
// This pair is used to // Interceptor
func WithInterceptor(v types.Interceptor) types.Pair {
	return types.Pair{
		Key:   Interceptor,
		Value: v,
	}
}

// WithLocation will apply location value to Options
// This pair is used to // Location specify the location for service or storage
func WithLocation(v string) types.Pair {
	return types.Pair{
		Key:   Location,
		Value: v,
	}
}

// WithName will apply name value to Options
// This pair is used to // Name specify the storage name
func WithName(v string) types.Pair {
	return types.Pair{
		Key:   Name,
		Value: v,
	}
}

// WithOffset will apply offset value to Options
// This pair is used to // Offset specify offset for this request, storage will seek to this offset before read
func WithOffset(v int64) types.Pair {
	return types.Pair{
		Key:   Offset,
		Value: v,
	}
}

// WithPairPolicy will apply pair_policy value to Options
// This pair is used to // PairPolicy
func WithPairPolicy(v types.PairPolicy) types.Pair {
	return types.Pair{
		Key:   PairPolicy,
		Value: v,
	}
}

// WithReadCallbackFunc will apply read_callback_func value to Options
// This pair is used to // ReadCallbackFunc specify what todo every time we read data from source
func WithReadCallbackFunc(v func([]byte)) types.Pair {
	return types.Pair{
		Key:   ReadCallbackFunc,
		Value: v,
	}
}

// WithSize will apply size value to Options
// This pair is used to // Size specify size for this request, storage will only read limited content data
func WithSize(v int64) types.Pair {
	return types.Pair{
		Key:   Size,
		Value: v,
	}
}

// WithStorageClass will apply storage_class value to Options
// This pair is used to // StorageClass
func WithStorageClass(v string) types.Pair {
	return types.Pair{
		Key:   StorageClass,
		Value: v,
	}
}

// WithWorkDir will apply work_dir value to Options
// This pair is used to // WorkDir specify the work dir for service or storage, every operation will be relative to this dir. work_dir MUST start with / for every storage services. work_dir will be default to / if not set.
//  For fs storage service on windows platform, the behavior is undefined.
func WithWorkDir(v string) types.Pair {
	return types.Pair{
		Key:   WorkDir,
		Value: v,
	}
}

// Parse will parse a k-v map to pairs slice.
func Parse(m map[string]interface{}) ([]types.Pair, error) {
	pairs := make([]types.Pair, 0, len(m))

	var err error

	for k, v := range m {
		var pv interface{}
		switch k {
		case ContentMd5:
			switch rv := v.(type) {
			case string:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   ContentMd5,
					Type:  "string",
					Value: v,
				}
			}
		case ContentType:
			switch rv := v.(type) {
			case string:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   ContentType,
					Type:  "string",
					Value: v,
				}
			}
		case Context:
			switch rv := v.(type) {
			case context.Context:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Context,
					Type:  "context.Context",
					Value: v,
				}
			}
		case ContinuationToken:
			switch rv := v.(type) {
			case string:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   ContinuationToken,
					Type:  "string",
					Value: v,
				}
			}
		case Credential:
			switch rv := v.(type) {
			case *credential.Provider:
				pv = rv
			case string:
				pv, err = credential.Parse(rv)
				if err != nil {
					return nil, &Error{
						Op:    "parse",
						Err:   err,
						Key:   Credential,
						Type:  "*credential.Provider",
						Value: v,
					}
				}
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Credential,
					Type:  "*credential.Provider",
					Value: v,
				}
			}
		case Endpoint:
			switch rv := v.(type) {
			case endpoint.Provider:
				pv = rv
			case string:
				pv, err = endpoint.Parse(rv)
				if err != nil {
					return nil, &Error{
						Op:    "parse",
						Err:   err,
						Key:   Endpoint,
						Type:  "endpoint.Provider",
						Value: v,
					}
				}
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Endpoint,
					Type:  "endpoint.Provider",
					Value: v,
				}
			}
		case Expire:
			switch rv := v.(type) {
			case int:
				pv = rv
			case string:
				pv, err = parseInt(rv)
				if err != nil {
					return nil, &Error{
						Op:    "parse",
						Err:   err,
						Key:   Expire,
						Type:  "int",
						Value: v,
					}
				}
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Expire,
					Type:  "int",
					Value: v,
				}
			}
		case HTTPClientOptions:
			switch rv := v.(type) {
			case *httpclient.Options:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   HTTPClientOptions,
					Type:  "*httpclient.Options",
					Value: v,
				}
			}
		case Interceptor:
			switch rv := v.(type) {
			case types.Interceptor:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Interceptor,
					Type:  "types.Interceptor",
					Value: v,
				}
			}
		case Location:
			switch rv := v.(type) {
			case string:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Location,
					Type:  "string",
					Value: v,
				}
			}
		case Name:
			switch rv := v.(type) {
			case string:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Name,
					Type:  "string",
					Value: v,
				}
			}
		case Offset:
			switch rv := v.(type) {
			case int64:
				pv = rv
			case string:
				pv, err = parseInt64(rv)
				if err != nil {
					return nil, &Error{
						Op:    "parse",
						Err:   err,
						Key:   Offset,
						Type:  "int64",
						Value: v,
					}
				}
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Offset,
					Type:  "int64",
					Value: v,
				}
			}
		case PairPolicy:
			switch rv := v.(type) {
			case types.PairPolicy:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   PairPolicy,
					Type:  "types.PairPolicy",
					Value: v,
				}
			}
		case ReadCallbackFunc:
			switch rv := v.(type) {
			case func([]byte):
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   ReadCallbackFunc,
					Type:  "func([]byte)",
					Value: v,
				}
			}
		case Size:
			switch rv := v.(type) {
			case int64:
				pv = rv
			case string:
				pv, err = parseInt64(rv)
				if err != nil {
					return nil, &Error{
						Op:    "parse",
						Err:   err,
						Key:   Size,
						Type:  "int64",
						Value: v,
					}
				}
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   Size,
					Type:  "int64",
					Value: v,
				}
			}
		case StorageClass:
			switch rv := v.(type) {
			case string:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   StorageClass,
					Type:  "string",
					Value: v,
				}
			}
		case WorkDir:
			switch rv := v.(type) {
			case string:
				pv = rv
			default:
				return nil, &Error{
					Op:    "parse",
					Err:   ErrPairTypeMismatch,
					Key:   WorkDir,
					Type:  "string",
					Value: v,
				}
			}
		default:
			continue
		}
		pairs = append(pairs, types.Pair{Key: k, Value: pv})
	}

	return pairs, nil
}
