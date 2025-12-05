# Examples

Small, runnable entrypoints that exercise each service in the SDK. Every example loads configuration from `.env` using `config.Load(".env")`, so create a file at the repo root (or export the variables) with the usual credentials.

## Prerequisites

- Go 1.22 or later
- A Port account with API credentials
- A `.env` file in the repository root (or environment variables set)

## Configuration

Create a `.env` file in the repository root with your credentials:

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

## Common Patterns

### Error Handling

All examples demonstrate proper error handling. The SDK uses typed errors (`porter.Error`) that include HTTP status codes:

```go
entity, err := cli.Entities().Get(ctx, "blueprint", "identifier")
if err != nil {
    var perr *porter.Error
    if errors.As(err, &perr) {
        switch perr.StatusCode {
        case 404:
            log.Println("Entity not found")
        case 401:
            log.Println("Authentication failed")
        default:
            log.Printf("API error: %v", err)
        }
    }
    return err
}
```

### Context Usage

All API methods accept a `context.Context` for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use ctx in API calls
entities, err := cli.Entities().List(ctx, "blueprint", nil)
```

### Pagination

For endpoints that return paginated results, use the `ListAll` or `ListAllBlueprint` helpers:

```go
opts := entities.SearchOptions{
    Query: map[string]any{
        "composite": map[string]any{
            "operator": "and",
            "rules": []any{},
        },
    },
    Limit: 100,
}
allEntities, err := cli.Entities().ListAllBlueprint(ctx, "my-blueprint", opts)
```

Or handle pagination manually:

```go
opts := entities.SearchOptions{Limit: 100}
var allEntities []entities.Entity

for {
    resp, err := cli.Entities().SearchBlueprint(ctx, "blueprint", opts)
    if err != nil {
        return err
    }
    allEntities = append(allEntities, resp.Entities...)
    if !resp.HasMore() {
        break
    }
    opts.From = resp.Next
}
```

### Resource Cleanup

Always close the client when done to release resources (especially important when using verbose logging):

```go
cli, err := client.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer cli.Close()
```

## Catalog

- **auth/**
  - `access_token`: call `/v1/auth/access_token` and print the prefix of the returned token.
  - `rotate_credentials`: rotate API credentials for the address in `PORT_INVITE_EMAIL`.
- **entities/**
  - `list`, `get`, `create`, `update`, `delete`: CRUD operations.
  - `upsert`, `bulk_upsert`: idempotent single/batch writes.
  - `bulk_delete`: remove batches with optional cascade.
  - `link` / `unlink`: manage relations.
  - `search`, `aggregate`, `aggregate_over_time`, `properties_history`: query/insight helpers.
- **organization/**
  - `get`: inspect org metadata/banner settings.
  - `patch`: update portal settings (title/icon/announcement).
  - `secrets`: list secrets and optionally create/update/delete one when env vars provided.
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
