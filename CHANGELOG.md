# Changelog

All notable changes to this project will be documented in this file.

## v0.2.1 - 2025-12-06

### Added
- Introduced `pkg/version` so the default user agent string and release metadata derive from a single source of truth.
- Added `scripts/release.sh` to automate semantic version bumps, README updates, and tag preparation.

### Changed
- Swapped the dependency on `github.com/joho/godotenv` with a lightweight in-repo `.env` loader (`pkg/config/dotenv.go`) to keep module consumers from hitting `go.sum` validation errors.
- Updated blueprint scaffolding and property history examples so every referenced field is defined by the sample blueprints/entities, enabling `go run examples/...` without schema violations.
- Cleaned up README guidance around version pinning to point at `v0.2.1`.

### Fixed
- Removed mirror property creation from `examples/blueprints/create`, preventing Port API `422 Unprocessable Entity` errors during repeated runs.
- Ensured entity history/lookups use records created by the standard samples, so all examples can be executed in a fresh workspace without manual data tweaks.
