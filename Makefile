.PHONY: dev-up dev-down db-migrate seed build-edge \
        test-go test-rust test-web test-all \
        lint-go lint-rust lint-web

# ── Dev environment ──────────────────────────────────────────────────────────

dev-up:
	infra/scripts/dev-up.sh

dev-down:
	infra/scripts/dev-down.sh

# ── Database ─────────────────────────────────────────────────────────────────

db-migrate:
	supabase db push

seed:
	supabase db reset --linked

# ── Build ─────────────────────────────────────────────────────────────────────

build-edge:
	cd edge && cargo build --release

# ── Tests ─────────────────────────────────────────────────────────────────────

test-go:
	go test ./services/... -count=1 -race

test-rust:
	cd edge && cargo test

test-web:
	pnpm test

test-all: test-go test-rust test-web

# ── Lint ──────────────────────────────────────────────────────────────────────

lint-go:
	golangci-lint run ./services/...

lint-rust:
	cd edge && cargo clippy -- -D warnings

lint-web:
	pnpm lint
