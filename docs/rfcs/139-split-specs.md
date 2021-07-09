- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-07-08
- RFC PR: [beyondstorage/specs#139](https://github.com/beyondstorage/specs/issues/139)
- Tracking Issue: [beyondstorage/go-storage#627](https://github.com/beyondstorage/go-storage/issues/627)

# GSP-139: Split Specs

Previous Discussion:

- [Move specs back into go-storage](https://github.com/beyondstorage/specs/issues/138)

## Background

`definitions` is a core idea of BeyondStorage's storage abstraction. We use `definitions` to carry the definitions of pairs, infos, and operations. In order to share `definitions` between different languages, we move `definitions` into [specs]. So we have cross-language **definitions**: `/definitions` with language-specific parsing utilities: `/go`, `/rust`.

[specs] is also used as `go-storage`'s RFCs center: we store our approved RFCs in `/rfcs` and `/specs`. In [specs] issues, we discuss ideas before sending formal RFCs.

However, with the growth of the community of BeyondStorage, [specs] lead to more and more confusion. [specs] looks like RFC centers for BeyondStorage, but only been used for `go-storage`. [beyond-tp](https://github.com/beyondstorage/beyond-tp) has its own RFCs storage.

After [GSP-128: Community Organization](./128-community-organization.md), the problem become more serious.

- Should [specs] become a separate project?
- Who will have write access over [specs]?

## Proposal

So I propose to split specs into projects:

- Move `/definitions` and `/go` into `go-storage`.
- Move `/rfcs` and `/spec` into `go-storage`.
- Move `/rust` into `rs-storage`.

After those changes, we can be more focused on building up features instead of caring about behavior across languages. And project's committer and the maintainer will have the ability to approved proposals without contact other teams.

## Rationale

### Why not maintain all projects RFC in specs?

Firstly, it could be confusing.

To not break all our links to GSP, we can't change the directory of `rfcs` and filename in it. So we will have the following directory names:

- rfcs
- btp / rfcs-btp or other directory names

Or filenames:

- 0-example.md
- btp-0-exmaple.md

Secondly, permission is hard to maintain.

- Only give access to a small group: every proposal needs their actions.
- Give access to all committers and maintainers: they may be affected by too many proposals.

### Why not make `definitions` a separate repo?

We used to treat `definitions` as a cross-language registry. But it proved to be over-designed: definition files are easy to maintain and we don't need to share them between different languages. After separation, we can be more focused on building up features instead of caring about behavior across languages.

So we will split current definitions into two copies, one for `go-storage` and the other one for `rs-storage`. Since then, we will not sync them anymore.

### RFCs in project

> Will it be not as clean? There will be many things in issues. BTW, issue labels may be enough for a small organization like us?

I think issue labels are enough for us. We will add a `proposal` label for proposal issues.

> Another small problem is the issue numbers for RFCs will become more discontinuous and bigger?

I think it's not a big problem.

> We cannot distinguish 2 LGTMs for RFCs and 1 LGTM for other PRs automatically?

For now, the restriction is executed by the committer/maintainer manually. In the further, we will assign labels like `need-2-lgtm`.

## Compatibility

To not break our links to existing proposals, we will:

- Move all issues into `go-storage`
    - They will be assigned a new number just like we create new issues in `go-storage`
- Copy `definitions` and `rfcs` to `go-storage`
- Convert `specs` into archived (read-only)

## Implementation

See [Compatibility](#Compatibility).

[specs]: https://github.com/beyondstorage/specs
