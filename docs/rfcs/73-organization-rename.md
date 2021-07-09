- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-05-19
- RFC PR: [beyondstorage/specs#73](https://github.com/beyondstorage/specs/issues/73)
- Tracking Issue: N/A

# GSP-73: Organization Rename

## Background

After voting on the [Organization Name Changing Proposal], our organization has been renamed to `Beyond Storage`. With this change, all our repositories don't match with the new name, and some links couldn't access anymore.

## Proposal

So I propose that all our repositories should reflect this change.

- Organization name or address should be updated in infrastructure.
- Module names, dependencies and related references should be updated in services.
- Organization information and related name should be updated in docs.
- Outside links should be accessed normally.

## Rationale
N/A

## Compatibility
All services, applications and infrastructure could be affected.

## Implementation

- Infrastructure
    - forum
    - chat
    - matrix
    - CI
- specs
    - Rename all the proposal names `AOS-*` to `GSP-*`.
    - Rename the interface `AosError` to `InternalError` and the method `IsAosError` to `IsInternalError` proposed in [GSP-51].
    - Update the module name for golang and rust, and all the reference links.
    - Update README doc.
- go-storage:
    - Update module name and import paths.
    - Rename `AosError` to `InternalError` and update the implementation.
    - Update README doc.
    - Bump to `github.com/beyondstorage/go-storage/v4`.
- go-integration-test:
    - Update module name and import paths.
    - Bump to `github.com/beyondstorage/go-integration-test/v4`.
- go-service-*:
    - Update the import paths, and the reference links.
    - Update module name and bump to a new major version.
- go-service-example:
    - Migrate to go-storage v4.0.0.
- go-storage-example:
    - Migrate to go-storage v4.0.0.
- go-coreutils:
    - Update the import paths and bump to a new version.
- dm:
    - Update import paths.
    - Update module name and bump to a new version.
- noah:
    - Update import paths.
    - Update module name and bump to a new version.
- rs-storage:
    - Update the package description and dependencies.
- site:
    - Add a blog to introduce why and what's the meaning of `Beyond Storage`.
    - Update docs, include reference links and organization related content.
- Update outside links.
    
[Organization Name Changing Proposal]: https://forum.aos.dev/t/organization-name-changing-proposal/38
[GSP-51]: ./51-distinguish-errors-by-isaoserror.md
