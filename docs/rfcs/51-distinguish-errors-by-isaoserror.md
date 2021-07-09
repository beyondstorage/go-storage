- Author: xxchan <xxchan22f@gmail.com>
- Start Date: 2021-05-10
- RFC PR: [beyondstorage/specs#51](https://github.com/beyondstorage/specs/issues/51)
- Tracking Issue: N/A

# AOS-51: Distinguish Errors by IsAosError

## Background

Currently we use a function named `formatError` in generated code to turn SDK errors into our errors, as defined in [AOS-47]. However, sometimes we have used our errors before `formatError`, so our errors will also be formatted. Then it's hard for us to distinguish between non-SDK errors and our errors.

Take `go-storage-s3` as an example:

```go
func formatError(err error) error {
	e, ok := err.(awserr.RequestFailure)
	if !ok {
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}
	switch e.Code() {
	// AWS SDK will use status code to generate awserr.Error, so "NotFound" should also be supported.
	case "NoSuchKey", "NotFound":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case "AccessDenied":
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}
}
```

This is problematic since our `PairRequiredError` will also be turned into `ErrUnexpected`.

## Proposal

So I propose to add support for recognizing errors defined in `go-storage`, including:

- Add a new unexported type `errorCode` and a function `NewErrorCode(string) error`, to replace sentinel errors created by `errors.New()`.
- Add a new interface `AosError`, which has a method `IsAosError()`.
- Top-level errors (`InitError`, `ServiceError`, `StorageError`) do not need to implement `IsAosError()`, since `AosError` is used when formatting the `Err` inside top-level errors.
- `errorCode` and every error `struct` SHOULD implement `IsAosError()`, which has an empty function body.

## Rationale

We cannot add methods to sentinel errors, so we need to replace them.

### Alternative

- Add a special `errorCode` `ErrAos`
- `ErrorCode` and every error `struct` should implement `Is(target error) bool { return target == ErrAos}`

Benchmark shows this approach is slower than the interface way, since `errors.Is` is more expensive than type assertion.

Besides, `errors.Is` is general-purpose while using a special-purpose interface can be a litter clearer.

## Compatibility

This proposal will not break users, since `errors.New()` create a private error type which cannot be used by users.

## Implementation

- `go-storage`:
  - Add `errorCode` `NewErrorCode(string) error`
  - Replace `errors.New` with `NewErrorCode`
  - Add interface `AosError`, and implement it for `errorCode` and error `struct`s
- `go-service-*`:
  - Replace `errors.New` with `NewErrorCode`
  - We don't have service error `struct` now, so don't need to implement `AosError`.

[AOS-47]: ./47-additional-error-specification.md
