---
author: Xuanwo <github@xuanwo.io>
status: candidate
updated_at: 2020-04-09
updates:
  - design/2-use-callback-in-list-operations.md
deprecates:
  - design/12-support-both-directory-and-prefix-based-list.md
---

# Proposal: Split storage list

## Background

proposal [support both directory and prefix based list][] has been proved to be a failure by practice. In this proposal, we introduce `ObjectFunc` for prefix based list support, and add many restriction for the usage of `FileFunc`, `DirFunc` and `ObjectFunc`. The problem is user don't know whether this storage service is prefix based or directory based. So they always fallback to the directory based list method which is not suffcient for object storage service.

## Proposal

So I propose following changes:

- Split `List` into `ListDir` and `ListPrefix`
- Remove `List` from `Storager`
- Add interface `DirLister` for `ListDir`
- Add interface `PrefixLister` for `ListPrefix`

So user need to assert to interface `DirLister` to use `ListDir`.

At the same time, we should:

- Rename `ListSegments` to `ListPrefixSegments` to match prefix changes
- Remove `ListSegments` from `Segmenter`
- Add interface `PrefixSegmentsLister` for `ListSegments`

## Rationale

None.

## Compatibility

All API call to `List` will be broken.

## Implementation

Most of the work would be done by the author of this proposal.

[support both directory and prefix based list]: ./12-support-both-directory-and-prefix-based-list.md