# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [v0.2.0] - 2019-10-17

### Added

- services/qingstor: Implement ListSegments (#7)

### Changed

- segment: Replace identifier and add ListSegments support (#6)

### Fixed

- services/qingstor: Fix ListDir not handled correctly (#2)
- services/qingstor: Fix object size and last modified not filled (#4)
- services/qingstor: Add stat updated at support (#5)

## v0.1.0 - 2019-10-12

### Added

- Add storager and servicer interface.
- Add segment and iterator support.
- Add pair based option and metadata support.
- Add qingstor services.

[Unreleased]: https://github.com/Xuanwo/storage/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/Xuanwo/storage/compare/v0.1.0...v0.2.0
