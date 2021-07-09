- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-06-29
- RFC PR: [beyondstorage/specs#128](https://github.com/beyondstorage/specs/pull/128)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-128: Community Organization

## Background

BeyondStorage has been grown up, and we need to define our community organization to make it clear that everyone's rights and responsibilities.

## Proposal

I propose to add five positions for BeyondStorage, and they are separate in different projects.

For now, there are following projects:

- `go-storage`: Focus on implement storage abstraction for golang, including `go-storage`, `go-service-*`, `go-endpoint` and so on.
- `dm`: Focus on implement data migration services, including `dm`.
- `specs`: Focus on the specs of BeyondStorage, new ideas, new proposals happened here, including `specs`.
- `community`: Focus on the development of BeyondStorage community, including `site`, `community` and tools like `go-community`.

### Leader

- ID: `leader`
- Permissions: `Admin`

Project core management team, involved in roadmap development of major community-related resolutions.

### Maintainer

- ID: `maintainer`
- Permissions: `Maintain`

The planner and designer of the project, with the authority to merge code into the trunk, from the Committer.

### Committer

- ID: `committer`
- Permissions: `Write`

The contributor who is recommended by Maintainer or Leader and has made an outstanding contribution to projects and needs to complete at least one feature or fix major bugs independently.

### Reviewer

- ID: `reviewer`
- Permissions: `Triage`

Advanced Contributor, responsible for reviewing community code, has LGTM (Looks Good To Me) approval for new submissions.

### Contributor

- ID: `contributor`
- Permissions: `Read`

Community contributors with more than one PR have been merged.

## Rationale

N/A

## Compatibility

No code changes.

## Implementation

- Create a repo called `community` in which maintains all members and teams
- Implement team member permission automatically maintain in `go-community`.
