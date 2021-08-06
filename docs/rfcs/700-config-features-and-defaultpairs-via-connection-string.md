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

### DefaultPairs

The format of default pair in connection string is the `key=value` pairs:

- `key` is the global or service pair name defined in `toml` and the format SHOULD be exactly the same.
- Operation has a different set of defaultable pairs. Defaultable pairs with the same name will share the default value.

A valid connection string containing default pairs could be:

- `s3://bucket_name/prefix?credential=hmac:xxxx:xxxx&server_side_encryption=xxxx&strorage_class=xxxx`

Also, users can configure the defaultable pair during initialization in the way of `WithXxx()`.

The connection string above is equivalent to:

```go
store, err := s3.NewStorage(
	ps.WithCredential("hmac:xxxx:xxxx"),
	ps.WithName("bucket_name"),
	ps.WithWorkDir("/prefix"),
    }), 
    s3.WithServiceSideEncryption("xxxx"),
    s3.WithStorageClass("xxxx"),
)
```

That means operations with the same defaultable pairs `server_side_encryption` and `storage_class` will share the default value.

### Implementation

#### Features

`*Features` are types defined in service.

**Feature Pairs Registry**

From services side:

We generate pairs like `enable_virtual_dir` according to `features` in `service.toml`. And register the feature pairs referring to [GSP-90: Re-support Initialization Via Connection String].

```go
var serviceFeaturesPairMap = map[string]string{
    // ...
    "enable_loose_pair": "bool",
}

var storageFeaturesPairMap = map[string]string{
    // ...
    "enable_loose_pair":  "bool",
    "enable_virtual_dir": "bool",
}

func init() {
	// ...
    // insert service and storage feature maps into `pairMap` for registry
	for k, v := range serviceFeaturesPairMap {
        pairMap[k] = v
    }
    for k, v := range storageFeaturesPairMap {
        pairMap[k] = v
    }
    services.RegisterSchema("<type>", pairMap)
}
```

**Parse features in New**

According to the current logic, `New*FromString` will first split connection string parse it into `ps []Pairs`, and finally call `New*(ty, ps)`. Then we can handle the default feature pairs to convert to `*Features` in `ParsePair*New()`.

#### Default Pairs

Default pairs are defined in services.

From [go-storage] side:

- Defaultable field can be added to `Pair` in `specs` to mark whether the pair could be set as default while initiate.
- Defaultable string array can be added to `Op` in `specs` to store the list of defaultable pair's name.
- Defaultable pair array can be added to `Functions` in `definitions` to store the list of defaultable pairs.

Based on the above, we can mark all the defaultable pairs for service and the defaultable pairs for a specific operation while parsing service toml.

From services side:

We can add a new field `defaultable` for operations, in which service can specify pairs that can be set as default, like:

```toml
[namespace.storage.op.write]
optional = ["content_md5", "content_type", "io_callback", "storage_class"]
defaultable = ["storage_class"]

[namespace.storage.op.create_append]
optional = ["content_type", "storage_class"]
defaultable = ["storage_class"]
```

**Default pairs registry**

The defaultable pairs also belong to global or system pairs and will be registered based on the existing logic.

**Parse pairs in `parsePair*New`**

We can add private fields correspond to defaultable pairs into the generated `Default*Pairs` to carry the shared pairs, like:

```go
type DefaultStoragePairs struct {
	// ...
	hasDefaultStorageClass  string
	defaultStorageClass     string
}
```

Then we can handle the defaultable pairs to assign the added fields of `Default*Pairs` in `ParsePair*New()`.

**Parse pairs in specific operation**

Default pairs values can be obtained from `defaultPairs` with type `Default*Pairs` stored in `Storage` and `Service`.
We can assign default value to the pair which are not passed in by args while parsing pairs for specific operation.

**Handle conflict**

When parsing pairs in specific operation：

- We should combine default pairs and `pairs from args`, and make sure that `pairs form args` can overwrite default pairs.
- When `WithDefautl*Pairs()` and `WithXxx()` are used simultaneously for initialization, the value passed in by `WithDefault*Pairs()` should be picked.
- The above conflict handling should be generated.

## Rationale

This design is based on [GSP-90: Re-support Initialization Via Connection String]. Basic connection string and parsable value types, pair registry, escaping limitation can be referred to it.

## Compatibility

No break changes.

## Implementation

- Add defaultable fields and implement feature pairs registry in [go-storage].
- Implement service code generate in [go-storage] definitions.

[GSP-90: Re-support Initialization Via Connection String]: ./90-re-support-initialization-via-connection-string.md
[go-storage]: https://github.com/beyondstorage/go-storage
