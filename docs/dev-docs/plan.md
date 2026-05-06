# Off-Planet CDN Simulator — Implementation & Testing Plan

> Audience: engineering leads and individual contributors  
> Status: living document — update as phases complete  
> Last revised: 2026-05-07

---

## 0. Project overview

Off-Planet CDN is a priority-aware content-distribution simulator for Moon/Mars habitats and other extreme-latency, bandwidth-scarce environments. The clean MVP is:

```
Next.js admin console
  + Go control API
  + Supabase (Postgres + Auth + Storage)
  + Redis / Dragonfly
  + Rust edge-agent / cache-proxy
  + Rust eviction engine
```

The full system adds semantic search (Qdrant), local analytics (DuckDB/Arrow), a WASM policy-plugin sandbox (Wasmtime), an MCP server for agentic operations, and longer-term security layers (Tetragon, PQC, RISC Zero, Teaclave).

---

## 1. Monorepo layout & toolchain bootstrap

### 1.1 Workspace setup

| Tool | Purpose |
|---|---|
| `pnpm` + `turbo.json` | JS/TS monorepo build orchestration |
| `go.work` | Multi-module Go workspace |
| `cargo` workspace | Rust multi-crate workspace under `edge/` |
| `supabase CLI` | Local Postgres + Auth + Storage + migrations |
| `docker compose` | All infra dependencies (Redis, Qdrant, OTel collector) |
| `make` | Top-level developer shortcuts |

**Deliverables**

- `pnpm-workspace.yaml` listing `apps/*`, `packages/*`
- `turbo.json` with `build`, `dev`, `test`, `lint` pipelines
- `go.work` referencing `services/control-api`, `services/scheduler`, `services/telemetry-ingest`, `services/policy-engine`, `services/mcp-server`, `packages/sdk-go`
- `edge/Cargo.toml` workspace referencing all crates
- `.env.example` with all required env-var keys
- `docker-compose.yml` with Redis/Dragonfly, Qdrant, OTel Collector, and optional Ollama services
- `Makefile` targets: `dev-up`, `dev-down`, `db-migrate`, `seed`, `build-edge`, `test-all`
- `infra/scripts/dev-up.sh` / `dev-down.sh` / `seed-demo.sh` / `build-edge.sh`

### 1.2 CI skeleton

Create `.github/workflows/`:

| File | Triggers | Jobs |
|---|---|---|
| `ci.yml` | PR + main push | lint, type-check, unit tests (all layers) |
| `go.yml` | changes to `services/**` | `go vet`, `golangci-lint`, `go test ./...` |
| `rust.yml` | changes to `edge/**` | `cargo clippy`, `cargo fmt --check`, `cargo test` |
| `web.yml` | changes to `apps/web/**` | `pnpm type-check`, `pnpm lint`, `pnpm test` |
| `docker.yml` | main push | build + push images to GHCR |
| `security.yml` | weekly + main push | `govulncheck`, `cargo audit`, `pnpm audit` |

---

## 2. Phase 0 — Repo bootstrap

**Goal:** runnable skeletons with health endpoints; nothing functional yet.

### 2.1 Supabase project

1. `supabase init` → `supabase/config.toml`
2. Write migrations in order (one PR per migration file):
   - `0001_init.sql` — `orgs`, `profiles`, `roles`
   - `0002_rls.sql` — enable RLS; add `org_id` policies for every table
   - `0003_cache_objects.sql` — `sites`, `nodes`, `node_heartbeats`, `cache_objects`, `cache_object_versions`, `cache_object_tags`, `priority_classes`
   - `0004_policies.sql` — `cache_policies`, `cache_policy_rules`, `cache_decisions`
   - `0005_jobs.sql` — `preload_jobs`, `preload_job_items`, `eviction_runs`, `eviction_candidates`, `bandwidth_windows`, `contact_windows`
   - `0006_telemetry.sql` — `mirror_sources`, `mirror_artifacts`, `package_registries`, `model_registries`, `content_requests`, `telemetry_events`
   - `0007_audit_logs.sql` — `incidents`, `audit_logs`, `wasm_plugins`, `agent_tool_calls`
