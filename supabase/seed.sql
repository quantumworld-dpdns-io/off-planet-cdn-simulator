-- Seed data for Off-Planet CDN Simulator
-- All inserts use ON CONFLICT DO NOTHING for idempotency.

-- ============================================================
-- 1. Demo org
-- ============================================================

INSERT INTO orgs (id, name, slug)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  'Demo Organization',
  'demo'
)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 2. Priority classes
-- ============================================================

INSERT INTO priority_classes (id, org_id, level, label, description, default_pinned)
VALUES
  (
    '00000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000001',
    'P0', 'Emergency Medical',
    'Critical life-safety content — always pinned and never evicted.',
    true
  ),
  (
    '00000000-0000-0000-0000-000000000011',
    '00000000-0000-0000-0000-000000000001',
    'P1', 'Engineering Manuals',
    'Habitat operations and maintenance documentation.',
    false
  ),
  (
    '00000000-0000-0000-0000-000000000012',
    '00000000-0000-0000-0000-000000000001',
    'P2', 'Education',
    'Curriculum, research papers, and training materials.',
    false
  ),
  (
    '00000000-0000-0000-0000-000000000013',
    '00000000-0000-0000-0000-000000000001',
    'P3', 'Package Registries',
    'Software packages and dependency mirrors.',
    false
  ),
  (
    '00000000-0000-0000-0000-000000000014',
    '00000000-0000-0000-0000-000000000001',
    'P4', 'Entertainment',
    'Media, games, and recreational content.',
    false
  ),
  (
    '00000000-0000-0000-0000-000000000015',
    '00000000-0000-0000-0000-000000000001',
    'P5', 'Stale Content',
    'Expired or superseded content — eviction candidates.',
    false
  )
ON CONFLICT DO NOTHING;

-- ============================================================
-- 3. Roles
-- ============================================================

INSERT INTO roles (id, org_id, name)
VALUES
  ('00000000-0000-0000-0000-000000000020', '00000000-0000-0000-0000-000000000001', 'admin'),
  ('00000000-0000-0000-0000-000000000021', '00000000-0000-0000-0000-000000000001', 'operator'),
  ('00000000-0000-0000-0000-000000000022', '00000000-0000-0000-0000-000000000001', 'viewer'),
  ('00000000-0000-0000-0000-000000000023', '00000000-0000-0000-0000-000000000001', 'user')
ON CONFLICT DO NOTHING;

-- ============================================================
-- 4. Sites
-- ============================================================

INSERT INTO sites (id, org_id, name, location)
VALUES
  (
    '00000000-0000-0000-0000-000000000030',
    '00000000-0000-0000-0000-000000000001',
    'Lunar Habitat Alpha',
    'Sea of Tranquility, Moon'
  ),
  (
    '00000000-0000-0000-0000-000000000031',
    '00000000-0000-0000-0000-000000000001',
    'Mars Base Camp',
    'Jezero Crater, Mars'
  )
ON CONFLICT DO NOTHING;

-- ============================================================
-- 5. Nodes (one per site)
-- ============================================================

INSERT INTO nodes (id, org_id, site_id, name, status, cache_dir, cache_max_bytes)
VALUES
  (
    '00000000-0000-0000-0000-000000000040',
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000030',
    'primary-node',
    'UNKNOWN',
    '/var/cache/cdn',
    107374182400  -- 100 GB
  ),
  (
    '00000000-0000-0000-0000-000000000041',
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000031',
    'primary-node',
    'UNKNOWN',
    '/var/cache/cdn',
    107374182400  -- 100 GB
  )
ON CONFLICT DO NOTHING;

-- ============================================================
-- 6. Cache policies (one per site)
-- ============================================================

INSERT INTO cache_policies (id, org_id, site_id, name, description, enabled)
VALUES
  (
    '00000000-0000-0000-0000-000000000050',
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000030',
    'Default Mission Policy',
    'Baseline caching policy for Lunar Habitat Alpha.',
    true
  ),
  (
    '00000000-0000-0000-0000-000000000051',
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000031',
    'Default Mission Policy',
    'Baseline caching policy for Mars Base Camp.',
    true
  )
ON CONFLICT DO NOTHING;

-- ============================================================
-- 7. Bandwidth window (Lunar site)
--    2026-05-08 00:00 UTC -> 2026-05-08 02:00 UTC
--    1 Mbps = 1048576 bps, reliability 0.9
-- ============================================================

INSERT INTO bandwidth_windows (id, org_id, site_id, label, window_start, window_end, bandwidth_bps, reliability_score)
VALUES (
  '00000000-0000-0000-0000-000000000060',
  '00000000-0000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000030',
  'Earth-Moon Link Window',
  '2026-05-08 00:00:00+00',
  '2026-05-08 02:00:00+00',
  1048576,
  0.9
)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 8. Contact window (Mars site)
--    Daily at 12:00 UTC, 3600 s duration
-- ============================================================

INSERT INTO contact_windows (id, org_id, site_id, label, rrule, duration_seconds)
VALUES (
  '00000000-0000-0000-0000-000000000070',
  '00000000-0000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000031',
  'Daily DSN Contact',
  'FREQ=DAILY;BYHOUR=12;BYMINUTE=0',
  3600
)
ON CONFLICT DO NOTHING;
