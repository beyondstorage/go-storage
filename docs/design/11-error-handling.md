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

- Add spec [1-error-handling](../spec/1-error-handling.md) to normalize error handling across the whole lib
- All error handling related code should be refactored

## Rationale

- <https://blog.golang.org/error-handling-and-go>
- <http://joeduffyblog.com/2016/02/07/the-error-model/>

## Compatibility

Error returned by [storage] could be changed.

## Implementation

Most of the work would be done by the author of this proposal.

[storage]: https://github.com/Xuanwo/storage