3. `supabase/seed.sql` — demo org, two sites (Lunar Hab Alpha, Mars Base Camp), one admin user, sample priority classes, sample policy
4. Supabase Edge Functions stubs: `webhook-node-alert`, `export-cache-report`

### 2.2 Go control API skeleton

- `services/control-api/cmd/api/main.go` — starts HTTP server (Gin or Chi), wires OTel, connects to Supabase Postgres via `pgx`
- `GET /v1/health` returns `{"status":"ok"}`
- `internal/config/` reads env vars with validation
- `internal/db/` exports a typed Postgres client
- `internal/middleware/` — request logger, OTel trace propagation, auth JWT verification (Supabase JWTs)

### 2.3 Rust edge-agent skeleton

- `edge/crates/edge-agent/src/main.rs` — starts Axum server
- `GET /local/health` returns `{"status":"ok"}`
- `edge/crates/shared/` — `types.rs` (priority enum, cache-score struct), `errors.rs`, `telemetry.rs` (OTel setup)

### 2.4 Next.js app shell

- `apps/web/` with Next.js 14 App Router, TypeScript, Tailwind CSS
- `middleware.ts` — Supabase Auth session check, redirect unauthenticated users to `/login`
- `/login/page.tsx` — Supabase Auth email+password form
- `/admin/layout.tsx` — sidebar nav linking all admin pages
- `/admin/dashboard/page.tsx` — placeholder "Coming soon"
- `/app/layout.tsx` — user portal layout
- `/app/home/page.tsx` — placeholder

---

## 3. Phase 1 — Functional cache control

**Goal:** full data lifecycle: catalog → policy → preload job → edge sync → cache status visible in admin dashboard.

### 3.1 Control API — full implementation

Implement all routes from the API design:

#### Sites & nodes
- `GET /v1/sites`, `POST /v1/sites`, `GET /v1/sites/:site_id`
- `GET /v1/nodes`, `POST /v1/nodes/register`
- `POST /v1/nodes/:node_id/heartbeat` — upserts `node_heartbeats`, updates `nodes.last_seen`
- `GET /v1/nodes/:node_id/status`

#### Cache objects
- `GET /v1/cache/objects` — filterable by `priority_class`, `site_id`, `tag`
- `POST /v1/cache/objects` — validates schema against `shared-schemas/cache-object.schema.json`
- `GET /v1/cache/objects/:object_id`
- `POST /v1/cache/objects/:object_id/pin` — sets `pinned=true`, records to `audit_logs`
- `POST /v1/cache/objects/:object_id/unpin`

#### Policies
- `GET /v1/policies`, `POST /v1/policies`, `PUT /v1/policies/:policy_id`

#### Preload jobs
- `POST /v1/preload/jobs` — creates `preload_jobs` + `preload_job_items`, enqueues to Redis
- `GET /v1/preload/jobs`, `GET /v1/preload/jobs/:job_id`
- `POST /v1/preload/jobs/:job_id/cancel`

#### Telemetry & audit
- `POST /v1/telemetry/events` — validates and writes to `telemetry_events`
- `GET /v1/audit-logs` — paginated, filterable by `actor`, `action`, `resource_type`

### 3.2 Go scheduler

- `cmd/scheduler/main.go` — polls Redis preload queue, dispatches fetch tasks to edge nodes
- `internal/contactwindows/` — reads `contact_windows` from Postgres, blocks dispatch outside windows
- `internal/queues/` — Redis client with typed job structs
- `internal/optimizer/` — basic FIFO ordering by priority class (full optimizer in Phase 2)

### 3.3 Rust edge-agent — real implementation

- `src/sync.rs` — polls control API for pending preload jobs, downloads content, writes to local disk
- `src/heartbeat.rs` — sends `POST /v1/nodes/:node_id/heartbeat` every 30 s
- `src/api.rs` — implements local API:
  - `GET /local/cache/status`
  - `POST /local/cache/fetch`
  - `POST /local/cache/preload`
  - `POST /local/policy/reload`

### 3.4 Rust cache-proxy

