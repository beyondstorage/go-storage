# Segmenter

## AbortSegment

AbortSegment will abort a segment.

Implementer:

- SHOULD return error while caller call AbortSegment without init.

Caller:

- SHOULD call InitIndexSegment before AbortSegment.

## CompleteSegment

CompleteSegment will complete a segment and merge them into a File.

Implementer:

- SHOULD return error while caller call CompleteSegment without init.

Caller:

- SHOULD call InitIndexSegment before CompleteSegment.