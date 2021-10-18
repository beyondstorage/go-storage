# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v2.3.0] - 2021-09-13

### Changed

- ci: Enable auto merge on dependabot
- ci: Cleanup Service Integration Tests (#57)
- docs: Update README (#59)

### Upgraded

- ci: Upgrade Xuanwo/fetch-metadata
- build(deps): bump github.com/tencentyun/cos-go-sdk-v5 to 0.7.31 (#58)

## [v2.2.0] - 2021-07-17

### Added

- ci: Add gofmt action (#37)
- ci: Add dependabot auto build support (#40)
- ci: Add diff check action (#41)

### Changed

- storage: Update types in service.toml to golang types (#46)
- storage: Implement GSP-654 Unify List Behavior (#46)

### Fixed

- ci: Fix auto-build not work correctly

## [v2.1.0] - 2021-06-29

### Added

- *: Implement GSP-87 Feature Gates (#27)
- storage: Add CreateDir (#28)
- storage: Implement GSP-97 Add Restrictions in Storage Metadata (#32)

### Changed

- *: Implement GSP-109 Redesign Features (#32)
- *: Implement GSP-117 Rename Service to System as the Opposite to Global (#32)

### Upgraded

- build(deps): bump github.com/tencentyun/cos-go-sdk-v5 to 0.7.27 (#34)

## [v2.0.0] - 2021-05-24

### Added

- storage: Implement SSE support (#17)
- services: implement GSP-47 & GSP-51 (#21)
- storage: Implement multipart support (#23)

### Changed

- storage: Idempotent storager delete operation (#20)
- *: Implement GSP-73 Organization rename (#24)

## [v1.1.0] - 2021-04-24

### Added

- pair: Implement default pair support for service (#4)
- storage: Implement Create API (#13)
- *: Add UnimplementedStub (#15)
- tests: Introduce STORAGE_COS_INTEGRATION_TEST (#16)
- tests: Add docs for how to run tests 
- storage: Implement GSP-40 (#18)

### Changed

- docs: Migrate zulip to matrix
- build: Fix build scripts
- ci: Only run Integration Test while push to master

### Upgraded

- build(deps): bump github.com/tencentyun/cos-go-sdk-v5 from 0.7.19 to 0.7.24

## v1.0.0 - 2021-02-08

### Added

- Implement cos services.

[v2.3.0]: https://github.com/beyondstorage/go-service-cos/compare/v2.2.0...v2.3.0
[v2.2.0]: https://github.com/beyondstorage/go-service-cos/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/beyondstorage/go-service-cos/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-cos/compare/v1.1.0...v2.0.0
[v1.1.0]: https://github.com/beyondstorage/go-service-cos/compare/v1.0.0...v1.1.0
