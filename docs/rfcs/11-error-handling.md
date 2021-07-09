- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2020-02-18
- RFC PR: N/A
- Tracking Issue: N/A

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
    Op  string
	Err error

	Seg *Segment
}

func (e *SegError) Error() string {
	return fmt.Sprintf("%s: %v: %s", e.Op, e.Seg, e.Err)
}

func (e *SegError) Unwrap() error {
	return e.Err
}
```

So `segment` can return error like following:

```go
func (s *Segment) ValidateParts() (err error) {
    ...

	// Zero Parts are not allowed, cause they can't be completed.
	if len(s.Parts) == 0 {
        return &SegError{"validate parts", ErrSegmentPartsEmpty, s}
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
