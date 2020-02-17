---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2020-02-16
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

- Expected error should always be formatted into a struct which carries contextual error information

## Caller

This section will describe error handling on package caller side.

[storage]: https://github.com/Xuanwo/storage
