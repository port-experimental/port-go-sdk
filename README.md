# Port Go SDK

Ground-up Go client for the Port REST API (`https://api.port.io/swagger/json`). Targets both EU (`https://api.port.io`) and US (`https://api.us.port.io`) regions, with `.env`/environment variable configuration.

## Installation

```bash
go get github.com/port-experimental/port-go-sdk
```

## Configuration

Environment variables (or `.env` file via `config.Load(".env")`):

| Variable | Description |
| --- | --- |
| `PORT_REGION` | `eu` (default) or `us`. |
| `PORT_BASE_URL` | Optional override for self-hosted environments. |
| `PORT_ACCESS_TOKEN` | Personal API token. |
| `PORT_CLIENT_ID` / `PORT_CLIENT_SECRET` | Client credentials (used when `PORT_ACCESS_TOKEN` empty). |

Programmatic configuration:

```go
cfg, _ := config.Load(".env")
cfg.Region = "us"
cli, _ := client.New(cfg, client.WithUserAgent("my-app/1.0"))
```

## Region handling

- EU (default): `https://api.port.io`
- US: `https://api.us.port.io`
- Override the base URL via `PORT_BASE_URL` or `cfg.BaseURL`

## Getting started

```bash
PORT_CLIENT_ID=... PORT_CLIENT_SECRET=... go run ./examples/entities/upsert
```

Or load from `.env` via `config.Load(".env")`.

## Error Handling

The SDK uses typed errors (`porter.Error`) for API responses. Check status codes to handle different error scenarios:

```go
entity, err := cli.Entities().Get(ctx, "blueprint", "identifier")
if err != nil {
    var perr *porter.Error
    if errors.As(err, &perr) {
        switch perr.StatusCode {
        case 404:
            log.Println("Entity not found")
        case 401, 403:
            log.Println("Authentication/authorization failed")
        case 429:
            log.Println("Rate limited - SDK will retry automatically")
        default:
            log.Printf("API error %d: %s", perr.StatusCode, perr.Message)
        }
    }
    return err
}
```

## Versioning

This SDK follows [Semantic Versioning](https://semver.org/). Version tags are available for pinning specific versions:

```bash
go get github.com/port-experimental/port-go-sdk@v0.2.1
```

See [CHANGELOG.md](CHANGELOG.md) for detailed release notes.

## Packages

| Package | Description |
|---------|-------------|
| `pkg/config` | Load environment variables and `.env` files, region-aware base URL resolution |
| `pkg/httpx` | Shared HTTP client with retry logic and connection pooling |
| `pkg/auth` | Token sources supporting client credentials and personal API tokens |
| `pkg/client` | Base API client with service accessors for all Port API endpoints |
| `pkg/entities` | Entity management: list, get, create, upsert, update, delete, relations, search, aggregation |
| `pkg/blueprints` | Blueprint management: list, get, create, upsert, delete, permissions |
| `pkg/datasources` | Data source and webhook configuration management |
| `pkg/automations` | Automation management: list, get, trigger, execution history |
| `pkg/organization` | Organization metadata and secret management |
| `pkg/users` | User and team management, role assignment |
| `pkg/webhooks` | Webhook utilities with HMAC SHA256 signature support |
| `pkg/porter` | Error types and helper functions for error handling |

## Advanced Usage

### Custom Retry Configuration

```go
cli, _ := client.New(cfg, 
    client.WithRetryAttempts(5), // Custom retry count
    client.WithUserAgent("my-app/1.0"),
)
defer cli.Close()
```

### Pagination

The SDK provides helpers for automatic pagination:

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
// Automatically fetches all pages
allEntities, err := cli.Entities().ListAllBlueprint(ctx, "blueprint", opts)
```

### Context and Timeouts

All API methods accept `context.Context` for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

entities, err := cli.Entities().List(ctx, "blueprint", nil)
```

## Examples

See `examples/README.md` for runnable snippets covering entities, blueprints, data sources, automations, organization, and users. Highlights:
- Entities: `examples/entities/{list,get,create,upsert,update,delete,bulk_upsert,bulk_delete,link,unlink,search,aggregate,aggregate_over_time,properties_history}`
- Blueprints: `examples/blueprints/{list,get,create,upsert,delete}`
- Automations: `examples/automations/{list,get,executions,trigger}`
- Data sources: `examples/datasources/{list,get,create,delete,rotate-secret,set-mapping}`
- Organization: `examples/organization/{get,patch,secrets}`
- Users: `examples/users/{list-users,list-teams,assign-role,invite}` (`invite` reads `PORT_INVITE_EMAIL`)

See `CHECKLIST.md` for remaining coverage work.