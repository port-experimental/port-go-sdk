# Port Go SDK Build Checklist

Top-level goal: implement a ground-up Go SDK (`github.com/port-experimental/port-go-sdk`) covering the entire Port REST API (https://api.port.io/swagger/json), with EU (`https://api.port.io`) default, US (`https://api.us.port.io`) optional, and documentation/examples for every endpoint.

## Environment / Setup
- [x] Confirm Go version (>=1.22), `go env`, and ensure `GOBIN`/`GOMODCACHE` paths writable.
- [ ] List/install optional tools (`golangci-lint`, `staticcheck`, `goimports`) if desired.
- [x] Create base folders: `cmd/`, `pkg/`, `examples/`, `docs/`, `.github/` (or CI equivalent), `.vscode/`.
- [x] Run `go mod init github.com/port-experimental/port-go-sdk`.
- [x] Configure `.gitignore`, `LICENSE`, `README.md` placeholders.
- [x] Add env loading helpers (`.env` support via `godotenv` or custom loader) and document expected variables.

## Core Infrastructure
- [x] Region/base URL resolver:
  - [x] Define config struct (region enum, base override).
  - [x] Default to `https://api.port.io`.
  - [x] Switch to `https://api.us.port.io` when region=US.
  - [x] Allow explicit base override for future regions/self-hosted.
- [x] HTTP client:
  - [x] Shared transport (timeouts, idle conns, TLS settings).
  - [x] User-agent helper (`port-go-sdk/<version>`).
  - [x] Retry middleware with exponential backoff + jitter.
  - [x] Request cloning so POST/PUT bodies replay on retry.
  - [x] Tests for retry (429, 5xx, Retry-After header, context cancel).
- [x] Authentication package:
  - [x] Client-credentials token exchange (`/v1/auth/access_token`).
  - [x] Personal API token passthrough.
  - [x] Automatic token refresh with buffer window.
  - [x] Tests for success/failure, region switching.
  - [x] Provide config via environment variables, `.env`, and functional options; validate required fields with clear errors.
- [x] Error handling:
  - [x] Define SDK error type (status code, message, body).
  - [x] Map Port error payloads into typed errors.
- [x] Webhook utility:
  - [x] HMAC SHA256 signature helper.
  - [x] Convenience POST method with signature + retry.
  - [x] Tests for signature formatting and failure paths.

## API Coverage (per Swagger)
- [ ] Schema modeling:
  - [ ] Generate or handcraft request/response structs for each endpoint group.
  - [ ] Document field tags and optional/required properties.
- [ ] Entities APIs:
  - [x] List entities.
  - [x] Get entity by identifier.
  - [x] Create/Upsert entity (with merge options).
  - [x] Update/patch entity.
  - [x] Delete entity.
  - [x] Manage relations (link/unlink).
- [ ] Blueprints APIs:
  - [x] List blueprints.
  - [x] Get blueprint.
  - [x] Create/update blueprint.
  - [x] Delete blueprint.
- [x] Data Sources / Webhooks:
  - [x] CRUD for data sources.
  - [x] Webhook secret rotation endpoints.
  - [x] Mapping upload/download.
- [ ] Automations / Triggers:
  - [x] List automations and executions.
  - [x] Trigger execution.
  - [x] Manage schedules (if exposed).
- [x] Access / Users / Teams:
  - [x] List users/teams/roles.
  - [x] Assign roles or permissions if API supports it.
- [ ] Misc endpoints from Swagger (dashboards, scorecards, etc. if present).
- [ ] Ensure every endpoint has a corresponding method and typed response.

## Examples & Documentation
- [x] Create `docs/overview.md` (installation, auth, region config).
- [x] Document configuration patterns (env vars, functional options).
- [x] `examples/README.md` describing layout.
- [ ] For each endpoint:
  - [ ] Add runnable example under `examples/<group>/<endpoint>.go`.
  - [ ] Include comments describing prerequisites (tokens, IDs).
  - [ ] Ensure examples compile via `go test ./examples/...` or `go vet`.
- [x] Update top-level README with:
  - [x] Quick start sample.
  - [x] Region selection snippet.
  - [x] Auth snippet (client credentials + personal token).
  - [x] Link matrix to docs/examples per API group.

## Tooling / Validation
- [x] Add `scripts/test.sh` to run `go test ./...` and build/ci-check examples.
- [ ] Optional `makefile` targets (`make test`, `make examples`, `make lint`).
- [ ] Integrate `go vet` and (optional) `golangci-lint`.
- [ ] Configure lint rules to enforce Go best practices (ineffassign, staticcheck, govet, gofmt) and document how to run them.
- [ ] When credentials available, smoke-test:
  - [ ] Auth + entity upsert on EU.
  - [ ] Auth + entity get on US.
- [ ] Prepare release notes and tag initial version (`v0.1.0`) after coverage and docs are complete.

Update this checklist as tasks progress (mark items `[x]` when done).
