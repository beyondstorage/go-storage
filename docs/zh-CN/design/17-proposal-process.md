---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2020-03-18
---

# Proposal: Proposal process

## Background

[storage][]'s development is proposal driven. We need to explain why and how we make changes to [storage][] so that we can understand why we are at here.

## Proposal

So I propose following process procedure:

**Simple changes**

- Send PR directly.

**BUG Fix**

- Create an issue
- Send a related PR to resolve it

**Big changes**

- Create an issue
- Send a PR with proposal
- Implement proposal

All steps do not need to be done by the same person. For example, issue could be created by user A, and proposal written by user B, and implemented by user C.

Changes level could be increased while needed. For example, user A sends an one line simple change, but it found out that we need a whole refactor on this package. At this time, we will need to follow the **Big changes** procedure.

Proposal's spec will be presented in spec [2-proposal][].

## Rationale

None

## Compatibility

None

## Implementation

No code related changes.

[storage]: https://github.com/Xuanwo/storage
[2-proposal]: ../spec/2-proposal.md