# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## v3.0.0 - 2021-10-23

### Changed

- feat(services/gcs): Move services gcs back (#896)
- ci(*): Upgrade minimum version to Go 1.16 (#916)

### Upgraded

- build(deps): bump cloud.google.com/go/storage in /services/gcs (#899)
- build(deps): bump google.golang.org/api in /services/gcs (#926)

## [v2.3.0] - 2021-09-13

### Added

- feat: Use credentials from the environment (#58)

### Changed

- ci: Enable auto merge for Dependabot
- ci: Cleanup Service Integration Tests (#67)
- docs: Update README (#68)

### Fixed

- fix: Regenerate code

### Refactor

- refactor: Rewrite credential parse in more directly way (#59)

### Upgraded

- ci: Upgrade fetch-metadata
- build(deps): Bump google.golang.org/api from 0.50.0 to 0.56.0 (#63)

## [v2.2.0] - 2021-07-21

### Added

- ci: Add gofmt action (#41)
- ci: Add diff check action (#44)
- ci: Add dependabot auto build support (#45)

### Changed

- storage: Update types in service.toml to golang types (#49)
- storage: Implement GSP-654 Unify List Behavior (#49)

### Fixed

- ci: Fix auto-build not work correctly

### Upgraded

- build(deps): Bump cloud.google.com/go/storage from 1.15.0 to 1.16.0 (#38)
- build(deps): Bump google.golang.org/api from 0.49.0 to 0.50.0 (#39)

## [v2.1.0] - 2021-06-29

### Added

- *: Implement GSP-87 Feature Gates (#31)
- storage: Implement GSP-93 Add ObjectMode Pair (#35)

### Changed

- *: Implement GSP-109: Redesign Features (#35)
- *: Implement GSP-117 Rename Service to System as the Opposite to Global (#35)

### Upgraded

- build(deps): Bump google.golang.org/api to 0.49.0 (#34)

## [v2.0.0] - 2021-05-24

### Added

- pair: Add gcs SSE-KMS support (#24)
- *: Implement GSP-47 & GSP-51 (#27)

### Changed

- storage: Idempotent storager delete operation (#26)
- *: Implement GSP-73 Organization rename (#29)

### Upgraded

- build(deps): Bump google.golang.org/api to 0.47.0 (#28)

## v1.0.0 - 2021-04-24

### Added

- Implement gcs service

[v2.3.0]: https://github.com/beyondstorage/go-service-gcs/compare/v2.2.0...v2.3.0
[v2.2.0]: https://github.com/beyondstorage/go-service-gcs/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/beyondstorage/go-service-gcs/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-gcs/compare/v1.0.0...v2.0.0
