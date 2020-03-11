---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2020-03-11
---

# Proposal: Loose mode

## Background

Current [storage]'s pair handle behavior is inconsistent.

In all `parseStoragePairXXX` functions, we will ignore not supported pairs via only pick supported one:

```go
v, ok = values[ps.DirFunc]
if ok {
    result.HasDirFunc = true
    result.DirFunc = v.(types.ObjectFunc)
}
```

But in other pair related logic, like `storage_class` support, we also returned errors:

```go
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassStandardIA, nil
	default:
		return "", &services.PairError{
			Op:    "parse storage class",
			Err:   services.ErrStorageClassNotSupported,
			Pairs: []*types.Pair{{Key: ps.StorageClass, Value: in}},
		}
	}
}
```

So users could be confused how we handle our compatibility related issues.

## Proposal

So I propose that all a `loose` mode for all services. `loose` mode will be `off` as default, and services will return error when they reach incompatible place. And when `loose` is on, all incompatible error will be ignored.

For example:

We have a Storager who doesn't support `Size` pair in `Read`.

`loose` on: This error will be ignored.
`loose` off: Storager returns a compatibility related error.

Currently, we mixed compatibility error and other pair related error in `PairError`, we should have a new `CompatibilityError` so that user can handle them by their self.

We return `CompatibilityError` when this operation can **work** but may not behavior as expected.

## Rationale

None.

## Compatibility

- More compatibility error could be returned as `loose` mode will on as default
- Some error could be returned as `CompatibilityError` instead of `PairError`

## Implementation

Most of the work would be done by the author of this proposal.

[storage]: https://github.com/Xuanwo/storage