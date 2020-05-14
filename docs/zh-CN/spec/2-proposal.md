---
author: Xuanwo <github@xuanwo.io>
status: candidate
updated_at: 2020-03-18
added_by:
  - design/17-proposal-process.md
---

# Spec: Proposal

## Format

proposal will be in `markdown` format.

## Title

proposal's filename will be like following:

`<proposal-number>-proposal-name.md`, for example, `16-loose-mode.md`

And they should be referred by `16-loose-mode` or `design/16-loose-mode.md`

## Metadata

Proposal will have metadata as front meta to carry more info.

Take [3-support-service-init-via-config-string][] as example:

```yaml
---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2019-12-23
updated_by:
  - design/4-credential-refactor.md
  - design/9-remove-storager-init.md
deprecated_by:
  - design/13-remove-config-string.md
---
```

`auther`, `status` and `updated_at` are required.

`auther` should be in format: `Name <Email>`

If this proposal affects other proposals, `updated_by`, `updates`, `deprecated_by`, `deprecates` and so on should also be added.

## Status

Proposal has following status: `draft`, `candidate`, `finished`.

- Proposal just created, but not implemented: `draft`
- Proposal implemented but doesn't have final confirmation: `candidate`
- Proposal has been included in a minor version: `finished`

So when we send a PR, we should set this proposal to `draft`.

After we implemented a proposal, we should set this proposal to `candidate`, at this stage, we can update our implementation based on actual feedback.

When we decide to have a minor release, we should take all `candidate` proposal into consideration. Dissatisfied proposal will be removed, and accepted proposals will be marked into `finished`.

After a proposal marked `finished`, we can't change it's content anymore (metadata could also be changed if updated). We need to submit a new proposal to change some proposal's behavior.

## Code

Related code should be sub directory with the proposal number.

[3-support-service-init-via-config-string]: ../design/3-support-service-init-via-config-string.md