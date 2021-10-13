- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-10-13
- RFC PR: [beyondstorage/go-storage#840](https://github.com/beyondstorage/go-storage/issues/840)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-840: Convert to monorepo

- Previous Discussion: [Proposal: Convert go-storage to monorepo](https://forum.beyondstorage.io/t/topic/251)

## Background

go-storage used to be a mono repo a long time ago, at the time when the project just starts. We have all our services at the same repo, under the same version.

However, this layout soon proved to be a huge burden for us. There are the following problems:

- go-storage will have too many dependencies, users don't want to include services that they don't want.
- We can't introduce changes: any API changes will require a major version bump.

So split go-storage into repos, now we have:

- go-storage
- go-service-xxx (up to 26 repos)
- go-endpoint
- go-credential
- ...

After this split, our community is struggling on other problems:

- Too many repos make it hard to track issues, we have to switch between different repos.
- Separate repos make it complex to do automation, we spent a log of time to tag/release, and so on.
- Separate repos make it hard to do huge refactor, we have to submit different PRs again and again.
- Separate repos make it hard to sync changes, we have to wait for release for go-storage.
- Add new services need to add a new repo which needs maintainers operations.
- ...

So I propose to convert go-storage to a mono repo with multiple go modules.

## Proposal

We will have all services inside the same repo, but with different go modules.

```text
go-storage
├── cmd
│   └── definitions
│       ├── bindata
│       └── testdata
├── definitions
├── pairs
├── pkg
│   ├── credential
│       ├── go.mod
│   ├── endpoint
│       ├── go.mod
│   ├── fswrap
│       ├── go.mod
│   ├── headers
│       ├── go.mod
│   ├── httpclient
│       ├── go.mod
│   ├── iowrap
│       ├── go.mod
│   └── randbytes
│       ├── go.mod
├── services
│   ├── s3
│       ├── go.mod
│   ├── gcs
│       ├── go.mod
│   ├── oss
│       ├── go.mod
├── tests
└── types
```

The module path will be:

- `github.com/beyondstorage/go-storage/v4/types` -> `beyondstorage.io/go/v4/types`
- `github.com/beyondstorage/go-service-s3/v2` -> `beyondstorage.io/go/services/s3/v2`
- `github.com/beyondstorage/go-service-gcs/v2` -> `beyondstorage.io/go/services/gcs/v2`

## Rationale

### How to release a version?

We will adopt [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) and to [Release Please](https://github.com/googleapis/release-please) to automate the release.

### How to update dependencies?

We will use [renovate](https://github.com/apps/renovate) to automate the dependencies update. As shown at [chore(all): update all](https://github.com/googleapis/google-cloud-go/pull/4971), renovate can upgrade the whole repo's dependencies at the same time. Which is far more suitable for us.

### How to run integration tests?

We will develop a simple tool like [changefinder](https://github.com/googleapis/google-cloud-go/tree/master/internal/actions/cmd/changefinder) to find what has been changed and only run corresponding tests.

## Compatibility

This is a huge change, so we will change the import path and bump into a major version of all our libraries. There will be no API changes in this GSP.

After updating the import path and running `go mod tidy`, all code should be work as expected.

## Implementation

No API changes.
