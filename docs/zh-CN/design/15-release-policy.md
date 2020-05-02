---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2020-03-10
---

# Proposal: Release policy

## Background

Golang team makes sure the go 1 compatibility in [Go 1 and the Future of Go Programs](https://golang.org/doc/go1compat):

> It is intended that programs written to the Go 1 specification will continue to compile and run correctly, unchanged, over the lifetime of that specification.

[storage](https://github.com/Xuanwo/storage) should have similar statement about its release policy.

## Proposal

So I propose that [storage](https://github.com/Xuanwo/storage) should obey the following policies:

### Semantic Versioning

**[storage](https://github.com/Xuanwo/storage) SHOULD follow [Semantic Versioning](https://semver.org/)**

- After `1.0.0` released, [storage](https://github.com/Xuanwo/storage) will follow semantic versioning strictly
- All exported items in [storage](https://github.com/Xuanwo/storage) except `internal` and `tests` will be included in semantic versioning

[storage](https://github.com/Xuanwo/storage) uses [dependbot](https://dependabot.com/) to upgrade its dependences automatically, the upgrade process will be like follows:

- [dependbot](https://dependabot.com/) will create new PRs to branch [dependence](https://github.com/Xuanwo/storage/tree/dependence) and merge them after build succeeded
- [storage](https://github.com/Xuanwo/storage) will merge branch [dependence](https://github.com/Xuanwo/storage/tree/dependence) in every release PR
- After release PR get merged, branch [dependence](https://github.com/Xuanwo/storage/tree/dependence) should be reset to branch master

### Target Golang Versions

**[storage](https://github.com/Xuanwo/storage) SHOULD be compatible with the last two golang major versions**

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