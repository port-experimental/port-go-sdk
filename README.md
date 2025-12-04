# Port Go SDK

Ground-up Go client for the Port REST API (`https://api.port.io/swagger/json`). Targets both EU (`https://api.port.io`) and US (`https://api.us.port.io`) regions, with `.env`/environment variable configuration.

## Getting started

```bash
go get github.com/port-experimental/port-go-sdk

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

Examples (per method):
- Entities: `examples/entities/{list,get,upsert,update,delete,link,unlink}`
- Blueprints: `examples/blueprints/{list,get,create,upsert,delete}`
- Automations: `examples/automations/{list,get,executions,trigger}`
- Data sources: `examples/datasources/{list,get,create,delete,rotate-secret,set-mapping}`
- Users: `examples/users/{list-users,list-teams,assign-role}`

See `CHECKLIST.md` for remaining coverage work.
