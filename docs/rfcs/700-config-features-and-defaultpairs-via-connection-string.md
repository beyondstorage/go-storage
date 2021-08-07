- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-08-03
- RFC PR: [beyondstorage/go-storage#700](https://github.com/beyondstorage/go-storage/issues/700)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-700: Config Features and DefaultPairs via Connection String

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

The format of features pairs in connection string is:

`enable_key=value`

- The whole key of the pair SHOULD have the prefix `enable_`.
- `key` is the pair name defined in `features.toml` and supported in service, and the format SHOULD be exactly the same.
- `value` SHOULD be `true` or `false`.

So a valid connection string containing features could be:

- `s3://bucket_name/prefix?credential=hmac:xxxx:xxxx&enable_loose_pair=false&enable_virtual_dir=true`

The connection string above is equivalent to:

```go
store, err := s3.NewStorage(
	ps.WithCredential("hmac:xxxx:xxxx"),
	ps.WithName("bucket_name"),
	ps.WithWorkDir("/prefix"),
	s3.WithStorageFeatures(s3.StorageFeaturs{LoosePair:  true,
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

Also, we generate the following function to support configure the default pair during initialization in the way of `WithXxx()`:

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

**Add pairs for features and default paris**

Default pairs and features can be added into `Optioanl` pairs of `New` for service, like:

```toml
[namespace.storage.new]
required = ["name"]
optional = ["storage_features", "default_storage_pairs", "default_storage_class", "enable_virtual_dir"]
```

When parsing pairs from toml, we should:

- Check whether the default pair belongs to global or system pair and feature belongs to the supported ones.
- Add pairs for the optional pairs of `New` with `default_` and `enable_` prefix.
- Mark them as defaultable, the same with global and system paris corresponding to the default ones. 

**Parse features in `parsePair*New`**

We can handle the optional defaultable pairs with prefix `enable_`, assign value to the corresponding field of `*Features` in `ParsePair*New()`.

**Parse default pairs in `parsePair*New`**

We can add private fields correspond to default pairs into the generated `Default*Pairs` to carry the shared values, like:

```go
type DefaultStoragePairs struct {
	// ...
	hasDefaultStorageClass  string
	defaultStorageClass     string
}
```

Then we can handle the optional defaultable pairs with prefix `default_` to assign the added fields of `Default*Pairs` in `ParsePair*New()`.

**Parse pairs in specific operation**

Based on the updated `Default*Pairs`, when parsing pairs in specific operation, we can assign default value to defaultable pair if the value not passed in from args.

**Handle conflict**

When parsing pairs in specific operation：

- We should combine default pairs and `pairs from args`, and make sure that `pairs form args` can overwrite default pairs.
- When using `WithDefautl*Pairs()` and `WithDefaultXxx()` at the same time for initialization, we will pick `Default*Pairs` first, and append pairs that parsed from `WithDefaultXxx()`.
- The above conflict handling should be generated.

## Rationale

This design is based on [GSP-90: Re-support Initialization Via Connection String]. Basic connection string and parsable value types, pair registry, escaping limitation can be referred to it.

## Compatibility

No break changes.

## Implementation

- Add pairs for default pairs and features in [go-storage].
- Implement service code generate in [go-storage] definitions.

[GSP-90: Re-support Initialization Via Connection String]: ./90-re-support-initialization-via-connection-string.md
[go-storage]: https://github.com/beyondstorage/go-storage