- `src/proxy.rs` — Axum HTTP proxy: checks local disk for cached content, serves it; on miss records `content_requests.miss=true` and (optionally) initiates upstream fetch
- `src/fetch.rs` — upstream HTTP fetch with timeout; writes to disk under a content-addressable path
- `src/range_requests.rs` — supports HTTP `Range` headers for large file streaming
- `src/headers.rs` — injects `X-Cache: HIT|MISS`, `X-Priority-Class`, `X-Offline-Available` response headers

### 3.5 Next.js admin console — Phase 1 pages

| Page | Key widgets |
|---|---|
| `/admin/dashboard` | Node health summary, cache fill %, recent preload jobs, recent audit events |
| `/admin/sites` | Table of sites; create-site modal |
| `/admin/sites/[siteId]` | Site detail: nodes list, cache summary |
| `/admin/nodes` | Node list with heartbeat age, status badge |
| `/admin/nodes/[nodeId]` | Node detail: cache usage, heartbeat history |
| `/admin/cache-policies` | Policy list; create / edit policy form |
| `/admin/cache-policies/[policyId]` | Policy rules, assigned sites, rule editor |
| `/admin/content-catalog` | Paginated object list, filter by priority / tag / site |
| `/admin/content-catalog/[objectId]` | Object detail: metadata, cache decisions, versions |
| `/admin/preload-jobs` | Job list with status badges |
| `/admin/preload-jobs/[jobId]` | Job detail: items, progress, cancel button |
| `/admin/audit-logs` | Paginated, filterable audit log table |

### 3.6 Next.js user portal — Phase 1 pages

| Page | Content |
|---|---|
| `/app/home` | Welcome, offline status banner, quick-search bar |
| `/app/manuals` | List of cached manuals (P0/P1), search |
| `/app/medical` | Filtered view of medical content |
| `/app/engineering` | Filtered view of engineering content |
| `/app/offline-status` | Cache fill %, last sync time, contact window schedule |
| `/app/downloads` | In-progress and completed downloads |

### 3.7 Shared schemas & SDK stubs

- `packages/shared-schemas/` — finalize JSON Schema files for `cache-object`, `cache-policy`, `node-heartbeat`, `preload-job`
- `packages/sdk-js/src/client.ts` — typed fetch wrapper for control API; auto-generated from OpenAPI spec
- `packages/sdk-go/offcdn/client.go` — same for Go

---

## 4. Phase 2 — Eviction engine, simulation, and analytics

### 4.1 Rust eviction engine

Implement the cache scoring formula in `edge/crates/eviction-engine/src/score.rs`:

```
cache_score =
  priority_weight
+ mission_relevance
+ predicted_demand
+ offline_criticality
+ revalidation_cost
+ fetch_latency_cost
+ package_dependency_score
- size_penalty
- staleness_penalty
- redundancy_penalty
```

Files:
- `src/score.rs` — `score_object(obj: &CacheObject, ctx: &ScoringContext) -> f64`
- `src/constraints.rs` — capacity limits, pinned-object guards (P0 never evicted)
- `src/simulator.rs` — dry-run eviction: given a target freed-bytes goal, return `Vec<EvictionCandidate>` without touching disk
- `src/tests.rs` — exhaustive unit tests (see §6.2)

### 4.2 Eviction API surface

Control API additions:
- `POST /v1/eviction/simulate` — calls policy-engine, returns ranked eviction candidates with scores
- `POST /v1/eviction/run` — executes eviction on a named edge node via local API

Edge local API additions:
- `POST /local/cache/evict` — receives eviction plan from control plane, deletes from disk, writes `eviction_runs` + `eviction_candidates` to Postgres via control API

### 4.3 Bandwidth window model

- Database: `bandwidth_windows` (site, start, end, bandwidth_bps, reliability_score) and `contact_windows` (site, comms-window schedule for Mars/Moon)
- Scheduler: read active window before dispatching; compute optimal batch size as `window_duration_s × bandwidth_bps × reliability_score`
- Admin page `/admin/bandwidth-windows`: timeline chart of upcoming windows, CRUD form

### 4.4 Preload optimizer

