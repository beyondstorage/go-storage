- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-08-03
- RFC PR: [beyondstorage/go-storage#700](https://github.com/beyondstorage/go-storage/issues/700)
- Tracking Issue: [beyondstorage/go-storage#704](https://github.com/beyondstorage/go-storage/issues/704)

# GSP-700: Config Features and DefaultPairs via Connection String

- Updated By:
  - [GSP-725](./725-add-defaultable-property-for-pair.md): Deprecate `defaultable` in namespace and split default pairs into global and system.

Previous discussion:

- [Connection string needs to support config StorageFeatures and DefaultPairs](https://github.com/beyondstorage/go-storage/issues/680)

## Background

In [GSP-90: Re-support Initialization Via Connection String](./90-re-support-initialization-via-connection-string.md), we introduced connection string to support service init and defined the format of connection string. 
The connection string exposes basic types and not support complex types like `struct` at present. So we can't set `*Feature` and `Default*Pairs` by connection string during initialization.

`Default*Pairs` is a slightly complex struct:

```go
// DefaultStoragePairs is default pairs for specific action
type DefaultStoragePairs struct {
	Create            []Pair
	Delete            []Pair
	List              []Pair
	Metadata          []Pair
	Read              []Pair
	Stat              []Pair
	Write             []Pair
}
```

A possible format that follows the struct in connection string for default storage paris may be:

`default_storage_pairs={write:[storage_class=STANDARD_IA,key2=value],read:[]}`

Obviously, it's complex for users and hard to pass in a valid string. In addition, it either exposes the internal details to end user directly (additional burden for end users) or it has to manipulate its config into our format (additional burden for application developers).

## Proposal

I propose to support config features and default pairs in plain string format in connection string.

### Features

The format for features setting in connection string is:

`enable_key`

- The whole key SHOULD have the prefix `enable_`.
- `key` is the feature name defined in `features.toml` and supported in service, and the format SHOULD be exactly the same.
- No value for the features setting string, or we can assume that the value is always true.

So a valid connection string containing features could be:

- `s3://bucket_name/prefix?credential=hmac:xxxx:xxxx&enable_loose_pair&enable_virtual_dir`

The connection string above is equivalent to:

```go
store, err := s3.NewStorage(
	ps.WithCredential("hmac:xxxx:xxxx"),
	ps.WithName("bucket_name"),
	ps.WithWorkDir("/prefix"),
	s3.WithStorageFeatures(s3.StorageFeaturs{
		LoosePair:  true,
		VirtualDir: true,
    })

)
```

### Default Pairs

The format of default pair in connection string is:

`default_key=value`

- The whole key of the pair SHOULD have the prefix `default_`.
- `key` is the global or service pair name defined in `toml`, and the format SHOULD be exactly the same.
- For `value`, parsable value types are `string`, `bool`, `int`, `int64`, `uint64`, `[]byte` and `time.Duration` at present.
- Pairs of operation with the same name will share the default value.

A valid connection string containing default pairs could be:

- `s3://bucket_name/prefix?credential=hmac:xxxx:xxxx&default_server_side_encryption=xxxx&default_strorage_class=xxxx`

To ensure the consistency of the semantics and behavior of the current two initialization methods for users, we generate the following function to support configure the default pair during initialization in the way of `WithXxx()` to achieve the same result as the way of connection string:

Take `default_storage_class` for example:

```go
func WithDefaultStorageClass(v string) Pair {
	reutrn Pair{
		Key:   "default_storage_class",
		Value: v,
	}
}
```

The connection string above is equivalent to:

```go
store, err := s3.NewStorage(
	ps.WithCredential("hmac:xxxx:xxxx"),
	ps.WithName("bucket_name"),
	ps.WithWorkDir("/prefix"),
    }), 
    s3.WithDefaultServiceSideEncryption("xxxx"),
    s3.WithDefaultStorageClass("xxxx"),
)
```

That means pair of operations with the same name `server_side_encryption` or `storage_class` will share the default value.

### Implementation

Default pairs and features are defined in service.

**Add `defaultable` in namespace**

Add a new field `defaultable` in namespace to specify the defaultable pair key list, like:

```toml
[namespace.storage]
features = ["virtual_dir"]
implement = ["direr", "multiparter"]
defaultable = ["excepted_bucket_owner"]
```

Pair keys listed in `defaultable` will be treated as allowing users to set default values. When parsing toml files, we should:

- Check whether the key in `defaultable` belongs to the pair list of operations in the namespace.
- Generate default pairs according to `defaultable` and `features` fields for pair registry and parsing.

**Parse features in `parsePair*New`**

We can get the feature pairs with prefix `enable_` from connection string or `WithEnableXxx()`, then update the corresponding field of `*Features` in `ParsePair*New()`.

**Parse default pairs in `parsePair*New`**

We can convert the default pairs with prefix `default_` passed in by connection string or `WithDefaultXxx()` to original pairs and append it to the pair array of supported operations in `Default*Pairs` in `ParsePair*New()`.

**Usage notes**

- Using `WithDefautl*Pairs()` and `WithDefaultXxx()` at the same time for initialization are not allowed at present.
- `WithEnableXxx()` will fill or overwrite the value of the corresponding field in `*Features` when it is used in conjunction with `With*Features()`.

## Rationale

This design is based on [GSP-90: Re-support Initialization Via Connection String]. Basic connection string and parsable value types, pair registry, escaping limitation can be referred to it.

### Alternative Implementation

We can add a new field `defaultable` in operations, like:

```toml
[namespace.storage.op.read]
optional = ["offset", "io_callback", "size", "excepted_bucket_owner"]
defaultable = ["excepted_bucket_owner"]
```

On this basis, we can parse the toml files and get a full list of defaultable pairs, and then, generate related pairs.

Compared with `add defaultable in namespace`, we have to list the same defaultable pair for different operations and traverse the pair lists of operations to get full defaultable pairs.

## Compatibility

No break changes.

## Implementation

- Add the new field into `service.toml`.
- Add pairs for default pairs and features in [go-storage].
- Implement service code generate in [go-storage] definitions.

[GSP-90: Re-support Initialization Via Connection String]: ./90-re-support-initialization-via-connection-string.md
[go-storage]: https://github.com/beyondstorage/go-storage
