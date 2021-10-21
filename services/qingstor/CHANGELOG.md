# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v3.3.0] - 2021-09-13

### Added

- rfcs: RFC-79 Add virtual link support (#79)
- feat: Implement CreateLink and setup linker test (#81)
- feat: WriteMultipart adds io_callback support (#86)
- feat: Implement StorageHTTPSigner (#88)

### Changed

- feat: Turn Expire into Duration for Reach (#76)
- ci: Enable auto merge for dependabot
- ci: Cleanup Service Integration Tests (#90)
- docs: Update README (#91)

### Fixed

- fix: Fixed append test failures (#84)

### Upgraded

- build(deps): Bump qingstor-sdk-go to version 4.4.0 (#87)

## [v3.2.0] - 2021-07-22

### Added

- ci: Add gofmt action (#62)
- ci: Add diff check action (#65)
- ci: Add dependabot auto build support (#66)

### Changed

- storage: Update types in service.toml to golang types (#71)
- storage: Update append as described in GSP-134 (#71)
- storage: Update list as described in GSP-654 (#71)
- build(deps): Migrate to go-endpoint (#74)

### Fixed

- ci: Fix auto-build not work correctly
- storage: Fix invalid argument for copy and move (#72)
- storage: Fix append behavior (#73)

### Upgraded

- build(deps): Bump github.com/google/uuid from 1.2.0 to 1.3.0 (#61)

## [v3.1.0] - 2021-06-29

### Added

- *: Implement GSP-87 Feature Gates (#53)
- storage: Implement GSP-93 Add ObjectMode Pair (#58)
- storage: Implement GSP-97 Add Restrictions In Storage Metadata (#58)

### Changed

- *: Implement GSP-109: Redesign Features (#58)
- *: Implement GSP-117 Rename Service to System as the Opposite to Global (#58)

### Upgraded

- build(deps): bump github.com/golang/mock from 1.5.0 to 1.6.0 (#56)

## [v3.0.0] - 2021-05-24

### Added

- storage: Add appender support (#40)
- *: Implement GSP-47 & GSP-51 (#46)
- storage: Implement GSP-61 Add object mode check for operations (#49)

### Changed

- service: Use path style instead of vhost (#43)
- service: Fix location not detected correctly (#45)
- storage: Idempotent storager delete operation (#44)
- storage: Implement GSP-62 WriteMultipart returns Part (#47)
- storage: Check if part number is valid when multipart upload (#48)
- *: Implement GSP-73 Organization rename (#51)

## [v2.1.0] - 2021-04-24

### Added

- *: Implement proposal unify object metadata (#25)
- storage: Normalize iterator next function names (#27)
- pair: Implement default pair support for service (#29)
- *: Set default pair when init (#31)
- storage: Implement Create API (#33)
- storage: Set multipart attributes when create multipart (#34)
- *: Add UnimplementedStub (#35)
- storage: Implement SSE support (#37)
- tests: Introduce STORAGE_QINGSTOR_INTEGRATION_TEST (#39)
- storage: Implement GSP-40 (#41)

### Changed

- storage: Clean up next page logic
- build: Make sure integration tests has been executed
- docs: Migrate zulip to matrix
- docs: Remove zulip
- ci: Only run Integration Test while push to master
- storage: Rename SSE related pairs to meet GSP-38 (#38)

### Fixed

- storage: Fix multipart integration tests (#36)

### Removed

- *: Remove parsed pairs pointer (#28)

### Upgrade

- build(deps): bump github.com/qingstor/qingstor-sdk-go/v4 (#26)

## [v2.0.0] - 2021-01-17

### Added

- tests: Add integration tests (#17)
- storage: Implement Fetcher (#19)
- storage: Implement proposal Unify List Operation (#20)
- *: Implement Segment API Redesign (#21)
- storage: Implement proposal Object Mode (#22)

### Changed

- Migrate to go-storage v3 (#23)

## v1.0.0 - 2020-11-12

### Added

- Implement qingstor services.

[v3.3.0]: https://github.com/beyondstorage/go-service-qingstor/compare/v3.2.0...v3.3.0
[v3.2.0]: https://github.com/beyondstorage/go-service-qingstor/compare/v3.1.0...v3.2.0
[v3.1.0]: https://github.com/beyondstorage/go-service-qingstor/compare/v3.0.0...v3.1.0
[v3.0.0]: https://github.com/beyondstorage/go-service-qingstor/compare/v2.1.0...v3.0.0
[v2.1.0]: https://github.com/beyondstorage/go-service-qingstor/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-qingstor/compare/v1.0.0...v2.0.0
