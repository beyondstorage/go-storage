- Author: xxchan <xxchan22f@gmail.com>
- Start Date: 2021-06-24
- RFC PR: [beyondstorage/specs#127](https://github.com/beyondstorage/specs/issues/127)
- Tracking Issue: [beyondstorage/go-storage#612](https://github.com/beyondstorage/go-storage/issues/612)

# GSP-127: Add ServiceInternel and RequestThrottled Errors

## Background

These are common errors, especially in HTTP-based services, but also apply to RPC-based services. And they are retryable errors. If we provide them, users can retry only them instead of all errors.

## Proposal

I propose to add the following global error codes:
- `ErrServiceInternal`: e.g., HTTP 5xx
- `ErrRequestThrottled`: e.g., HTTP 429 Too Many Requests/Limit Exceeded, 503 SlowDown, ...

## Rationale

Alternative names
- `ErrServer`/`ErrInternal` & `ErrThrottling`/`ErrThrottled`: short, but not conform to `Err<Noun><Predicate>`
- `ErrServerThrottling`: maybe not as common as `RequestThrottled`?
- Other throttling error code names:
  - TooManyRequests
  - ProvisionedThroughputExceeded
  - TransactionInProgress
  - RequestLimitExceeded
  - BandwidthLimitExceeded
  - LimitExceeded
  - SlowDown
  - PriorRequestNotComplete

## Compatibility

New error codes, do not break users.

## Implementation

- Add the error codes in `go-storage`.
- Add them in the [doc](https://beyondstorage.io/docs/go-storage/internal/handling-errors#list-of-global-error-codes).
- Use them in `formatError` in `go-service-*`
