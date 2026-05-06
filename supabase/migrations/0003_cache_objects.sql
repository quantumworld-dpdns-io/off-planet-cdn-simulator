-- Migration 0003: Cache objects, versions, and tags

-- ============================================================
-- ENUM types
-- ============================================================

CREATE TYPE cache_object_status AS ENUM ('ACTIVE', 'DEPRECATED', 'DELETED');

-- ============================================================
-- cache_objects
-- ============================================================

CREATE TABLE cache_objects (
  id                 uuid                PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id             uuid                NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id            uuid                REFERENCES sites(id) ON DELETE SET NULL,
  priority_class_id  uuid                REFERENCES priority_classes(id) ON DELETE SET NULL,
  name               text                NOT NULL,
  content_type       text,
  source_url         text,
  content_hash       text,
  size_bytes         bigint              NOT NULL DEFAULT 0,
  pinned             boolean             NOT NULL DEFAULT false,
  status             cache_object_status NOT NULL DEFAULT 'ACTIVE',
  tags               text[]              NOT NULL DEFAULT '{}',
  metadata           jsonb               NOT NULL DEFAULT '{}',
  current_version_id uuid,
  created_at         timestamptz         NOT NULL DEFAULT now(),
  updated_at         timestamptz         NOT NULL DEFAULT now()
);

-- ============================================================
-- cache_object_versions
-- ============================================================

CREATE TABLE cache_object_versions (
  id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  object_id    uuid        NOT NULL REFERENCES cache_objects(id) ON DELETE CASCADE,
  version      int         NOT NULL,
  content_hash text        NOT NULL,
  size_bytes   bigint      NOT NULL,
  storage_path text,
  created_at   timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- cache_object_tags
-- ============================================================

CREATE TABLE cache_object_tags (
  id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id     uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  object_id  uuid        NOT NULL REFERENCES cache_objects(id) ON DELETE CASCADE,
  tag        text        NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_cache_object_tags_object_tag UNIQUE (object_id, tag)
);

-- ============================================================
-- Self-referential FK: cache_objects.current_version_id -> cache_object_versions
-- Added after both tables exist
-- ============================================================

ALTER TABLE cache_objects
  ADD CONSTRAINT fk_cache_objects_current_version
  FOREIGN KEY (current_version_id) REFERENCES cache_object_versions(id) ON DELETE SET NULL;

-- ============================================================
-- updated_at trigger
-- ============================================================

CREATE TRIGGER trg_cache_objects_updated_at
  BEFORE UPDATE ON cache_objects
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- RLS
-- ============================================================

ALTER TABLE cache_objects         ENABLE ROW LEVEL SECURITY;
ALTER TABLE cache_object_versions ENABLE ROW LEVEL SECURITY;
ALTER TABLE cache_object_tags     ENABLE ROW LEVEL SECURITY;

CREATE POLICY "cache_objects_org_isolation" ON cache_objects
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "cache_object_versions_org_isolation" ON cache_object_versions
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

CREATE POLICY "cache_object_tags_org_isolation" ON cache_object_tags
  USING (org_id = (auth.jwt() ->> 'org_id')::uuid);

-- ============================================================
-- Indexes
-- ============================================================

CREATE INDEX idx_cache_objects_site_priority
  ON cache_objects(site_id, priority_class_id);

CREATE INDEX idx_cache_objects_tags
  ON cache_objects USING gin(tags);

CREATE INDEX idx_cache_object_versions_object
  ON cache_object_versions(object_id, version DESC);
