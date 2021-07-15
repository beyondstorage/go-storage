# RFCs

This folder maintains RFCs for go-storage.

## Explanation

### `RFC`/`Proposal`

We use `RFC` or `Porposal` to represent any non-trivial changes to our codebase, including add, change, deprecate, remove public APIs, behavior changes, or huge impact code refactor.

### `GSP`

We use `GSP`(a.k.a `go-storage proposal`) to represent proposals that apply to `go-storage`.

To get a GSP approved, we need at least two reviewers' approval, and one of them should be committer or maintainer.

- [All Reviewers](https://github.com/orgs/beyondstorage/teams/go-storage-reviewer)
- [All Committers](https://github.com/orgs/beyondstorage/teams/go-storage-committer)
- [All Maintainers](https://github.com/orgs/beyondstorage/teams/go-storage-maintainer)

We use GitHub pull request number as the GSP number.

For historic reasons:

- GSP-1 to GSP-25 numbered by auto-increment id.
- GSP-38 to GSP-139 numbered by PR numbers in [specs](https://github.com/beyondstorage/specs).
- New proposals after GSP-139 will numbered by PR number in this repo (changed by [GSP-139](./139-split-specs.md)).

## Process

- (Optional) Submit an issue in `go-storage` or [forum](https://forum.beyondstorage.io/) as a pre-proposal to discuss the idea.
- Copy `docs/rfcs/0-example.md` to `docs/rfcs/0-my-feature.md` (where "my-feature" is descriptive). Don't assign an RFC number yet; This is going to be the PR number.
- Fill in the RFC.
- Submit a pull request. Use the issue number of the PR to update the filename.
- Discuss the RFC pull request with reviewers and make edits. If the proposal is huge or complex, the reviewer could require the author to implement a demo.
- If at least 2 reviewers approve the RFC, it is accepted.
- Before an accepted RFC gets merged, a tracking issue should be opened in the corresponding repo, and the RFC metadata should be updated. If the author doesn't do so in time, a maintainer can help complete it and then merge the PR.

### After an RFC is approved

Once the feature is implemented in `go-storage`, the tracking issue can be closed. At the same time, tracking issues for all the services should be created (can be automated using `go-community`). 

Now the feature is viewed as implemented, and `go-storage` CAN release a new version.

Regarding as the service tracking issues:
- If it is only required to update `go-storage`, the task CAN be delayed when `go-storage` releases a new version and automatically done by dependabot.
- If manual work is required, implement it correspondingly. This CAN be done immediately or when `go-storage` releases a new version.
  
  It is possible that some other tasks (only need to update `go-storage`) are done by the way, and their tracking issues can be closed.