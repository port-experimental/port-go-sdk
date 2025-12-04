# Port Go SDK Overview

## Installation

```bash
go get github.com/port-experimental/port-go-sdk
```

## Configuration

Environment variables (or `.env` file using `config.Load(".env")`):

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
- Override: set `PORT_BASE_URL` or `cfg.BaseURL`.

## Examples

See `examples/README.md` for runnable snippets covering entities, blueprints, data sources, automations, and users.
