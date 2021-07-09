- Author: xxchan <xxchan22f@gmail.com>
- Start Date: 2021-06-17
- RFC PR: [beyondstorage/specs#90](https://github.com/beyondstorage/specs/issues/90)
- Tracking Issue: [beyondstorage/go-storage#600](https://github.com/beyondstorage/go-storage/issues/600)

# GSP-90: Re-support Initialization Via Connection String

- Updates:
  - [GSP-13](./13-remove-config-string.md): Deprecates it
  - [GSP-3](./3-support-service-init-via-config-string.md): Re-support this feature, which was deprecated by [GSP-13](./13-remove-config-string.md)

## Background

We had config strings like `s3://bucket_name?access_key=xxxx&secret_key=xxxx` before. This is introduced in [GSP-3] but deprecated by [GSP-13].

There are two reasons why config strings are removed, as mentioned in [GSP-13]:

1. Usability: If an application uses [go-storage]'s config string, then it either exposes the format of config string to end user directly (additional burden for end users) or it has to manipulate its config into our format (additional burden for application developers). And if the application wants to combine its other config into a config string, this is also inconvenient.
2. Implementation: We didn't find a good way to parse pairs from strings easily and safely.

So we let the users of [go-storage] construct pairs directly using `WithXxx()`. 

But passing string config is indeed more convenient, if not much. And actually parsing string into pairs is not that difficult. We can use a registry mechanism like [GSP-48]'s implementation.

## Proposal

So I propose to support service init from connection string (We rename config string to connection string since this name is more common):

We add the following APIs:
```go
func NewServicerFromString(connStr string) (types.Servicer, error) {}
func NewStoragerFromString(connStr string) (types.Storager, error) {}
```

### Format

The format of the connection string is (optional parts marked by squared brackets):

`<type>://[<name>][<work_dir>][?key1=value1&...&keyN=valueN]`

- name: storage name, e.g., bucket name. MUST NOT contain /
- work_dir: For object storage, it is prefix; for fs, it is directory path. MUST start with / for every storage services.
- For the `key=value` pairs:
  - If `=value` is missing, we just ignore the pair. But `key=` means a pair with a blank value.
  - The `key` is the pair name defined in `toml` and the format SHOULD be exactly the same.
- If their are multiple pairs with the same key, the first one will be picked.

So a valid connection string could be:

- `s3://bucket_name`
- `s3://bucket_name/prefix`
- `s3://?credential=hmac:xxxx:xxxx&endpoint=http://s3.us-east-2.amazonaws.com`
- `s3://bucket_name/prefix?credential=hmac:xxxx:xxxx&endpoint=http://s3.us-east-2.amazonaws.com`
- `fs://`
- `fs:///tmp`

#### Parseable Value Types

- `string`
- `bool`
- `int`
- `int64`
- `[]byte`: input should be base-64 encoded

For types like `credential` and `endpoint`, they are parsed in specific service's `new` functions, and treated as plain string here.


In the future, we may consider supporting complex types like `struct` (`DefaultStoragePairs`).

#### Escaping?

We don't require (and don't allow) users to escape the options. However, `&` SHOULD NOT be used since we use it as the delimiter. Other characters are OK (`/`, `=`, `+`, ...).

In the future, we may consider escaping and enabling `&`.

### Implementation

`New*FromString` will first split connection string into `(ty string, m map[string]string)`, and then parse it into `ps []Pairs`, and finally call `New*(ty, ps)`.
To support this, we have to know that a name is a pair (global or service), and its type, so we implement:

#### Pair Registry

We register the types of global pairs and service pairs.

From [go-storage] side:

```go
// RegisterSchema will register a service's pair map.
//
// Users SHOULD NOT call this function.
func RegisterSchema(ty string, m map[string]string) {}
```

From services side, we can generate following code:

```go
var pairMap = map[string]string{
	"xxxxx": "xxxxx",
}

func init() {
	// ...
	services.RegisterSchema("<type>", pairMap)
}
```

## Rationale

### Alternative Format

We can also add `<credential>@<endpoint>` in the config string, i.e.:

`type://<credential>@<endpoint>/<name><work_dir>?key1=value1&...&keyN=valueN`

However, our `credential` and `endpoint` is different from the `username:password@hostname` in URL, and we cannot easily utilize `url.Parse` to reduce our work. So we just keep the format simple currently, and we can add them in the future if really needed.

### Richer Schema?

In the pair registry, we may need more pair information beside types in the future, so we can also create a custom type `PairInfo` as the value type of the map. But to keep things simple, we'd better not overdesign now.

### Pair Conflict?

After [GSP-105], global pairs can be safely merged into the service pairs of any service without any conflict. And we can treat them in the same way.

### From Map?

We can also easily support `New*FromMap(ty string, m map[string]string)` now, but we may also want to support `New*FromMap(ty string, m map[string]interface{})`. We leave them to the future when real needs occur, instead of exposing these APIs early.

### Parsing APIs?

We can also export parsing APIs like 
```go 
func ParseMap(ty string, m map[string]string) []Pair {} 
func ParseString(connStr string) []Pair {}
```
Pros:
- Users can combine other pairs into a connection string flexibly. (Solve the usability problem mentioned in the background section.)

Cons:
- We only want to make initialization easier, instead of providing a general pair parser, which may involve additional complexity. But if we provide an API like `ParseMap` then we can't stop users from using it in other places, and thus it's effectively a general pair parser.

## Compatibility

New additional utility, no break change.

## Implementation

- Implement pair registry and parsing in [go-storage].
- Implement service code generate in [go-storage] definitions.
- Make sure all service has been updated.


[GSP-3]: ./3-support-service-init-via-config-string.md
[GSP-13]: ./13-remove-config-string.md
[GSP-48]: ./48-service-registry.md
[GSP-105]: ./105-trim-service-pair-prefix.md
[go-storage]: https://github.com/beyondstorage/go-storage
