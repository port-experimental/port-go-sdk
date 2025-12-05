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

## Packages

- `pkg/config` — load env/.env, region-aware base URL.
- `pkg/httpx` — shared HTTP client + retry logic.
- `pkg/auth` — token sources (client credentials or personal token).
- `pkg/client` — base API client + service accessors.
- `pkg/entities` — entity helpers (list/upsert/get/update/delete/relations).
- `pkg/blueprints` — blueprint list/get/upsert/delete.
- `pkg/datasources`, `pkg/users` — early scaffolds for additional API groups.
- `pkg/automations` — list/get/trigger/list-executions helpers.
- `pkg/webhooks` — signed webhook sender.

## Examples

See `examples/README.md` for runnable snippets covering entities, blueprints, data sources, automations, and users. Highlights:
- Entities: `examples/entities/{list,get,create,upsert,update,delete,bulk_upsert,bulk_delete,link,unlink,search,aggregate,aggregate_over_time,properties_history}`
- Blueprints: `examples/blueprints/{list,get,create,upsert,delete}`
- Automations: `examples/automations/{list,get,executions,trigger}`
- Data sources: `examples/datasources/{list,get,create,delete,rotate-secret,set-mapping}`
- Users: `examples/users/{list-users,list-teams,assign-role,invite}` (`invite` reads `PORT_INVITE_EMAIL`)

See `CHECKLIST.md` for remaining coverage work.
