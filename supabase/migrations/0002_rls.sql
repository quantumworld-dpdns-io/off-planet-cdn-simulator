-- Migration 0002: Row-Level Security, updated_at trigger, and org-isolation policies

-- ============================================================
-- Trigger function: set_updated_at
-- ============================================================

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;

-- ============================================================
-- Attach updated_at trigger to all mutable tables
-- ============================================================

CREATE TRIGGER trg_orgs_updated_at
  BEFORE UPDATE ON orgs
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_profiles_updated_at
  BEFORE UPDATE ON profiles
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_sites_updated_at
  BEFORE UPDATE ON sites
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_nodes_updated_at
  BEFORE UPDATE ON nodes
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- Enable RLS on all tables defined so far
-- ============================================================

ALTER TABLE orgs              ENABLE ROW LEVEL SECURITY;
ALTER TABLE roles             ENABLE ROW LEVEL SECURITY;
ALTER TABLE profiles          ENABLE ROW LEVEL SECURITY;
ALTER TABLE priority_classes  ENABLE ROW LEVEL SECURITY;
ALTER TABLE sites             ENABLE ROW LEVEL SECURITY;
ALTER TABLE nodes             ENABLE ROW LEVEL SECURITY;
ALTER TABLE node_heartbeats   ENABLE ROW LEVEL SECURITY;

-- ============================================================
-- Org isolation policies
-- ============================================================

-- orgs: users can only see their own org
CREATE POLICY "orgs_org_isolation" ON orgs
  USING (id = (auth.jwt() ->> 'org_id')::uuid);

-- roles
CREATE POLICY "roles_org_isolation" ON roles
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- profiles: org isolation + self access
CREATE POLICY "profiles_org_isolation" ON profiles
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "profiles_self_access" ON profiles
  USING (id = auth.uid());

CREATE POLICY "profiles_self_update" ON profiles
  FOR UPDATE
  USING (id = auth.uid());

-- priority_classes
CREATE POLICY "priority_classes_org_isolation" ON priority_classes
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- sites
CREATE POLICY "sites_org_isolation" ON sites
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- nodes
CREATE POLICY "nodes_org_isolation" ON nodes
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- node_heartbeats: SELECT only
CREATE POLICY "node_heartbeats_org_isolation" ON node_heartbeats
  FOR SELECT
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);
