-- Migration 0007: Incidents, audit_logs, wasm_plugins, agent_tool_calls
-- Also finalizes the deferred FK from cache_policies -> wasm_plugins

-- ============================================================
-- ENUM types
-- ============================================================

CREATE TYPE incident_severity AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL');

-- ============================================================
-- incidents
-- ============================================================

CREATE TABLE incidents (
  id          uuid              PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id      uuid              NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id     uuid              REFERENCES sites(id) ON DELETE SET NULL,
  node_id     uuid              REFERENCES nodes(id) ON DELETE SET NULL,
  title       text              NOT NULL,
  description text,
  severity    incident_severity NOT NULL DEFAULT 'LOW',
  resolved_at timestamptz,
  created_at  timestamptz       NOT NULL DEFAULT now(),
  updated_at  timestamptz       NOT NULL DEFAULT now()
);

-- ============================================================
-- audit_logs
-- ============================================================

CREATE TABLE audit_logs (
  id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  actor_id      uuid        REFERENCES profiles(id) ON DELETE SET NULL,
  action        text        NOT NULL,
  resource_type text        NOT NULL,
  resource_id   uuid,
  old_value     jsonb,
  new_value     jsonb,
  ip_address    inet,
  created_at    timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- wasm_plugins
-- ============================================================

CREATE TABLE wasm_plugins (
  id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  name         text        NOT NULL,
  description  text,
  wasm_hash    text        NOT NULL,
  storage_path text        NOT NULL,
  enabled      boolean     NOT NULL DEFAULT false,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- agent_tool_calls
-- ============================================================

CREATE TABLE agent_tool_calls (
  id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  tool_name     text        NOT NULL,
  input         jsonb       NOT NULL DEFAULT '{}',
  output        jsonb,
  duration_ms   int,
  error_message text,
  created_at    timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- Deferred FK: cache_policies.wasm_plugin_id -> wasm_plugins
-- ============================================================

ALTER TABLE cache_policies
  ADD CONSTRAINT fk_cache_policies_wasm_plugin
  FOREIGN KEY (wasm_plugin_id) REFERENCES wasm_plugins(id) ON DELETE SET NULL;

-- ============================================================
-- updated_at triggers
-- ============================================================

CREATE TRIGGER trg_incidents_updated_at
  BEFORE UPDATE ON incidents
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_wasm_plugins_updated_at
  BEFORE UPDATE ON wasm_plugins
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- RLS
-- ============================================================

ALTER TABLE incidents       ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs      ENABLE ROW LEVEL SECURITY;
ALTER TABLE wasm_plugins    ENABLE ROW LEVEL SECURITY;
ALTER TABLE agent_tool_calls ENABLE ROW LEVEL SECURITY;

CREATE POLICY "incidents_org_isolation" ON incidents
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- audit_logs: SELECT only via client
CREATE POLICY "audit_logs_org_isolation" ON audit_logs
  FOR SELECT
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "wasm_plugins_org_isolation" ON wasm_plugins
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "agent_tool_calls_org_isolation" ON agent_tool_calls
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- ============================================================
-- Indexes
-- ============================================================

CREATE INDEX idx_audit_logs_actor_created
  ON audit_logs(actor_id, created_at DESC);

CREATE INDEX idx_audit_logs_resource
  ON audit_logs(resource_type, resource_id);