`services/scheduler/internal/optimizer/` — given pending preload job items and active bandwidth window, sort items by cache score descending; fill window greedily by size; emit sorted dispatch list to Redis queue.

### 4.5 DuckDB / Arrow local analytics

- `edge/crates/eviction-engine/` or standalone `edge/crates/content-indexer/` — append `cache_decisions` and `eviction_runs` rows to a local DuckDB file
- `GET /local/analytics/bandwidth` — query DuckDB; return bandwidth usage per hour as Arrow IPC
- Admin page `/admin/dashboard` chart: bandwidth usage over time, cache hit rate over time

### 4.6 Admin eviction simulator page

`/admin/eviction-simulator`:
- Input: site, target freed-bytes, dry-run toggle
- Calls `POST /v1/eviction/simulate`
- Displays ranked candidate table: object, priority class, current score, freed bytes, reason
- "Confirm eviction" button calls `POST /v1/eviction/run`

---

## 5. Phase 3 — Registry mirrors

### 5.1 Rust registry-mirror crate

`edge/crates/registry-mirror/src/`:

| Module | Registry | Protocol |
|---|---|---|
| `npm.rs` | npm | tarballs + `registry.npmjs.org` metadata |
| `pypi.rs` | PyPI | simple index + sdists/wheels |
| `crates_io.rs` | crates.io | sparse index + `.crate` tarballs |
| `oci.rs` | OCI (Docker Hub, GHCR) | OCI distribution spec |
| `model_registry.rs` | Hugging Face / Ollama | model file + manifest |
| `manifest.rs` | shared | manifest struct, signing, validation |

Control API additions:
- `GET /v1/mirrors/sources`, `POST /v1/mirrors/sources`
- `POST /v1/mirrors/sync` — enqueues a mirror sync job

Edge local API:
- `POST /local/mirror/sync` — pulls from upstream into local mirror directory

Admin pages: `/admin/package-mirrors`, `/admin/model-mirrors`

User pages: `/app/packages`, `/app/models`

---

## 6. Phase 4 — AI and agentic operations

### 6.1 MCP server (Go)

`services/mcp-server/internal/tools/`:

| Tool | Description |
|---|---|
| `cache_status.go` | Query current cache fill, pinned objects, top-N by score |
| `generate_preload_plan.go` | Given a bandwidth window and priority class filter, return prioritized preload plan |
| `inspect_node.go` | Return full node status, heartbeat age, local cache stats |
| `simulate_eviction.go` | Call eviction simulator; return candidate list |
| `summarize_incident.go` | Summarize an incident from `incidents` table with audit context |

Expose tools via MCP JSON-RPC and via OpenAPI Tool Calling (`services/mcp-server/openapi.yaml`).

### 6.2 Semantic search with Qdrant

- `edge/crates/content-indexer/src/`:
  - `chunk.rs` — split documents (PDF, Markdown, plain text) into ~512-token chunks
  - `embeddings.rs` — call Ollama `/api/embeddings` or a remote embedding endpoint
  - `qdrant.rs` — upsert chunks as Qdrant point vectors; query for nearest neighbours
  - `duckdb_export.rs` — export embedding metadata to DuckDB for analytics

- Control API: `GET /v1/search?q=&site_id=&priority_class=` — proxies Qdrant query, enriches with Postgres metadata

- User page `/app/search` — full semantic search across cached manuals, medical docs, engineering content

### 6.3 Admin AI copilot

- Ollama (`llama3` or similar) serves locally via `docker compose`
- System prompt includes current cache state from `cache_status` MCP tool
- Admin page `/admin/settings` — AI copilot toggle; surfaced as a floating chat panel
- Arize Phoenix or Weave for AI trace evaluation (optional, dev-only)

---

## 7. Phase 5 — Security and verifiability

### 7.1 WASM policy plugin sandbox

`edge/crates/wasm-policy-runner/src/`:
- `runtime.rs` — Wasmtime engine setup with WASI; fuel limit for plugin execution budget
- `host.rs` — host functions: `get_object_metadata`, `get_site_context`, `emit_score`
- `wit.rs` — WIT interface: `cache-policy.wit` defines `score-object: func(obj: cache-object) -> f64`

