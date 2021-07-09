- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-06-17
- RFC PR: [beyondstorage/specs#111](https://github.com/beyondstorage/specs/issues/111)
- Tracking Issue: [beyondstorage/go-storage#602](https://github.com/beyondstorage/go-storage/issues/602)

# GSP-111: Add System Metadata in Storage Metadata

## Background

In [GSP-6], we split storage meta and object meta. In [GSP-40], we split all metadata into four groups: `global metadata`, `standard metadata`, `service metadata` and `user metadata`. To avoid concept confusion, we decide to rename `service metadata` to `system metadata` in [Idea: Find a new word to represent Service]. For storage, there's no `standard metadata` and user defined metadata is not allowed.

For now, storage related information carried in `StorageMeta`:

```go
type StorageMeta struct {
	location string
	Name     string
	WorkDir  string

	// bit used as a bitmap for object value, 0 means not set, 1 means set
	bit uint64
	m   map[string]interface{}
}
```

It only defines the `global metadata` for storage which could be used in all services, but no field to carry the system defined metadata in special services.

## Proposal

So I propose to add `system metadata` in storage metadata.

The `global metadata` has been included in object type. `system metadata` will be `systemMetadata interface{}` in `StorageMeta`. 

`system metadata` for storage will be stored in a struct that defined by service.

For services:

Example of adding system pair for storage metadata in `service.toml`:

```go
[infos.storage.meta.<system-meta>]
type = "<type>"
description = "<description>"
```

For `system metadata`, we will introduce `Strong Typed System Metadata`. We can generate following struct for specific service according to system pairs:

```go
type StorageSystemMetadata struct {
    <system meta>  <type>
    ...
}
```

And add following generated functions in service packages:

```go
// Only be used in service to set SystemMetadata into StorageMeta.
func setStorageSystemMetadata(s *StorageMeta, sm StorageSystemMetadata) {}
// GetStorageSystemMetadata will get SystemMetadata from StorageMeta.
func GetStorageSystemMetadata(s *StorageMeta) StorageSystemMetadata {}
```

## Rationale

This design is highly influenced by [GSP-40].

## Compatibility

No breaking changes.

## Implementation

- `specs`
  - Add `system-metadata` with type `any` in [info_storage_meta.toml].
- `go-storage`
  - Generate `systemMetadata` and corresponding `get/set` functions into `StorageMeta`.
  - Support generate `StorageSystemMetadata` struct to set `system metadata` for storage.


[GSP-6]: ./6-normalize-metadata.md
[GSP-40]: ./40-unify-object-metadata.md
[Idea: Find a new word to represent Service]: https://github.com/beyondstorage/specs/issues/114
[info_storage_meta.toml]: ../definitions/info_storage_meta.toml
