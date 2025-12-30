# bookmark-management

Gin-based HTTP service (Go) providing:

- `GET /health-check` for service health/status
- `GET /gen-pass` for generating a random password

## Requirements

- Go `1.24.4` (per `go.mod`)

## Quick start

### 1) Configure

The service reads configuration from environment variables (see `internal/api/config.go`). Defaults are provided.

Supported variables:

- `APP_PORT` (default: `8080`)
- `SERVICE_NAME` (default: `bookmark-management`)
- `INSTANCE_ID` (default: auto-generated UUID if empty)

Note: the application does not automatically load `.env` (there is no dotenv loader in the code). If you want to use it, you must export these variables in your shell/session before running.

### 2) Run the API

```bash
go run ./cmd/api
```

The server listens on `:${APP_PORT}`.

## Endpoints

### Health check

`GET /health-check`

Response:

```json
{
  "message": "OK",
  "serviceName": "bookmark-management",
  "instanceID": "<string>"
}
```

Example:

```bash
curl -s http://localhost:8080/health-check
```

### Generate password

`GET /gen-pass`

Returns a random password string of length `10` (characters: `A-Z`, `a-z`, `0-9`).

Example:

```bash
curl -s http://localhost:8080/gen-pass
```

## Testing

Run all tests:

```bash
go test ./...
```

The repo includes:

- handler-level unit tests under `internal/handler/*_test.go`
- endpoint tests under `internal/test/endpoint/*_test.go` using `httptest` and the API engine

## Project structure

- `cmd/api` - application entrypoint (`main.go`)
- `internal/api` - Gin engine setup, endpoint registration, config loading
- `internal/handler` - HTTP handlers
- `internal/service` - business logic (health check, password generation)
- `internal/test/endpoint` - black-box style endpoint tests

## Notes

- There is currently no persistence layer wired in (folders `internal/model` and `internal/repository` exist but are empty).