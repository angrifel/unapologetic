# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- Claude Code `update-changelog` skill for automated changelog updates
- MIT license

[Unreleased]: https://github.com/angrifel/unapologetic/compare/ea00b4371869d02656bbd97841caff4a76bc451d...master
