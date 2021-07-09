---
author: Xuanwo <github@xuanwo.io>
status: candidate
updated_at: 2020-03-18
added_by:
- rfcs/17-proposal-process.md
---

# Spec: Proposal

## Format

proposal will be in `markdown` format.

## Title

proposal's filename will be like following:

`<proposal-number>-proposal-name.md`, for example, `16-loose-mode.md`

And they should be referred by `GSP-16` with hyperlink: [GSP-16]

## Metadata

Proposal will have metadata at the beginning.

```markdown
- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2019-11-15
- RFC PR: [beyondstorage/specs#100](https://github.com/beyondstorage/specs/issues/100)
- Tracking Issue: [beyondstorage/go-storage#100](https://github.com/beyondstorage/go-storage/issues/100)

# GSP-0: <proposal name>

- Updates:
  - [GSP-20](./20-abc): Deletes something
- Updated By: 
  - [GSP-10](./10-do-be-do-be-do): Adds something
  - [GSP-1000](./1000-lalala): Deprecates this RFC
```

`Author`, `Start Date`, `RFC PR` and `Tracking Issue` are required.

If this proposal affects or is affected by other proposals, `Updated By`, `Updates`, should be added below the title.

## Code

Related code should be sub directory with the proposal number.

[GSP-16]: ../rfcs/16-loose-mode.md