`wasm-plugins/examples/`:
- `medical_first/` — boosts P0/P1 by 100 points
- `entertainment_last/` — reduces P4 by 50 points
- `package_dependency_boost/` — raises score when object is a transitive dep of a pinned package

Control API: `POST /v1/policies` accepts optional `wasm_plugin_id`; policy engine loads the plugin when scoring.

### 7.2 Tetragon runtime monitoring

`infra/tetragon/`:
- `tracing-policy-cache-agent.yaml` — alerts on unexpected `exec`, `open`, or network calls from `edge-agent`
- `tracing-policy-sensitive-files.yaml` — alerts on access to private key material or sensitive manifests

### 7.3 Manifest signing

`edge/crates/shared/src/crypto.rs`:
- Sign cache manifests with Ed25519 before writing to Supabase Storage
- Verify signature on load in `manifests.rs`
- PQC-ready interface: stub `sign_pqc` / `verify_pqc` for future CRYSTALS-Dilithium replacement

### 7.4 Proof-of-cache prototype (Phase 3+ roadmap)

- Stub RISC Zero or Noir circuit: given a cache manifest hash and a set of eviction decisions, produce a succinct proof that the decisions followed the declared policy
- Not required for MVP; architecture doc in `docs/security-model.md`

---

## 8. Testing strategy

### 8.1 Testing layers overview

| Layer | Tools | Location |
|---|---|---|
| Rust unit | `cargo test`, `proptest` | `edge/crates/*/src/tests.rs` or inline |
| Go unit | `go test ./...`, `testify` | `services/*/internal/**/*_test.go` |
| TS/React unit | `vitest`, `@testing-library/react` | `apps/web/src/**/*.test.tsx` |
| Integration (Go ↔ Supabase) | `go test` with `testcontainers-go` | `tests/integration/` |
| Integration (Rust ↔ disk) | `cargo test` with `tempfile` | `edge/crates/*/tests/` |
| API contract | `hurl` or `bruno` | `tests/api/` |
| E2E | `playwright` | `tests/e2e/` |
| Load | `k6` | `tests/load/` |

### 8.2 Eviction engine — unit tests (critical path)

File: `edge/crates/eviction-engine/src/tests.rs`

**Score correctness tests**

| Test | Assertion |
|---|---|
| `p0_score_always_above_threshold` | P0 object score ≥ 10,000 regardless of size |
| `p4_score_below_p0` | P4 score < P0 score for equivalent metadata |
| `large_object_penalized` | 10 GB object scores lower than 1 MB object at same priority |
| `stale_object_penalized` | Object last accessed 90 days ago scores lower than same object accessed today |
| `pinned_object_immune_from_eviction` | `constraints::can_evict(pinned=true) == false` |
| `score_is_deterministic` | Same inputs produce same score across 100 runs |

**Simulator tests**

| Test | Assertion |
|---|---|
| `simulate_respects_capacity_target` | Sum of `freed_bytes` of candidates ≥ `target_freed_bytes` |
| `simulate_never_evicts_pinned` | No pinned or P0 object in simulator output |
| `simulate_evicts_p4_before_p3` | First candidate is always lower priority than last |
| `simulate_empty_cache` | Returns empty list, no panic |
| `simulate_insufficient_space` | Returns partial list, emits `InsufficientSpaceWarning` |

**Property-based tests** (proptest)

| Test | Strategy |
|---|---|
| `score_monotone_in_priority` | `∀ a,b: priority(a) < priority(b) ⟹ score(a) < score(b)` (all else equal) |
| `eviction_plan_valid_for_any_cache_state` | Random cache inventory → simulator → plan is valid |

### 8.3 Cache scoring formula — Go unit tests

File: `services/policy-engine/internal/scoring/*_test.go`

| Test | Assertion |
|---|---|
| `TestScoreWeightsSum` | All weight fields positive; no NaN/Inf |
| `TestMissionRelevanceOverride` | Site in "emergency mode" boosts P1 by ≥ 50 |
| `TestPredictedDemandDecay` | Demand score decays correctly over time |
| `TestPackageDependencyBoost` | Transitive dependency scores higher than isolated package |

