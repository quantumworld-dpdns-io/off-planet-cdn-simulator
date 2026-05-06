# Integration Tests

Live integration tests for the Off-Planet CDN Simulator control API. These tests require:

- The control API running (`go run ./cmd/api` from `services/control-api/`)
- A Supabase local dev database (started via `supabase start` or `docker compose up`)

## Running

```sh
SUPABASE_DB_URL=postgresql://postgres:postgres@localhost:54322/postgres \
CONTROL_API_URL=http://localhost:8080 \
go test ./... -v -count=1
```

Environment variables default to the values shown above if not set.

## Test coverage

- `health_test.go` — `GET /v1/health` liveness check
- `sites_test.go` — list, create, get, and 404/400 error cases for sites
- `nodes_test.go` — register nodes, heartbeat updates `last_seen` and `node_heartbeats` table, node status endpoint, and error cases
