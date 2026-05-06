-- Migration 0004: Cache policies, rules, and decisions
-- Note: cache_policies.wasm_plugin_id FK is deferred to 0007_audit_logs.sql
--       after wasm_plugins table is created.

-- ============================================================
-- ENUM types
-- ============================================================

CREATE TYPE cache_policy_action AS ENUM ('PIN', 'PREFETCH', 'EVICT_FIRST', 'IGNORE');

-- ============================================================
-- cache_policies
-- ============================================================

CREATE TABLE cache_policies (
  id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id         uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id        uuid        REFERENCES sites(id) ON DELETE SET NULL,
  name           text        NOT NULL,
  description    text,
  enabled        boolean     NOT NULL DEFAULT true,
  -- wasm_plugin_id FK added in migration 0007 after wasm_plugins exists
  wasm_plugin_id uuid,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- cache_policy_rules
-- ============================================================

CREATE TABLE cache_policy_rules (
  id                uuid               PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id            uuid               NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  policy_id         uuid               NOT NULL REFERENCES cache_policies(id) ON DELETE CASCADE,
  priority_class_id uuid               REFERENCES priority_classes(id) ON DELETE SET NULL,
  action            cache_policy_action NOT NULL,
  params            jsonb              NOT NULL DEFAULT '{}',
  sort_order        int                NOT NULL DEFAULT 0,
  created_at        timestamptz        NOT NULL DEFAULT now()
);

-- ============================================================
-- cache_decisions
-- ============================================================

CREATE TABLE cache_decisions (
  id        uuid               PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id    uuid               NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  node_id   uuid               REFERENCES nodes(id) ON DELETE SET NULL,
  object_id uuid               REFERENCES cache_objects(id) ON DELETE SET NULL,
  policy_id uuid               REFERENCES cache_policies(id) ON DELETE SET NULL,
  decision  cache_policy_action NOT NULL,
  score     float8,
  reason    text,
  created_at timestamptz       NOT NULL DEFAULT now()
);

-- ============================================================
-- updated_at trigger
-- ============================================================

CREATE TRIGGER trg_cache_policies_updated_at
  BEFORE UPDATE ON cache_policies
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- RLS
-- ============================================================

ALTER TABLE cache_policies      ENABLE ROW LEVEL SECURITY;
ALTER TABLE cache_policy_rules  ENABLE ROW LEVEL SECURITY;
ALTER TABLE cache_decisions     ENABLE ROW LEVEL SECURITY;

CREATE POLICY "cache_policies_org_isolation" ON cache_policies
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "cache_policy_rules_org_isolation" ON cache_policy_rules
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "cache_decisions_org_isolation" ON cache_decisions
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- ============================================================
-- Indexes
-- ============================================================

CREATE INDEX idx_cache_decisions_node_object
  ON cache_decisions(node_id, object_id, created_at DESC);