### 8.4 Control API — integration tests

Use `testcontainers-go` to spin up Postgres + Redis; run real migrations before each test suite.

| Test | Assertions |
|---|---|
| `TestRegisterNode` | Node appears in Postgres; heartbeat row created |
| `TestHeartbeatUpdatesLastSeen` | `last_seen` timestamp advances on repeated heartbeats |
| `TestCreateCacheObjectAndPin` | Object persists; `pinned=true`; audit log row created |
| `TestCreatePreloadJobEnqueuedToRedis` | Job row in Postgres; job ID appears in Redis queue |
| `TestEvictionSimulate` | Returns candidate list sorted by score ascending |
| `TestAuditLogCreatedOnEviction` | After eviction, `audit_logs` has a row with `action=evict` |
| `TestRLSIsolation` | Org A cannot read Org B's sites (use two Supabase JWT tokens) |

### 8.5 Edge-agent integration tests

Use `tempfile` for disk, mock HTTP server (Axum `TestClient`) for control API.

| Test | Assertion |
|---|---|
| `test_heartbeat_sends_correctly` | Control API mock receives heartbeat within 35 s |
| `test_fetch_writes_to_disk` | After `POST /local/cache/fetch`, file exists at expected path |
| `test_preload_respects_priority_order` | P0 item fetched before P4 item in same job |
| `test_evict_removes_from_disk` | After `POST /local/cache/evict`, file no longer exists |
| `test_proxy_cache_hit` | Proxy responds with `X-Cache: HIT` for cached content |
| `test_proxy_cache_miss_fetches_upstream` | Proxy fetches upstream on miss; file written to disk |

### 8.6 Registry mirror tests

| Test | Assertion |
|---|---|
| `test_npm_manifest_parse` | Real `package.json` parses to `MirrorArtifact` correctly |
| `test_pypi_simple_index_parse` | Simple HTML index parsed to package list |
| `test_crates_io_index_parse` | Sparse index entry parsed to `CrateVersion` |
| `test_oci_manifest_pull` | OCI manifest + layer digests extracted correctly |
| `test_model_registry_manifest` | HuggingFace model card metadata parsed |

### 8.7 React component tests

Use `vitest` + `@testing-library/react` + MSW (mock service worker) for API mocking.

| Test | Assertion |
|---|---|
| `Dashboard renders node summary` | Shows correct node count from mock API |
| `PreloadJobList shows status badges` | RUNNING / COMPLETED / CANCELLED badges render |
| `CacheObjectDetail shows priority class` | P0 shown as red badge, P4 as grey |
| `EvictionSimulatorPage dry run works` | Submitting dry-run form shows candidate table |
| `UserPortal OfflineStatus shows last sync` | Banner shows correct timestamp from mock |
| `AuthRedirect sends unauthenticated to login` | Middleware redirects 401 → `/login` |

### 8.8 E2E tests (Playwright)

Targets local dev stack (`pnpm dev:web` + `go run ./cmd/api` + `cargo run -p edge-agent`).

| Scenario | Steps |
|---|---|
| Admin login → dashboard visible | Navigate to `/`, redirected to `/login`, sign in, land on `/admin/dashboard` |
| Create site → visible in list | Fill create-site form, submit, site appears in `/admin/sites` |
| Register node → heartbeat badge green | Register node via API, load `/admin/nodes`, badge shows ONLINE |
| Create content object → pin it | Create object, navigate to detail page, click Pin, badge shows PINNED |
| Submit preload job → appears as RUNNING | Submit job form, job list shows RUNNING |
| Eviction simulator dry run → candidate list shown | Fill eviction form, submit, table appears |

### 8.9 API contract tests

Use `hurl` (or Bruno) to define HTTP request/response pairs as plain text fixtures in `tests/api/`.

