# Changelog

All notable changes to this project will be documented in this file.

## v0.3.1 - 2026-02-03

### Added
- Introduced update/delete scorecard examples that reuse a shared bootstrap helper so each sample can run independently.
- Added `examples/scorecards/internal/setup` to provision the sandbox blueprint and sample scorecard before every example.

### Changed
- Scorecard list/get/create examples now seed consistent data before running, preventing 422 errors in fresh environments.

## v0.3.0 - 2026-02-03

### Added
- Implemented `pkg/scorecards` with list/get/create/update/delete helpers wired into `client.Client`.
- Added runnable scorecard examples (`examples/scorecards/{list,get,create}`) and documented the new capability in the README/checklist.

### Changed
- Regenerated `docs/port-openapi.json` from the latest Port swagger (`https://api.port.io/swagger/json`) to ensure the SDK reflects the current API surface and OpenAPI 3.1.0 schema validation.


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
