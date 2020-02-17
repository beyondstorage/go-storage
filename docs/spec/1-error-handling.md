---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2020-02-17
added_by:
  - design/11-error-handling.md
---

# Spec: Error Handling

This spec will describe how to handle errors in [storage].

## Definitions

- `error`: program is not running properly
- `package`: All valid go package in [storage]
- `implementer`: People who implement `package`
- `caller`: People who use/call `package`

## Error

From [storage]'s side, errors can be classified as following:

- Expected errors
- Unexpected errors

Expected errors are errors that implementer expected. For example, package `segment`'s implementer expect following errors:

```go
ErrPartSizeInvalid     = errors.New("part size invalid")
ErrPartIntersected     = errors.New("part intersected")
ErrSegmentNotInitiated = errors.New("segment not initiated")
ErrSegmentPartsEmpty   = errors.New("segment Parts are empty")
ErrSegmentNotFulfilled = errors.New("segment not fulfilled")
```

Those errors should be defined with enough comments and any changes to them should be documented.

Unexpected errors are errors that implementer unexpected. There are many reasons for them:

- Implementation BUG
- Invalid input
- Third-party libraries
- ...

Those errors could be changed or disappeared while dependence upgraded and no changelog for them.

## Implementer

This section will describe error handling on package implementer side.

- Expected error SHOULD always be formatted into a struct which carries contextual error information
- Expected error belongs to the package who declared, only this package CAN return this error
- Implementers CAN panic while they make sure this operation can't move on or this operation will affect or destroy data and SHOULD NOT recover
- Unexpected error SHOULD return as is or with [error warp]

## Caller

This section will describe error handling on package caller side.

- Caller SHOULD only check package's expected error and don't check errors returned by package's imported libs
- [storage]'s package CAN panic while operations can't move on, caller SHOULD recover them by self

[storage]: https://github.com/Xuanwo/storage
[error warp]: https://blog.golang.org/go1.13-errors