| Test file | Coverage |
|---|---|
| `health.hurl` | `GET /v1/health` → 200 |
| `sites.hurl` | Create, list, get site |
| `nodes.hurl` | Register, heartbeat, status |
| `cache-objects.hurl` | Create, list, pin, unpin |
| `policies.hurl` | Create, list, update |
| `preload-jobs.hurl` | Create, list, cancel |
| `eviction.hurl` | Simulate and run |
| `audit-logs.hurl` | Pagination, filter by action |

### 8.10 Load tests (k6)

File: `tests/load/cache-write-read.js`

Scenarios:
- **Heartbeat storm:** 500 nodes sending heartbeats concurrently; assert p95 < 200 ms
- **Preload job fanout:** 100 operators submitting jobs simultaneously; assert no queue overflow
- **Proxy hit rate:** 1000 RPS against cache-proxy with 90% hit rate; assert p99 < 50 ms
- **Eviction simulation at scale:** 10,000-object catalog; simulate runs in < 500 ms

---

## 9. Supabase schema detail

### Key design decisions

- Every table has `org_id uuid NOT NULL REFERENCES orgs(id)` → RLS policy: `org_id = auth.jwt() ->> 'org_id'`
- `cache_objects` versioned via `cache_object_versions`; latest version pointed to by `cache_objects.current_version_id`
- `priority_classes` is a lookup table (`P0`–`P5`) seeded in `seed.sql`; FK from `cache_objects.priority_class_id`
- `preload_job_items.status` enum: `PENDING | FETCHING | DONE | FAILED`
- `eviction_runs.status` enum: `DRY_RUN | COMMITTED | PARTIAL | FAILED`
- `node_heartbeats` partitioned by month (Postgres native partitioning) to avoid unbounded growth
- `telemetry_events` similarly partitioned; TTL policy managed by a scheduled Supabase Edge Function
- Supabase Realtime enabled on `nodes` and `preload_jobs` tables for dashboard live-updates

### Index recommendations

```sql
-- Most common query patterns
CREATE INDEX idx_cache_objects_site_priority ON cache_objects(site_id, priority_class_id);
CREATE INDEX idx_preload_job_items_job_status ON preload_job_items(job_id, status);
CREATE INDEX idx_audit_logs_actor_created ON audit_logs(actor_id, created_at DESC);
CREATE INDEX idx_node_heartbeats_node_created ON node_heartbeats(node_id, created_at DESC);
CREATE INDEX idx_telemetry_events_site_created ON telemetry_events(site_id, created_at DESC);
```

---

## 10. OpenAPI spec

`services/control-api/api/openapi.yaml` is the source of truth. SDK clients in `packages/sdk-js/` and `packages/sdk-go/` are generated from it using `oapi-codegen` (Go) and `openapi-typescript` (JS).

Spec validation runs in CI (`spectral lint openapi.yaml`) on every PR that touches `services/control-api/`.

---

## 11. Observability

### OTel instrumentation points

| Service | Spans |
|---|---|
| control-api | Every HTTP handler; Supabase Postgres query; Redis enqueue |
| scheduler | Contact-window check; optimizer run; dispatch loop iteration |
| edge-agent | Heartbeat; sync loop; fetch; eviction |
| cache-proxy | Proxy request; cache hit/miss; upstream fetch |
| eviction-engine | Scoring batch; simulator run |

Attributes on every span: `org_id`, `site_id`, `node_id`, `priority_class`.

### Key metrics (Prometheus via OTel exporter)

| Metric | Type | Description |
|---|---|---|
| `cdn_cache_fill_ratio` | Gauge | Local cache fill % per node |
| `cdn_cache_hit_total` | Counter | Cache proxy hits |
| `cdn_cache_miss_total` | Counter | Cache proxy misses |
| `cdn_eviction_objects_total` | Counter | Objects evicted per run |
| `cdn_preload_job_duration_seconds` | Histogram | End-to-end preload job time |
| `cdn_heartbeat_age_seconds` | Gauge | Age of last heartbeat per node |
| `cdn_scoring_duration_seconds` | Histogram | Time to score a batch of objects |

---

## 12. Infrastructure & deployment

### Docker Compose services (local dev)

