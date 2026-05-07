package db

import (
	"context"

	"github.com/off-planet-cdn/control-api/internal/models"
)

func (c *Client) ListMirrorSources(ctx context.Context, orgID string) ([]models.MirrorSource, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, registry_type, upstream_url, COALESCE(label,''), enabled, created_at, updated_at
         FROM mirror_sources WHERE org_id = $1 ORDER BY registry_type, created_at ASC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := []models.MirrorSource{}
	for rows.Next() {
		var s models.MirrorSource
		if err := rows.Scan(&s.ID, &s.OrgID, &s.RegistryType, &s.UpstreamURL, &s.Label, &s.Enabled, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

func (c *Client) CreateMirrorSource(ctx context.Context, orgID, registryType, upstreamURL, label string) (*models.MirrorSource, error) {
	var s models.MirrorSource
	err := c.Pool.QueryRow(ctx,
		`INSERT INTO mirror_sources (org_id, registry_type, upstream_url, label)
         VALUES ($1, $2, $3, NULLIF($4,''))
         RETURNING id, org_id, registry_type, upstream_url, COALESCE(label,''), enabled, created_at, updated_at`,
		orgID, registryType, upstreamURL, label,
	).Scan(&s.ID, &s.OrgID, &s.RegistryType, &s.UpstreamURL, &s.Label, &s.Enabled, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) ListMirrorArtifacts(ctx context.Context, orgID, sourceID string) ([]models.MirrorArtifact, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, source_id, name, version, COALESCE(size_bytes,0), COALESCE(storage_path,''), synced_at, created_at
         FROM mirror_artifacts
         WHERE org_id = $1 AND ($2 = '' OR source_id::text = $2)
         ORDER BY created_at DESC LIMIT 100`,
		orgID, sourceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artifacts := []models.MirrorArtifact{}
	for rows.Next() {
		var a models.MirrorArtifact
		if err := rows.Scan(&a.ID, &a.OrgID, &a.SourceID, &a.Name, &a.Version, &a.SizeBytes, &a.StoragePath, &a.SyncedAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		artifacts = append(artifacts, a)
	}
	return artifacts, rows.Err()
}
