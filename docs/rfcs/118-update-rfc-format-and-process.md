- Author: xxchan <xxchan22f@gmail.com>
- Start Date: 2021-06-22
- RFC PR: [beyondstorage/specs#118](https://github.com/beyondstorage/specs/issues/118)
- Tracking Issue: [beyondstorage/specs#121](https://github.com/beyondstorage/specs/issues/121)

# GSP-118: Update RFC Format and Process

- Updates:
  - [GSP-17](./17-proposal-process.md): Updates the RFC spec and process

## Background

Early discussion:
- [beyondstorage/specs#110](https://github.com/beyondstorage/specs/issues/110)

The current proposal process and format is specified by [GSP-17](./17-proposal-process.md) and [spec/2-proposal](../spec/2-proposal.md).

The metadata `status` and `updated_at` are required.

For `status`, the transition `draft` -> `candidate` -> `finished` is troublesome (have to start a PR to update status). In practice, this is not followed (There are 18 rfcs with status `draft`.
As an alternative, we can link an issue to track the implementation status.

For `updated_at`, we often forgot to update this when a proposal was discussed for days. And I think maybe the date information is not very important.

## Proposal

So I propose to update the metadata format and requirements in the RFC [template](./0-example.md) and RFC [spec](../spec/2-proposal.md):
- Remove front matter, put metadata in the text.
- Change metadata requirements:
  - Remove `status`, and add `Tracking Issue`.
  - Remove `updated_at`, and add `Start Date`.
  - Add `RFC PR`.

And specify the RFC process in [README](../README.md).

## Rationale

Principles:

- Little trouble for proposers:
  - Don't need to maintain trivial information like `updated_at`
- Good readability for reviewers and future readers:
  - Attach a link whenever possible: `Tracking Issue` and `RFC PR`.

Reference: [Rust RFC](https://github.com/rust-lang/rfcs) has the following metadata:

```markdown
- Feature Name: (fill me in with a unique ident, `my_awesome_feature`)
- Start Date: (fill me in with today's date, YYYY-MM-DD)
- RFC PR: [rust-lang/rfcs#0000](https://github.com/rust-lang/rfcs/pull/0000)
- Rust Issue: [rust-lang/rust#0000](https://github.com/rust-lang/rust/issues/0000)
```

[Rust RFC website](https://rust-lang.github.io/rfcs/)

## Compatibility

N/A

## Implementation

Migrate plan: 
- Update existing RFCs to the new format. 
  - Set `Start Date` to `updated_at`.
  - Set `Tracking Issue` and `RFC PR` to `N/A` if not applicable.
- RFCs that haven't been merged before this RFC is merged should follow the new format.
