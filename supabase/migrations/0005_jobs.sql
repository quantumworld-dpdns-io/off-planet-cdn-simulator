-- Migration 0005: Preload jobs, eviction runs, bandwidth/contact windows

-- ============================================================
-- ENUM types
-- ============================================================

CREATE TYPE job_status AS ENUM ('PENDING', 'RUNNING', 'DONE', 'FAILED', 'CANCELLED');

CREATE TYPE preload_job_item_status AS ENUM ('PENDING', 'FETCHING', 'DONE', 'FAILED');

CREATE TYPE eviction_run_status AS ENUM ('DRY_RUN', 'COMMITTED', 'PARTIAL', 'FAILED');

-- ============================================================
-- preload_jobs
-- ============================================================

CREATE TABLE preload_jobs (
  id                     uuid                  PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id                 uuid                  NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id                uuid                  REFERENCES sites(id) ON DELETE SET NULL,
  name                   text                  NOT NULL,
  status                 job_status            NOT NULL DEFAULT 'PENDING',
  priority_class_filter  priority_class_level[],
  bandwidth_budget_bytes bigint,
  started_at             timestamptz,
  completed_at           timestamptz,
  created_at             timestamptz           NOT NULL DEFAULT now(),
  updated_at             timestamptz           NOT NULL DEFAULT now()
);

-- ============================================================
-- preload_job_items
-- ============================================================

CREATE TABLE preload_job_items (
  id                uuid                    PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id            uuid                    NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  job_id            uuid                    NOT NULL REFERENCES preload_jobs(id) ON DELETE CASCADE,
  object_id         uuid                    REFERENCES cache_objects(id) ON DELETE SET NULL,
  status            preload_job_item_status NOT NULL DEFAULT 'PENDING',
  bytes_transferred bigint                  NOT NULL DEFAULT 0,
  error_message     text,
  created_at        timestamptz             NOT NULL DEFAULT now(),
  updated_at        timestamptz             NOT NULL DEFAULT now()
);

-- ============================================================
-- eviction_runs
-- ============================================================

CREATE TABLE eviction_runs (
  id                 uuid               PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id             uuid               NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  node_id            uuid               REFERENCES nodes(id) ON DELETE SET NULL,
  status             eviction_run_status NOT NULL DEFAULT 'DRY_RUN',
  target_freed_bytes bigint             NOT NULL,
  actual_freed_bytes bigint             NOT NULL DEFAULT 0,
  policy_id          uuid               REFERENCES cache_policies(id) ON DELETE SET NULL,
  created_at         timestamptz        NOT NULL DEFAULT now(),
  updated_at         timestamptz        NOT NULL DEFAULT now()
);

-- ============================================================
-- eviction_candidates
-- ============================================================

CREATE TABLE eviction_candidates (
  id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id     uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  run_id     uuid        NOT NULL REFERENCES eviction_runs(id) ON DELETE CASCADE,
  object_id  uuid        REFERENCES cache_objects(id) ON DELETE SET NULL,
  score      float8      NOT NULL,
  size_bytes bigint      NOT NULL,
  evicted    boolean     NOT NULL DEFAULT false,
  reason     text,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- bandwidth_windows
-- ============================================================

CREATE TABLE bandwidth_windows (
  id                uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id            uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id           uuid        REFERENCES sites(id) ON DELETE CASCADE,
  label             text,
  window_start      timestamptz NOT NULL,
  window_end        timestamptz NOT NULL,
  bandwidth_bps     bigint      NOT NULL,
  reliability_score float8      NOT NULL DEFAULT 1.0,
  created_at        timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- contact_windows
-- ============================================================

CREATE TABLE contact_windows (
  id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id           uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id          uuid        REFERENCES sites(id) ON DELETE CASCADE,
  label            text,
  rrule            text        NOT NULL,
  duration_seconds int         NOT NULL,
  next_window_at   timestamptz,
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- updated_at triggers
-- ============================================================

CREATE TRIGGER trg_preload_jobs_updated_at
  BEFORE UPDATE ON preload_jobs
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_preload_job_items_updated_at
  BEFORE UPDATE ON preload_job_items
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_eviction_runs_updated_at
  BEFORE UPDATE ON eviction_runs
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_contact_windows_updated_at
  BEFORE UPDATE ON contact_windows
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- RLS
-- ============================================================

ALTER TABLE preload_jobs        ENABLE ROW LEVEL SECURITY;
ALTER TABLE preload_job_items   ENABLE ROW LEVEL SECURITY;
ALTER TABLE eviction_runs       ENABLE ROW LEVEL SECURITY;
ALTER TABLE eviction_candidates ENABLE ROW LEVEL SECURITY;
ALTER TABLE bandwidth_windows   ENABLE ROW LEVEL SECURITY;
ALTER TABLE contact_windows     ENABLE ROW LEVEL SECURITY;

CREATE POLICY "preload_jobs_org_isolation" ON preload_jobs
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "preload_job_items_org_isolation" ON preload_job_items
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "eviction_runs_org_isolation" ON eviction_runs
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "eviction_candidates_org_isolation" ON eviction_candidates
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "bandwidth_windows_org_isolation" ON bandwidth_windows
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "contact_windows_org_isolation" ON contact_windows
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- ============================================================
-- Indexes
-- ============================================================

CREATE INDEX idx_preload_job_items_job_status
  ON preload_job_items(job_id, status);
