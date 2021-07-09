- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-01-12
- RFC PR: N/A
- Tracking Issue: N/A

# Proposal: Segment API Redesign

## Background

Like List Operations, Segment related APIs also changed a lot.

The first design looks like:

```go
type Segmenter interface {
    ListSegments(path string, pairs ...*types.Pair) (err error)
    InitSegment(path string, pairs ...*types.Pair) (id string, err error)
    WriteSegment(id string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error)
    CompleteSegment(id string, pairs ...*types.Pair) (err error)
    AbortSegment(id string, pairs ...*types.Pair) (err error)
}
```

But this design is buggy: a single segment/multipart can't be identified by ID. So [Proposal: Return segment interface instead](18-return-segment-interface-instead.md) introduced a `Segment` interface:

```go
type Segment interface {
    String() string

    ID() string
    Path() string
}
```

Changing the `Segmenter` interface into:

```go
type Segmenter interface {
    ListSegments(path string, pairs ...*types.Pair) (err error)
    InitSegment(path string, pairs ...*types.Pair) (seg segment.Segment, err error)
    WriteSegment(seg segment.Segment, r io.Reader, pairs ...*types.Pair) (err error)
    CompleteSegment(seg segment.Segment, pairs ...*types.Pair) (err error)
    AbortSegment(seg segment.Segment, pairs ...*types.Pair) (err error)
}
```

And [Proposal: Split Segmenter](21-split-segmenter.md) added new interfaces:

```go
type IndexSegmenter interface {
    InitIndexSegment(path string, pairs ...Pair) (seg Segment, err error)
    WriteIndexSegment(seg Segment, r io.Reader, index int, size int64, pairs ...Pair) (err error)
}
```

Those changes make Segment much more complex. The more difficult question is our design requires maintaining segment status by storager which makes it is not recovery friendly.

## Proposal

So I propose following changes:

- Don't maintain the internal segment status anymore
- Remove `Segment` interface and add a same name struct instead.
- Add `Segmenter`, `IndexSegmenter` and `OffsetSegmenter` interfaces.

So our interfaces look like following:

```go
type Segment struct {
    Path string
    ID string
}

type Part struct {
	Index int
	Size int64
	ETag string
}

type Segmenter interface {
    ListSegments(pairs ...Pair) SegmentIterator
    InitSegment(path string, pairs ...Pair) Segment
    AbortSegment(seg Segment) error
}

type IndexSegmenter interface {
    Segmenter

    ListIndexSegment(seg Segment) (PartIterator, error)
    CompleteIndexSegment(seg Segment, parts []*Part, pairs ...Pair) error
    WriteIndexSegment(seg Segment, r io.Reader, index int, size int64, pairs ...Pair) (err error)
}

type OffsetSegmenter interface {
    Segmenter

    WriteOffsetSegment(seg Segment, r io.Reader, offset int64, size int64, pairs ...Pair) (err error)
}
```

- `Segmenter`'s `ListSegments` also supports `ListType`.
- `OffsetSegmenter` don't need exclusively `CompleteSegment` operation.
- `OffsetSegmenter` don't have `ListSegment` operation.
- One storager could implement only one Segmenter here.

## Rationale

No rationale content.

## Compatibility

Segment related API could be broken.

## Implementation

Most of the work would be done by the author of this proposal.
