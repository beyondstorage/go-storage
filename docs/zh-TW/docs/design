---
author: Xuanwo <github@xuanwo.io>
status: candidate
updated_at: 2020-04-27
---

# Proposal: Split Segmenter

## Background

In proposal [18-return-segment-interface-instead]((./18-return-segment-interface-instead.md)), we introduced `segment.Segmenter` interface, but also introduce a new problem: user don't know which type of segment that service supported.

In the `Rationale` section of [18-return-segment-interface-instead]((./18-return-segment-interface-instead.md)), we have already discussed two implementations:

1. `OffsetBasedSegmenter` vs `IndexBasedSegmenter`
2. `WriteSegmentViaIndex` vs `WriteSegmentViaOffset`

We said:

> No.1 implementation's problem is these add too much work for services to support both of them.

But maybe this is the way, we can build new segmenter upon this idea.

## Proposal

Based on the idea of `OffsetBasedSegmenter` and `IndexBasedSegmenter`, I propose following changes:

Add new interface `IndexSegmenter` that implement following functions:

- `InitIndexSegment(path string, pairs ...*types.Pair) (seg segment.Segment, err error)`
- `WriteIndexSegment(seg segment.Segment, r io.Reader, index int, size int64, pairs ...*types.Pair) (err error)`
- `CompleteSegment(seg segment.Segment, pairs ...*types.Pair) (err error)`
- `AbortSegment(seg segment.Segment, pairs ...*types.Pair) (err error)`

Extract `CompleteSegment` and `AbortSegment` into a non-exported interface `segmenter`.

Embed `segmenter` into `IndexSegmenter`, `DirSegmentsLister` and `PrefixSegmentsLister`

Remove `Segmenter`

### Add new interface `IndexSegmenter`

It's easy to find out that we don't need to implement an entire new segmenter interface. After introducing `segment.Segmenter` interface, we can reuse already implemented `CompleteSegment` and `AbortSegment`.

After this change, the workload for service implementer reduced as the same. If new segment method needs to implemented, we only need to add `OffsetSegmenter`.

### Extract non-exported interface and embed into others

We already added `DirSegmentsLister` and `PrefixSegmentsLister`, but list only interface is not useful, we always need to convert to `Segmenter`. After extracting `CompleteSegment` and `AbortSegment` into a non-exported interface, we can make `DirSegmentsLister` and `PrefixSegmentsLister` more useful.

### Remove `Segmenter`

Yes, I know, I introduced an API breaking changes here.

I'm sorry, but I don't have the time and energy to maintain both v1 and v2 for now. Since [storage](https://github.com/Xuanwo/storage) is nearly no one except me used, I think it's safe to introduce this change here and not break others.

I'm aware the risk, and I will be responsible for this.

## Rationale

None.

## Compatibility

All API call related `Segmenter` will be borken.

## Implementation

Most of the work would be done by the author of this proposal.