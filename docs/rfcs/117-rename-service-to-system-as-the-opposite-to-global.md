- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-06-21
- RFC PR: [beyondstorage/specs#117](https://github.com/beyondstorage/specs/issues/117)
- Tracking Issue: [beyondstorage/go-storage#604](https://github.com/beyondstorage/go-storage/issues/604)

# GSP-117: Rename Service to System as the Opposite to Global

- Updates:
    - [GSP-38](./38-service-pair-naming-style.md): Rename `service pair` to `system pair`
    - [GSP-40](./40-unify-object-metadata.md): Rename `service metadata` to `system metadata`
    - [GSP-105](./105-trim-service-pair-prefix.md): Rename `service pair` to `system pair`

## Background

The word `Service` has different meanings among our repositories:

- In [go-storage], `Storager` is the main interface for storage service, and `Servicer` is the container for storage. `Storage` and `Service` correspond to the abstraction layers `Storager` and `Servicer` respectively. `ServiceError` and `StorageError` follow the concept.
- In [specs], `Service Pair` or `Service Info` means the pair or info could only be used in the current service, which is opposite to `Global Pair` or `Global Info`. So we have `Service Level Object Metadata` to `serviceMetadata` in `Object`.

These seem to be confusing for developers and users. So we decide to use `system` as service or system defined information opposite to `global` to fix the confusion on `service pair`, `service metadata`, etc in [Idea: Find a new word to represent Service].

## Proposal

All our repositories should reflect this change.

- Rename `Service Pair` to `System Pair` to represent the input argument for operations in current service, which is opposite to `Global Pair`.
- Rename `Service Info` to `System Info` to represent the returning information from the current service, which is opposite to `Global Info`.
- Rename `service Metadata` to `system metadata` to represent system-defined metadata.

## Rationale

N/A

## Compatibility

All API call that used object system metadata could be affected. We could migrate as follows:

- Add `systemMetadata` in `Object` and mark `serviceMetadata`, `ObjectMetadata` related as deprecated.
- Release a new version for [go-storage] and all services bump to this version with all references to `Object.serviceMetadata`, etc updated.
- Remove deprecated fields in `Object` in the next major version.

## Implementation

- Mark `service metadata` related as deprecated.
- Add `system metadata` in specs and go-storage.
- Update all references to `Object.serviceMetadata` in services.
- Update comments about `Service Pair`, `Service Info` and so on in site.


[go-storage]: https://github.com/beyondstorage/go-storage
[specs]: https://github.com/beyondstorage/specs
[Idea: Find a new word to represent Service]: https://github.com/beyondstorage/specs/issues/114
