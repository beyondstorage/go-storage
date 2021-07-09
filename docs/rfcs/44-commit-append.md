- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-04-26
- RFC PR: [beyondstorage/specs#44](https://github.com/beyondstorage/specs/issues/44)
- Tracking Issue: N/A

# AOS-44: Add CommitAppend in Appender

## Background

`Appender` is designed for `append` operation. We need to use `CreateAppend` to create an appendable object, and `WriteAppend` to write data at the tail of appendable object.

This design implies the fact that: an appendable object is readable. But it's not true for every service. Service like dropbox needs a `close` operation to mark this appends process has been finished. And before the `close` has been called, the object is not readable for user. 

## Proposal

So I propose to add `CommitAppend` for `Appender`.

`CommitAppend` will commit and finish an append process.

For those services that commit in every `WriteAppend` operation, they don't need to do anything in this operation. For those services that requires an explict close operation, they can do this job in `CommitAppend`.

After `CommitAppend` been introduced, users SHOULD call `CommitAppend` to finish an append process.

Company with `CommitAppend`, we will introduce following object metadata to allow user get the limitation of storage service:

- `append-number-maximum`: Maximum append numbers.
- `append-size-maximum`: Maximum size for single `WriteAppend` operation. 
- `append-total-size-maximum`: Maximum total size for the object.

Service SHOULD provide these values if there are limitations for append.

User MAY check these values while doing append operations. If the values isn't exist, user can append without limitation. The service may return an error regardless of whether the check is performed.

## Rationale

N/A

## Compatibility

The proposal is compatible with appender.

## Implementation

N/A
