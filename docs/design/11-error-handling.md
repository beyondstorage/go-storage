---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2020-02-17
adds:
  - spec/1-error-handling.md
---

# Proposal: Error Handling

## Background

[storage] intends to be a production ready storage layer, and error handling is one of the most important parts.

While facing an error, user should be capable to do following things: 

- Knowing what error happened
- Deciding how to deal with
- Digging why error occurred

In order to provide those capabilities, we should return error with contextual information.

## Proposal

So I propose following changes:

- Add spec [spec: error handling] to normalize error handling across the whole lib
- All error handling related code should be refactored

Take package `segment` as example, as described in [spec: error handling], we will have an `Error` struct to carry contextual information:

```go
type SegError struct {
	Err error
    Op  string

	Seg *Segment
}

func (e *SegError) Error() string {
	return fmt.Sprintf("%v %v: %s", e.Seg, e.Op, e.Err)
}

func (e *SegError) Unwrap() error {
	return e.Err
}

func newSegmentError(err error, op string, seg *Segment) *SegError {
	return &SegError{
		Err: err,
        Op: op,
		Seg: seg,
	}
}
```

So `segment` can return error like following:

```go
func (s *Segment) ValidateParts() (err error) {
    ...

	// Zero Parts are not allowed, cause they can't be completed.
	if len(s.Parts) == 0 {
        return newSegmentError(ErrSegmentPartsEmpty, OpValidateParts, s)
	}

    ...
}
```

And caller can check those errors:

```go
err := s.ValidateParts()
```

If we don't care errors returned by `segment`:

```go
if err != nil {
    return err
}
```

If we want to handle some state:

```go
if err != nil && errors.Is(err, segment.ErrSegmentPartsEmpty) {
    log.Print("segment is empty")
}
```

If we want to get more detail error contextual information:

```go
var e SegError
if err != nil && errors.As(err, &e) {
    log.Print(e.Segment)
}
```


## Rationale

- <https://blog.golang.org/error-handling-and-go>
- <http://joeduffyblog.com/2016/02/07/the-error-model/>
- <https://blog.golang.org/go1.13-errors>

## Compatibility

Error returned by [storage] could be changed.

## Implementation

Most of the work would be done by the author of this proposal.

[storage]: https://github.com/Xuanwo/storage
[spec: error handling]: ../spec/1-error-handling.md
