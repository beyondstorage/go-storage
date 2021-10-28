# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## v3.0.1 - 2021-10-27

### Fixed

- fix(services/s3): Output constant type should not be changed (#971)

## v3.0.0 - 2021-10-23

### Added

- feat(services/s3): Add missing option back (#880)

### Changed

- feat(services/s3): Move services s3 back (#853)
- refactor(services/s3): Cleanup the logic of client init (#876)
- ci(*): Upgrade minimum version to Go 1.16 (#916)

### Upgraded

- build(deps): bump github.com/aws/aws-sdk-go-v2/service/s3 (#940)
- build(deps): bump github.com/aws/aws-sdk-go-v2/config in /services/s3 (#941)

## [v2.5.0] - 2021-10-09

### Added

- test: Add new test case
- ci: Enable auto merge for dependabot
- docs: Add Compatible Services section in README
- feat: Add support for Write Empty File Behavior (#228)
- feat: Add Multipart HTTP Signer support (#237)

### Changed

- docs: Update README (#208)
- ci: Sync github actions (#220, #224, #239)

### Upgraded

- build(deps): Bump github.com/aws/aws-sdk-go to 1.40.58

## [v2.4.0] - 2021-09-03

### Added

- feat: Support offset and size for read (#157)
- feat: Implement CreateLink and setup linker test in go-service-s3 (#175)
- rfcs: RFC-178 Add Virtual Link Support (#178)
- feat: Implement StorageHTTPSigner (#191)
- feat: Add IoCallback for multipart (#196)

### Changed

- ci: Migrate integration tests to self hosted teamcity (#190)

### Fixed

- feat: Resolve null pointer problem when metadata undefined (#185)
- feat: IoCallback in WriteMultipart not supported correctly (#198)

### Upgraded

- build(deps): Bump github.com/aws/aws-sdk-go from 1.40.1 to 1.40.36 (#205)

## [v2.3.0] - 2021-07-16

### Added

- ci: Add gofmt action (#140)
- ci: Add dependabot auto build support (#139)
- ci: Add diff check action (#145)

### Changed

- storage: Update types in service.toml to golang types (#150)
- storage: Implement GSP-654 Unify List Behavior  (#150)

### Fixed

- ci: Fix auto-build not work correctly

### Upgraded

- build(deps): Bump github.com/aws/aws-sdk-go from 1.38.68 to 1.40.1 (#151)

## [v2.2.0] - 2021-06-29

### Added

- storage: Add CreateDir (#126)
- storage: Implement GSP-93 Add ObjectMode Pair (#126)
- storage: Implement GSP-97 Add Restrictions in Storage Metadata (#126)

### Changed

- *: Implement GSP-109 Redesign Features (#126)
- *: Implement GSP-117 Rename Service to System as the Opposite to Global (#126)

### Upgraded

- build(deps): Bump aws-go-sdk to 1.38.68 (#127)

## [v2.1.0] - 2021-06-15

### Added

- *: Implement GSP-87 Feature Gates (#101)

### Fixed

- service: Fix endpoint not supported (#115)
- service: Fix incorrect usage of endpoint (#117)

### Upgraded

- build(deps): Bump go-endpoint to 1.0.1
- build(deps): Bump aws-go-sdk to 1.38.61 (#116)

## [v2.0.0] - 2021-05-24

### Added

- *: Implement GSP-47 Additional Error Specification (#78)
- *: Implement GSP-51 Distinguish Errors by IsInternalError (#83)
- storage: Implement GSP-61 Add object mode check for operations (#89)

### Changed

- storage: Idempotent storager delete operation (#82)
- storage: Implement GSP-62 WriteMultipart returns Part (#88)
- *: Implement GSP-73 Organization rename (#94)

### Upgraded

- build(deps): Bump aws-go-sdk to 1.38.45 (#95)

## [v1.2.0] - 2021-04-24

### Added

- *: Add UnimplementedStub (#61)
- tests: Introduce STORAGE_S3_INTEGRATION_TEST (#65)
- storage: Add configurations and support SSE (#63)
- storage: Implement GSP-40 (#68)

### Changed

- ci: Only run Integration Test while push to master
- storage: Update SSE pair description (#66)

### Upgraded

- build(deps): bump github.com/aws/aws-sdk-go from 1.38.19 to 1.38.23 (#64)
- build(deps): bump github.com/aws/aws-sdk-go from 1.38.23 to 1.38.24 (#67)

## [v1.1.0] - 2021-04-14

### Added

- *: Implement default pair support for service (#16)
- storage: Implement Create (#56)

### Changed

- build(deps): bump github.com/aws/aws-sdk-go to 1.38.19 (#57)

## v1.0.0 - 2021-02-07

### Added

- Implement s3 services.

[v2.5.0]: https://github.com/beyondstorage/go-service-s3/compare/v2.4.0...v2.5.0
[v2.4.0]: https://github.com/beyondstorage/go-service-s3/compare/v2.3.0...v2.4.0
[v2.3.0]: https://github.com/beyondstorage/go-service-s3/compare/v2.2.0...v2.3.0
[v2.2.0]: https://github.com/beyondstorage/go-service-s3/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/beyondstorage/go-service-s3/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-s3/compare/v1.2.0...v2.0.0
[v1.2.0]: https://github.com/beyondstorage/go-service-s3/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/beyondstorage/go-service-s3/compare/v1.0.0...v1.1.0
