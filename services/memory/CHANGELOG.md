# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v0.3.0] - 2021-09-13

### Added

- feat: Implement IoCallback support (#23)
- feat: Add storage_features and default_storage_pairs support

### Changed

- ci: Enable auto merge for Dependabot
- ci: Upgrade fetch-metadata
- docs: Update README (#27)

## [v0.2.0] - 2021-08-23

### Added

- ci: Add intergration tests (#15)
- feat: Add support for Copier, Mover, Appender, Direr (#14)

### Fixed

- fix: stat root path will return ErrObjectNotExist (#17)
- fix: Object size calculated incorrectly while short write (#21)

### Refactor

- refactor: Add name field in object (#19)

## v0.1.0 - 2021-07-26

### Added

- Implement memory services.

[v0.3.0]: https://github.com/beyondstorage/go-service-memory/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/beyondstorage/go-service-memory/compare/v0.1.0...v0.2.0 