| Service | Image | Port |
|---|---|---|
| `supabase` | managed via CLI | 54321 (API), 54322 (Studio) |
| `redis` | `redis:7-alpine` | 6379 |
| `qdrant` | `qdrant/qdrant` | 6333 |
| `otel-collector` | `otel/opentelemetry-collector` | 4317 (gRPC), 4318 (HTTP) |
| `ollama` | `ollama/ollama` (optional) | 11434 |

### Kubernetes manifests (`infra/k8s/`)

| Manifest | Deployed component |
|---|---|
| `namespace.yaml` | `offplanet` namespace |
| `control-api.yaml` | Deployment + Service + HPA |
| `scheduler.yaml` | Deployment + ConfigMap |
| `telemetry-ingest.yaml` | Deployment + Service |
| `redis.yaml` | StatefulSet + PVC |
| `qdrant.yaml` | StatefulSet + PVC |
| `otel-collector.yaml` | DaemonSet + ConfigMap |

### Terraform modules (`infra/terraform/`)

- `supabase/` — Supabase project, database password, auth config
- `cloudflare/` — DNS, R2 bucket for content manifests
- `object-storage/` — S3-compatible bucket policy, CORS

---

## 13. Implementation order (recommended sprint sequence)

| Sprint | Deliverable | Done when |
|---|---|---|
| S1 | Monorepo + CI + Supabase schema (Phase 0) | `make dev-up` passes; all 3 health endpoints return 200 |
| S2 | Go control API — sites, nodes, heartbeat | Integration tests pass; heartbeat visible in Postgres |
| S3 | Rust edge-agent — heartbeat + fetch | Node appears ONLINE in admin dashboard |
| S4 | Cache objects + policies + preload jobs (control API) | Full CRUD; job enqueued in Redis |
| S5 | Go scheduler — dispatch to edge | Edge agent downloads content from preload job |
| S6 | Rust cache-proxy — hit/miss + X-Cache headers | Playwright E2E: serve cached file |
| S7 | Next.js admin console — all Phase 1 pages | All pages render with real data |
| S8 | Next.js user portal — Phase 1 pages | User can browse cached manuals offline |
| S9 | Rust eviction engine — scoring + simulator | Unit tests pass; eviction simulator page functional |
| S10 | Bandwidth windows + preload optimizer | Scheduler respects contact windows |
| S11 | DuckDB analytics + dashboard charts | Hit-rate and bandwidth charts render |
| S12 | Registry mirrors (npm, PyPI, crates.io, OCI, model) | User page `/app/packages` lists mirrored packages |
| S13 | MCP server + semantic search (Qdrant) | `cache_status` tool returns real data; `/app/search` works |
| S14 | WASM policy runner + example plugins | `medical_first` plugin raises P0 scores correctly |
| S15 | Tetragon policies + manifest signing | CI `security.yml` passes `cargo audit`, signing round-trip test passes |

---

## 14. Non-goals for MVP (do not implement)

- Real satellite communication integration or DTN protocol
- Real-time video CDN optimization
- Autonomous irreversible cache deletion (simulation only)
- On-chain settlement or token economics
- Full confidential-computing deployment (Teaclave stubs only)
- PQC in production (interface stubs only)

---

## 15. Open questions & decisions needed

| # | Question | Owner | Target sprint |
|---|---|---|---|
| 1 | Redis vs. Dragonfly — do we need Dragonfly's multi-threading for job queue throughput at MVP scale? | infra lead | S1 |
| 2 | Qdrant vs. LanceDB — Qdrant is more mature; LanceDB is pure Rust. Given edge-agent is Rust, LanceDB may be preferable long-term. | edge lead | S12 |
| 3 | `go.work` vs. independent `go.mod` per service — `go.work` simplifies local dev but complicates independent versioning. | backend lead | S1 |
| 4 | Auth role model — separate `admin` / `operator` / `viewer` / `user` roles, or Supabase Custom Claims? | backend lead | S2 |
| 5 | Cache storage format — content-addressable (SHA256 path) vs. flat directory keyed by `object_id`. SHA256 deduplicates automatically. | edge lead | S3 |
| 6 | Edge agent deployment model for demo — Docker Compose sidecar vs. separate VM? | infra lead | S1 |
