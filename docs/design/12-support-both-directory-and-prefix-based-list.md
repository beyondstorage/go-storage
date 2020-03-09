---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2020-03-02
updates:
  - design/2-use-callback-in-list-operations.md
---

# Proposal: Support both directory and prefix based list

## Background

We removed recursive support in [2-use-callback-in-list-operations], and in the same proposal, we outlined that:

> Directory based storage will only list one directory, and prefix based storage will only list one prefix without a delimiter.

It's easy for implementation, but doesn't meet demands. There are two main scenarios:

- Work with directories on prefix based storage services, like sync files
- List without delimiters so that we can list faster if we don't care about directories

As we can see, the former proposal fix the list without delimiters scenario but eliminate the possibility that works with directories on prefix based storage services.

## Proposal

So I propose following changes:

- Add `ObjectFunc` for prefix based list
- Treat `FileFunc` and `DirFunc` as directory based list

Storager's behavior will keep following rules:

- `ObjectFunc` can't be passed with `FileFunc`
- `ObjectFunc` can't be passed with `DirFunc`
- `directory based list support` is required
- `prefix based list support` is optional

And we will do validate in `internal/cmd/service`, so that implementers don't need to check them.

## Rationale

None

## Compatibility

No breaking changes.

## Implementation

Most of the work would be done by the author of this proposal.

[2-use-callback-in-list-operations]: ./2-use-callback-in-list-operations.md