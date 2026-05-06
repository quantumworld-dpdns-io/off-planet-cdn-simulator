-- Migration 0001: Core tables — orgs, roles, profiles, priority_classes, sites, nodes, node_heartbeats

-- ============================================================
-- ENUM types
-- ============================================================

CREATE TYPE node_status AS ENUM ('ONLINE', 'OFFLINE', 'DEGRADED', 'UNKNOWN');

CREATE TYPE priority_class_level AS ENUM ('P0', 'P1', 'P2', 'P3', 'P4', 'P5');

-- ============================================================
-- orgs
-- ============================================================

CREATE TABLE orgs (
  id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  name       text        NOT NULL,
  slug       text        NOT NULL UNIQUE,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- roles
-- ============================================================

CREATE TABLE roles (
  id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id     uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  name       text        NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- profiles
-- ============================================================

CREATE TABLE profiles (
  id           uuid        PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  org_id       uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  role_id      uuid        REFERENCES roles(id) ON DELETE SET NULL,
  display_name text,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- priority_classes
-- ============================================================

CREATE TABLE priority_classes (
  id             uuid                 PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id         uuid                 NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  level          priority_class_level NOT NULL,
  label          text                 NOT NULL,
  description    text,
  default_pinned boolean              NOT NULL DEFAULT false,
  created_at     timestamptz          NOT NULL DEFAULT now()
);

-- ============================================================
-- sites
-- ============================================================

CREATE TABLE sites (
  id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id      uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  name        text        NOT NULL,
  location    text,
  description text,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- nodes
-- ============================================================

CREATE TABLE nodes (
  id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id           uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  site_id          uuid        REFERENCES sites(id) ON DELETE CASCADE,
  name             text        NOT NULL,
  status           node_status NOT NULL DEFAULT 'UNKNOWN',
  cache_dir        text        NOT NULL,
  cache_max_bytes  bigint      NOT NULL,
  cache_used_bytes bigint      NOT NULL DEFAULT 0,
  last_seen        timestamptz,
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- node_heartbeats (range-partitioned by created_at)
-- ============================================================

CREATE TABLE node_heartbeats (
  id               uuid        NOT NULL DEFAULT gen_random_uuid(),
  org_id           uuid        NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  node_id          uuid        NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
  status           node_status NOT NULL,
  cache_used_bytes bigint      NOT NULL,
  cache_max_bytes  bigint      NOT NULL,
  agent_version    text,
  created_at       timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Monthly partitions for 2026
CREATE TABLE node_heartbeats_y2026m01 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE node_heartbeats_y2026m02 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE node_heartbeats_y2026m03 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE node_heartbeats_y2026m04 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE node_heartbeats_y2026m05 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE node_heartbeats_y2026m06 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE node_heartbeats_y2026m07 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE node_heartbeats_y2026m08 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE node_heartbeats_y2026m09 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE node_heartbeats_y2026m10 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE node_heartbeats_y2026m11 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE node_heartbeats_y2026m12 PARTITION OF node_heartbeats FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');
