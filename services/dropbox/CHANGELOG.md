# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v2.3.0] - 2021-09-13

### Changed

- ci: Enable auto merge for dependabot
- ci: Bump Xuanwo/fetch-metadata
- docs: Update README (#42)
- ci: Cleanup Service Integration Tests (#43)

### Fixed

- fix: Regenerate code

## [v2.2.0] - 2021-07-21

### Added

- ci: Add gofmt action (#26)
- ci: Add dependabot auto build support (#29)
- ci: Add diff check action (#30)

### Changed

- storage: Implement GSP-134 Write Behavior Consistency (#34)

### Fixed

- ci: Fix auto-build not work correctly

### Upgraded

- build(deps): Bump github.com/dropbox/dropbox-sdk-go-unofficial/v6 from 6.0.1 to 6.0.2 (#24)

## [v2.1.0] - 2021-06-29

### Added

- *: Implement GSP-87 Feature Gates (#19)
- storage: Add CreateDir (#22)
- storage: Implement GSP-97 Add Restrictions in Storage Metadata (#22)

### Changed

- *: Implement GSP-109 Redesign Features (#22)
- *: Implement GSP-117 Rename Service to System as the Opposite to Global (#22)

## [v2.0.0] - 2021-05-24

### Added

- docs: Add README
- storage: Add appender support (#12)
- *: Implement GSP-47 & GSP-51 (#13)
- Implement GSP-61 Add object mode check for operations (#16)

### Changed

- *: Implement GSP-73 Organization rename (#17)

### Upgraded

- build(deps): Bump dropbox-sdk-go to 6.0.1 (#95)

## v1.0.0 - 2021-04-24

### Added

- Implement dropbox services.

[v2.3.0]: https://github.com/beyondstorage/go-service-dropbox/compare/v2.2.0...v2.3.0
[v2.2.0]: https://github.com/beyondstorage/go-service-dropbox/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/beyondstorage/go-service-dropbox/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-dropbox/compare/v1.0.0...v2.0.0
