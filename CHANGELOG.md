# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

 - Support CSS styling of rendered content
 - `version` subcommand

## [0.3.1] - 2023-11-09

### Added

 - Subcommand functionality
 - `--host` option

### Fixed

 - Re-add `stop` subcommand
 - Cancel server context on server error

## [0.3.0] - 2023-11-08

### Changed

 - Refactor project into server and cmd packages
 - Replace goreleaser with justfile and scripts
 - Require POST method for /shutdown endpoint

## [0.2.0] - 2023-03-12

### Added

 - Enable stopping server using cli flag

## [0.1.1] - 2023-02-27

### Fixed

 - File server now uses embedded static assets

## [0.1.0] - 2023-02-26

Initial version, providing a working proof-of-concept.

[Unreleased]: https://github.com/cluttrdev/showdown/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/cluttrdev/showdown/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/cluttrdev/showdown/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/cluttrdev/showdown/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/cluttrdev/showdown/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/cluttrdev/showdown/releases/tag/v0.1.0

