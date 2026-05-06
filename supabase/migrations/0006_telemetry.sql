-- Migration 0006: Mirror sources/artifacts, package/model registries, content requests, telemetry events

-- ============================================================
-- mirror_sources
-- ============================================================

CREATE TABLE mirror_sources (
  id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  registry_type text        NOT NULL, -- npm | pypi | crates_io | oci | model
  upstream_url  text        NOT NULL,
  label         text,
  enabled       boolean     NOT NULL DEFAULT true,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- mirror_artifacts
-- ============================================================

CREATE TABLE mirror_artifacts (
  id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  source_id    uuid        REFERENCES mirror_sources(id) ON DELETE CASCADE,
  name         text        NOT NULL,
  version      text        NOT NULL,
  size_bytes   bigint,
  storage_path text,
  synced_at    timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_mirror_artifacts_source_name_version UNIQUE (source_id, name, version)
);

-- ============================================================
-- package_registries
-- ============================================================

CREATE TABLE package_registries (
  id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id     uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id    uuid        REFERENCES sites(id) ON DELETE CASCADE,
  source_id  uuid        REFERENCES mirror_sources(id) ON DELETE SET NULL,
  local_port int,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- model_registries
-- ============================================================

CREATE TABLE model_registries (
  id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id     uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id    uuid        REFERENCES sites(id) ON DELETE CASCADE,
  source_id  uuid        REFERENCES mirror_sources(id) ON DELETE SET NULL,
  local_port int,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- content_requests
-- ============================================================

CREATE TABLE content_requests (
  id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  node_id      uuid        REFERENCES nodes(id) ON DELETE SET NULL,
  object_id    uuid        REFERENCES cache_objects(id) ON DELETE SET NULL,
  miss         boolean     NOT NULL DEFAULT false,
  bytes_served bigint      NOT NULL DEFAULT 0,
  duration_ms  int,
  created_at   timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- telemetry_events (range-partitioned by created_at)
-- ============================================================

CREATE TABLE telemetry_events (
  id         uuid        NOT NULL DEFAULT gen_random_uuid(),
  org_id     uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id    uuid        REFERENCES sites(id) ON DELETE SET NULL,
  node_id    uuid        REFERENCES nodes(id) ON DELETE SET NULL,
  event_type text        NOT NULL,
  payload    jsonb       NOT NULL DEFAULT '{}',
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Monthly partitions for 2026
CREATE TABLE telemetry_events_y2026m01 PARTITION OF telemetry_events FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE telemetry_events_y2026m02 PARTITION OF telemetry_events FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE telemetry_events_y2026m03 PARTITION OF telemetry_events FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE telemetry_events_y2026m04 PARTITION OF telemetry_events FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE telemetry_events_y2026m05 PARTITION OF telemetry_events FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE telemetry_events_y2026m06 PARTITION OF telemetry_events FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE telemetry_events_y2026m07 PARTITION OF telemetry_events FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE telemetry_events_y2026m08 PARTITION OF telemetry_events FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE telemetry_events_y2026m09 PARTITION OF telemetry_events FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE telemetry_events_y2026m10 PARTITION OF telemetry_events FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE telemetry_events_y2026m11 PARTITION OF telemetry_events FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE telemetry_events_y2026m12 PARTITION OF telemetry_events FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

-- ============================================================
-- updated_at triggers
-- ============================================================

CREATE TRIGGER trg_mirror_sources_updated_at
  BEFORE UPDATE ON mirror_sources
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- RLS
-- ============================================================

ALTER TABLE mirror_sources    ENABLE ROW LEVEL SECURITY;
ALTER TABLE mirror_artifacts  ENABLE ROW LEVEL SECURITY;
ALTER TABLE package_registries ENABLE ROW LEVEL SECURITY;
ALTER TABLE model_registries  ENABLE ROW LEVEL SECURITY;
ALTER TABLE content_requests  ENABLE ROW LEVEL SECURITY;
ALTER TABLE telemetry_events  ENABLE ROW LEVEL SECURITY;

CREATE POLICY "mirror_sources_org_isolation" ON mirror_sources
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "mirror_artifacts_org_isolation" ON mirror_artifacts
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "package_registries_org_isolation" ON package_registries
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "model_registries_org_isolation" ON model_registries
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "content_requests_org_isolation" ON content_requests
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- telemetry_events: SELECT only via client
CREATE POLICY "telemetry_events_org_isolation" ON telemetry_events
  FOR SELECT
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- ============================================================
-- Indexes
-- ============================================================

CREATE INDEX idx_telemetry_events_site_created
  ON telemetry_events(site_id, created_at DESC);
