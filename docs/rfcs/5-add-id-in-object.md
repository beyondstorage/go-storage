- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2020-01-03
- RFC PR: N/A
- Tracking Issue: N/A

# Proposal: Add ID in Object struct

## Background

PR [services: Add dropbox basic support](https://github.com/Xuanwo/storage/pull/53) prompts a great problem: Should we have an ID metadata?

The difference between `struct value` and `metadata` is:

- All value in `struct` are required, caller can use them safely.
- All value in `metadata` are optional, caller need to check them before using.

First of all, it's obvious all storage services have an ID for Object:

- For local file system: `ID` could be whole path towards root.
- For object storage: `ID` could be the whole key.
- For dropbox alike SaaS: `ID` could be their ID in business.

Then, user need to access ID for some reason:

- Distinguish files with same name (Some SaaS allow same name file in the same folder.)
- Upper application needs ID for their business logic.

## Proposal

So I propose following changes:

- Add `ID string` in `Object` struct
- Make sure every services filled `Object.ID`

## Rationale

None.

## Compatibility

No breaking changes.

## Implementation

Implemented as Proposal.
