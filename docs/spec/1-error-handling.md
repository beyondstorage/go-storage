---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2020-02-18
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

Expected errors are errors that implementer expected. Those errors SHOULD be defined with enough comments and any changes to them should be documented.

Unexpected errors are errors that implementer unexpected. Those errors COULD be changed or disappeared while dependence upgraded and no changelog for them.

Error should always represent as a struct which carries contextual error information. Depends on package implementation, package could have more than one error struct. 

```go
type Error struct {
	Op  string
	Err error

	ContextA string
	ContextB structB
	...
}
```

- `Op` means in which operation this error triggered.
- `Err` carries underlying error
  - For expected error, the related error SHOULD be used directly
  - For unexpected error, the error SHOULD be passed as is or warped
- `ContextX` carries contextual error information, every error context struct should implement `String() string`

Every error struct SHOULD implement following methods:

- `Error() string`
- `Unwrap() error`

String returned in `Error()` SHOULD be in the same format: 

`{Op}: {ContextA}, {ContextB}: {Err}`

`Unwrap` SHOULD always return underlying error without any operation.

## Implementer

This section will describe error handling on package implementer side.

- Expected error belongs to the package who declared, only this package CAN return this error
- Implementers CAN panic while they make sure this operation can't move on or this operation will affect or destroy data and SHOULD NOT recover

## Caller

This section will describe error handling on package caller side.

- Caller SHOULD only check package's expected error and don't check errors returned by package's imported libs
- [storage]'s package CAN panic while operations can't move on, caller SHOULD recover them by self

## Example

Expected errors

```go
var (
	// ErrUnsupportedProtocol will return if protocol is unsupported.
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
	// ErrInvalidValue means value is invalid.
	ErrInvalidValue = errors.New("invalid value")
)
```

Error struct

```go
// Error represents error related to endpoint.
type Error struct {
	Op string
	Err error

	Protocol string
	Values   []string
}

func (e *Error) Error() string {
	if e.Values == nil {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Protocol, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s, %s: %s", e.Op, e.Protocol, e.Values, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *Error) Unwrap() error {
	return e.Err
}
```

Expected error occurs

```go
err = &Error{"parse", s[0], nil, ErrUnsupportedProtocol}
```

Unexpected error occurs

```go
port, err := strconv.ParseInt(s[2], 10, 64)
if err != nil {
	return nil, &Error{"parse", ProtocolHTTP, s[1:], err}
}
```

[storage]: https://github.com/Xuanwo/storage
