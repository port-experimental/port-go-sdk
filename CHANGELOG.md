# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Added `WithRetryAttempts` option to Client for configurable retry behavior
- Added `HasMore()` method to `ListResponse` for pagination checking
- Added `ListAll()` and `ListAllBlueprint()` helper methods for automatic pagination
- Added comprehensive godoc documentation to all service methods with context timeout recommendations
- Enhanced examples README with error handling patterns, context usage, and pagination examples
- Added pagination test coverage

### Fixed
- Fixed deprecated `rand.Seed` usage in `pkg/httpx` (replaced with `rand.New`)
- Fixed file handle leak in verbose logger (added `Close()` method to Client)
- Added proper error handling for `io.ReadAll` calls throughout the codebase
- Added error handling for `http.NewRequestWithContext` calls
- Fixed `go.mod` dependency tracking (godotenv is now a direct dependency)

### Added
- Added `Close()` method to Client for proper resource cleanup
- Added input validation for bulk operations (max 20 for BulkUpsert, max 100 for BulkDelete)
- Added `VerifySignature` function to webhooks package for verifying incoming webhook signatures
- Added LICENSE file (MIT License)
- Added Makefile with common development tasks
- Added CHANGELOG.md
- Added package-level godoc documentation
- Added sanitization to verbose logging to prevent secret exposure

### Security
- Sanitized verbose logging to prevent logging of sensitive data (tokens, secrets, etc.)

## [0.1.0] - 2024-XX-XX

### Added
- Initial release of Port Go SDK
- Support for EU and US regions
- Client credentials and API token authentication
- Entity CRUD operations (create, read, update, delete, upsert)
- Bulk entity operations (bulk upsert, bulk delete)
- Entity relations (link, unlink)
- Entity search and aggregation
- Blueprint management
- Data source management
- Automation triggers and execution listing
- User and team management
- Organization management
- Webhook utilities with HMAC SHA256 signing
- Retry logic with exponential backoff
- Verbose logging support

[Unreleased]: https://github.com/port-experimental/port-go-sdk/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/port-experimental/port-go-sdk/releases/tag/v0.1.0

