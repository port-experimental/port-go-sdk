# Examples

Small, runnable entrypoints that exercise each service in the SDK. Every example loads configuration from `.env` using `config.Load(".env")`, so create a file at the repo root (or export the variables) with the usual credentials:

| Variable | Purpose |
| --- | --- |
| `PORT_CLIENT_ID` / `PORT_CLIENT_SECRET` | Client-credential flow used by most examples. |
| `PORT_ACCESS_TOKEN` | Optional personal token (when set, client credentials are ignored). |
| `PORT_REGION` | `eu` (default) or `us`. |
| `PORT_BASE_URL` | Override for self-hosted/staging environments. |
| `PORT_INVITE_EMAIL` | Required only by invite/credential examples (see below). |

> Go 1.22+ is recommended. Run everything from the repo root so relative module paths resolve correctly.

## Running an example

```bash
# With env vars already exported
go run ./examples/entities/list

# Or point at a .env
PORT_CLIENT_ID=... PORT_CLIENT_SECRET=... go run ./examples/entities/upsert
```

Several programs contain placeholder blueprint/entity/datasource IDs—replace them with IDs from your Port account before running.

## Catalog

- **auth/**
  - `access_token`: call `/v1/auth/access_token` and print the prefix of the returned token.
  - `rotate_credentials`: rotate API credentials for the address in `PORT_INVITE_EMAIL`.
- **entities/**
  - `list`, `get`, `create`, `update`, `delete`: CRUD operations.
  - `upsert`, `bulk_upsert`: idempotent single/batch writes.
  - `bulk_delete`: remove batches with optional cascade.
  - `link` / `unlink`: manage relations.
  - `search`: filtered results (blueprint scoped).
- **blueprints/**
  - `list`, `get`: enumerate definitions.
  - `create`, `upsert`, `update`, `delete`: mutate blueprint schemas.
  - `permissions`: read blueprint-level permission rules.
- **datasources/**
  - `list`, `get`, `create`, `delete`: manage webhook/data source definitions.
  - `rotate-secret`: rotates the shared secret for a data source.
  - `set-mapping`: upload JSON mapping content.
- **automations/**
  - `list`, `get`: inspect automation definitions.
  - `trigger`: invoke an automation.
  - `executions`: list execution history.
- **users/**
  - `list-users`, `list-teams`: enumerate users/teams.
  - `assign-role`: assign a given role to a user (fill in IDs).
  - `invite`: invite `PORT_INVITE_EMAIL` with the default `Member` role.

Each directory is its own Go module entry, so you can `go run ./examples/<group>/<example>` independently. Use them as copy/paste references while building your own integrations.
