- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-08-03
- RFC PR: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP: Config Features and DefaultPairs via Connection String

Previous discussion:

- [Connection string needs to support config StorageFeatures and DefaultPairs](https://github.com/beyondstorage/go-storage/issues/680)

## Background

In [GSP-90: Re-support Initialization Via Connection String](./90-re-support-initialization-via-connection-string.md), we introduced connection string to support service init and defined the format of connection string. 
The connection string exposes basic types and not support complex types like `struct` at present. So we can't set `ServiceFeature`/`StorageFeatures` and `DefaultServicePairs`/`DefaultStoragePairs` by connection string during initialization.

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

Obviously, It's complex for users and hard to pass in a valid string. In addition, it either exposes the internal details to end user directly (additional burden for end users) or it has to manipulate its config into our format (additional burden for application developers).

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

The format of default pair in connection string is:

`default_key=value`

- The whole key of the pair SHOULD have the prefix `default_`.
- `key` is the pair name defined in `toml` and the format SHOULD be exactly the same.
- Operation has a different set of default pairs. Pairs in operation pair list will share the same default value if have the same name .

So a valid connection string containing default pairs could be:

- `s3://bucket_name/prefix?credential=hmac:xxxx:xxxx&default_server_side_encryption=xxxx&default_strorage_class=xxxx`

Also, we generate the following function for default pairs to support configure the default pair during initialization in the way of `WithXxx()`:

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

### Implementation

#### Features

`*Features` are types defined in service.

**Feature Pairs Registry**

From services side:

We generate pairs like `enable_virtual_dir` according to `features` in `service.toml`. And register the feature pairs referring to [GSP-90: Re-support Initialization Via Connection String].

```go
var serviceFeaturesPairMap = map[string]string{
    // ...
    "default_loose_pair": "bool",
}

var storageFeaturesPairMap = map[string]string{
    // ...
    "default_loose_pair":  "bool",
    "default_virtual_dir": "bool",
}

func init() {
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

According to the current logic, `New*FromString` will first split connection string parse it into `ps []Pairs`, and finally call `New*(ty, ps)`. Then we can handle the default feature pairs to convert to `*Features` in `ParsePair*New()`.

#### Default Pairs

Default pairs are defined in services.

From [go-storage] side:

- Defaultable field can be added to `Pair` in `specs` to mark whether the pair could be set as default while initiate.
- Defaultable string array can be added to `Op` in `specs` to store the list of defaultable pair's name.
- Defaultable pair array can be added to `Functions` in `definitions` to store the list of defaultable pairs.

Based on the above, we can record all the default pairs and the default pairs for a specific operation while parsing service toml.

From services side:

We can add new field `defaultable` for operations, in which service can specify pairs that can be set as default, like:

```toml
[namespace.storage.op.write]
optional = ["content_md5", "content_type", "io_callback", "storage_class"]
defaultable = ["storage_class"]

[namespace.storage.op.create_append]
optional = ["content_type", "storage_class"]
defaultable = ["storage_class"]
```

**Default pairs registry**

We generate pairs like `default_storage_class` according to the defaultable pairs in the service. Then insert the pairs into `pairMap` like what we did in [GSP-90: Re-support Initialization Via Connection String]:

```go
var pairMap[string]string {
	// ...
	"default_storage_class" "string",
}
```

The defaultable pairs will be registered based on the existing logic.

**Parse pairs in `parsePair*New`**

We can generate the following struct to carry the shared default pairs:

```go
type DefaultConfigs struct {
	HasDefaultStorageClass  bool
	DefaultStorageClass     string
	...
}
```

Then we can handle the default pairs to assign `DefaultConfigs` in `ParsePair*New()`.

**Parse pairs in specific operation**

Default pairs values come from `DefaultConfigs` stored in `Storage` and `Service`. When parsing pairs in specific operation, we should combine default pairs and `pairs from args`, and make sure that `pairs form args` can overwrite default pairs, and this should be generated.

## Rationale

This design is based on [GSP-90: Re-support Initialization Via Connection String]. Basic connection string and parsable value types, pair registry, escaping limitation can be referred to it.

## Compatibility

Default pairs value need to be saved during initialization. New added field in `Storage` and `Service` will break the build.
We could migrate as follows:

- Release a new version for [go-storage] and all services bump to this version.
- Add a new field with type `DefaultConfig` in `Service` and `Storage` and assign it in `New*r()` for services.

## Implementation

- Add defaultable fields and implement feature pairs registry and parsing in [go-storage].
- Implement service code generate in [go-storage] definitions.

[GSP-90: Re-support Initialization Via Connection String]: ./90-re-support-initialization-via-connection-string.md
[go-storage]: https://github.com/beyondstorage/go-storage
