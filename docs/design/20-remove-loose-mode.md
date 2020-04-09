---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2020-04-09
---

# Proposal: Remove loose mode

## Background

In order to resolve the problem of pair handle behavior is inconsistent, we introduce loose mode in proposal [16-loose-mode](./16-loose-mode.md). However, enabling the loose mode introduces more problems.

First, we need to figure out issue [error could be returned too early while in loose mode](https://github.com/Xuanwo/storage/issues/233). We have check whether the error can be ignored in loose mode.

Then, nearly no people want the API call move on when the behavior is inconsistent. Any inconsistent behavior should be handled with returning error.

## Proposal

So I proposal to remove loose mode and return error for all inconsistent behavior.

## Rationale

None.

## Compatibility

Loose mode has been removed

## Implementation

Most of the work would be done by the author of this proposal.
