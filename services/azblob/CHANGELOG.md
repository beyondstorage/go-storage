# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## v3.0.0 - 2021-10-23

### Added

- feat: Add support for Write Empty File Behavior (https://github.com/beyondstorage/go-service-azblob/pull/54)

### Changed

- feat(services/azblob): Move services azblob back (#875)
- ci(*): Upgrade minimum version to Go 1.16 (#916)

## [v2.3.0] - 2021-09-13

### Changed

- ci: Auto merge updates from dependabot if CI passed
- ci: Allow PRs that not act by dependabot
- ci: Cleanup Service Integration Tests (#50)
- docs: Update README (#51)

## [v2.2.0] - 2021-07-21

### Added

- ci: Add gofmt action (#33)
- ci: Add dependabot auto build support (#36)
- ci: Add diff check action (#37)

### Changed

- storage: Update types in service.toml to golang types (#32)
- storage: Implement GSP-654 Unify List Behavior (#32)

### Fixed

- ci: Fix checkout wrong commit in auto-build (#40)
- ci: Fix auto build not work correctly

### Upgraded

- build(deps): Bump github.com/Azure/azure-storage-blob-go from 0.13.0 to 0.14.0 (#30)

## [v2.1.0] - 2021-06-29

### Added

- *: Implement GSP-87 Feature Gates (#25)
- storage: Implement CreateDir (#26)
- storage: Implement GSP-97 Add Restrictions in Storage Metadata (#27)

### Changed

- *: Implement GSP-109 Redesign Features (#27)
- *: Implement GSP-117 Rename Service to System as the Opposite to Global (#27)

## [v2.0.0] - 2021-05-24

### Added

- storage: Add SSE part in formatFileObject (#17)
- storage: Add appender support (#18)
- storage: Implement CommitAppend (#19)
- *: Implement GSP-47 & GSP-51 (#21)
- storage: Implement GSP-61 Add object mode check for operations (#22)

### Changed

- docs: Update links to matrix
- storage: Idempotent storager delete operation (#20)
- *: Implement GSP-73 Organization rename (#23)

## [v1.1.0] - 2021-04-24

### Added

- storage: Add AccessTier, ContentType, ContentMD5 support for write (#4)
- *: Implement default pair support for service (#6)
- storage: Implement Create API (#10)
- *: Add UnimplementedStub (#11)
- tests: Introduce STORAGE_AZBLOB_INTEGRATION_TEST (#12)
- storage: Implement SSE support (#13)
- storage: Implement GSP-40 (#15)

### Changed

- ci: Only run Integration Test while push to master

## v1.0.0 - 2021-02-18

### Added

- Implement azblob services.

[v2.3.0]: https://github.com/beyondstorage/go-service-azblob/compare/v2.2.0...v2.3.0
[v2.2.0]: https://github.com/beyondstorage/go-service-azblob/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/beyondstorage/go-service-azblob/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-azblob/compare/v1.1.0...v2.0.0
[v1.1.0]: https://github.com/beyondstorage/go-service-azblob/compare/v1.0.0...v1.1.0
