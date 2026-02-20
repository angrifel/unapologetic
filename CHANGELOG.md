# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `release-management.sh` script consolidating release management commands: `last-release-tag`, `last-release-rev`, `semver-prerelease`, `prepare-changelog-for-release`, `extract-version-notes`, and `prepare-prerelease-and-create-draft-release-on-side-tag`

### Changed

- Makefile `semver-prerelease` build metadata now includes a UTC timestamp (`YYYYMMDD-HHMMSS`) prepended to the short git hash
- Claude Code `update-changelog` skill updated to use `./release-management.sh last-release-rev` for base commit resolution, `git diff --name-only` to list changed files, and `git diff` per file for detailed examination
- GitHub Actions release workflow updated to use `./release-management.sh extract-version-notes` instead of `make changelog-release-notes`

### Removed

- Makefile `git-diff-since-last-release` target
- Makefile `semver-prerelease` target (moved to `release-management.sh semver-prerelease`)
- Makefile `prepare-release` target (moved to `release-management.sh prepare-changelog-for-release`)
- Makefile `changelog-release-notes` target (moved to `release-management.sh extract-version-notes`)

## [v0.1.0] - 2026-02-18

### Added

- `httpaux` package with HTTP response cloning (`CloneHTTPResponseWithBody`), body buffering (`BufferResponseBody`), and `RoundTripperFunc` adapter for using functions as `http.RoundTripper`
- `httpspy` package with HTTP request and response capture functions supporting configurable header censoring and optional body collection
- `ioaux` package with I/O function adapters (`ReaderFunc`, `CloserFunc`) and `ReadSeekCloser` for adding seek capabilities to `io.ReadCloser`
- `iospy` package with witness wrappers (`WitnessReader`, `WitnessCloser`) for recording I/O calls, `ReaderWithEOFError` for custom EOF replacement, and `LimitReaderWithError`
- Package-level documentation (`doc.go`) and examples for all public packages
- GitHub Actions CI workflow for build, lint, and test across multiple Go versions and OS targets
- GitHub Actions release workflow with source archives (tar.gz, zip), SBOM generation, SHA-256 checksums, and build provenance attestations
- Pre-commit git hook running lint and test
- Makefile with targets for testing, linting, module tidying, local GitHub Actions execution, and release diff helpers
- Makefile `prepare-release` target for automating changelog version bumps, committing, and tagging releases
- Makefile `semver-prerelease` target for generating pre-release version strings
- Claude Code `update-changelog` skill for automated changelog updates
- MIT license
- README with project overview, Go Reference documentation link, CI status badge, and license badge

[Unreleased]: https://github.com/angrifel/unapologetic/compare/v0.1.0...update-commands
[v0.1.0]: https://github.com/angrifel/unapologetic/compare/ea00b4371869d02656bbd97841caff4a76bc451d...v0.1.0
