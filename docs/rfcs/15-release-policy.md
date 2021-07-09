- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2020-03-10
- RFC PR: N/A
- Tracking Issue: N/A

# Proposal: Release policy

## Background

Golang team makes sure the go 1 compatibility in [Go 1 and the Future of Go Programs]:

> It is intended that programs written to the Go 1 specification will continue to compile and run correctly, unchanged, over the lifetime of that specification. 

[storage] should have similar statement about its release policy.

## Proposal

So I propose that [storage] should obey the following policies:

### Semantic Versioning

**[storage] SHOULD follow [Semantic Versioning](https://semver.org/)**

- After `1.0.0` released, [storage] will follow semantic versioning strictly
- All exported items in [storage] except `internal` and `tests` will be included in semantic versioning

[storage] uses [dependbot] to upgrade its dependences automatically, the upgrade process will be like follows:

- [dependbot] will create new PRs to branch [dependence] and merge them after build succeeded
- [storage] will merge branch [dependence] in every release PR
- After release PR get merged, branch [dependence] should be reset to branch master

### Target Golang Versions

**[storage] SHOULD be compatible with the last two golang major versions**

Assuming the current version is go 1.14

- Developer SHOULD develop with go 1.13 **OR** go 1.14
- CI should be passed on **BOTH** go 1.13 **AND** go 1.14
- Any error/bug report on go 1.12 MAY mark as `wontfix`
- New features included in go 1.14 SHOULD NOT be included

## Rationale

None

## Compatibility

None

## Implementation

No code related changes.

[Go 1 and the Future of Go Programs]: https://golang.org/doc/go1compat
[storage]: https://github.com/Xuanwo/storage
[dependbot]: https://dependabot.com/
[dependence]: https://github.com/Xuanwo/storage/tree/dependence
