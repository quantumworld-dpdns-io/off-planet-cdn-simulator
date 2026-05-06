package db

import (
	"context"
	"github.com/off-planet-cdn/control-api/internal/models"
)

// ListObjects returns cache objects for an org, with optional filters.
// siteID, priorityClassID, tag are optional — empty string means no filter.
func (c *Client) ListObjects(ctx context.Context, orgID, siteID, priorityClassID, tag string) ([]models.CacheObject, error) {
	query := `
        SELECT id, org_id, site_id, priority_class_id, name,
               COALESCE(content_type,''), COALESCE(source_url,''), COALESCE(content_hash,''),
               size_bytes, pinned, status::text, COALESCE(tags, '{}'), created_at, updated_at
        FROM cache_objects
        WHERE org_id = $1
          AND status != 'DELETED'
          AND ($2 = '' OR site_id::text = $2)
          AND ($3 = '' OR priority_class_id::text = $3)
          AND ($4 = '' OR $4 = ANY(tags))
        ORDER BY created_at DESC
        LIMIT 500`

	rows, err := c.Pool.Query(ctx, query, orgID, siteID, priorityClassID, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []models.CacheObject
	for rows.Next() {
		var o models.CacheObject
		if err := rows.Scan(
			&o.ID, &o.OrgID, &o.SiteID, &o.PriorityClassID, &o.Name,
			&o.ContentType, &o.SourceURL, &o.ContentHash,
			&o.SizeBytes, &o.Pinned, &o.Status, &o.Tags, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		objects = append(objects, o)
	}
	return objects, rows.Err()
}

func (c *Client) CreateObject(ctx context.Context, orgID, siteID, priorityClassID, name, contentType, sourceURL string, sizeBytes int64, tags []string) (*models.CacheObject, error) {
	if tags == nil {
		tags = []string{}
	}
	var o models.CacheObject
	err := c.Pool.QueryRow(ctx,
		`INSERT INTO cache_objects (org_id, site_id, priority_class_id, name, content_type, source_url, size_bytes, tags, status)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'ACTIVE')
         RETURNING id, org_id, site_id, priority_class_id, name,
                   COALESCE(content_type,''), COALESCE(source_url,''), COALESCE(content_hash,''),
                   size_bytes, pinned, status::text, COALESCE(tags,'{}'), created_at, updated_at`,
		orgID, siteID, priorityClassID, name, contentType, sourceURL, sizeBytes, tags,
	).Scan(
		&o.ID, &o.OrgID, &o.SiteID, &o.PriorityClassID, &o.Name,
		&o.ContentType, &o.SourceURL, &o.ContentHash,
		&o.SizeBytes, &o.Pinned, &o.Status, &o.Tags, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (c *Client) GetObject(ctx context.Context, orgID, objectID string) (*models.CacheObject, error) {
	var o models.CacheObject
	err := c.Pool.QueryRow(ctx,
		`SELECT id, org_id, site_id, priority_class_id, name,
                COALESCE(content_type,''), COALESCE(source_url,''), COALESCE(content_hash,''),
                size_bytes, pinned, status::text, COALESCE(tags,'{}'), created_at, updated_at
         FROM cache_objects WHERE id = $1 AND org_id = $2 AND status != 'DELETED'`,
		objectID, orgID,
	).Scan(
		&o.ID, &o.OrgID, &o.SiteID, &o.PriorityClassID, &o.Name,
		&o.ContentType, &o.SourceURL, &o.ContentHash,
		&o.SizeBytes, &o.Pinned, &o.Status, &o.Tags, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (c *Client) SetPinned(ctx context.Context, orgID, objectID string, pinned bool) (*models.CacheObject, error) {
	var o models.CacheObject
	err := c.Pool.QueryRow(ctx,
		`UPDATE cache_objects SET pinned = $1, updated_at = now()
         WHERE id = $2 AND org_id = $3
         RETURNING id, org_id, site_id, priority_class_id, name,
                   COALESCE(content_type,''), COALESCE(source_url,''), COALESCE(content_hash,''),
                   size_bytes, pinned, status::text, COALESCE(tags,'{}'), created_at, updated_at`,
		pinned, objectID, orgID,
	).Scan(
		&o.ID, &o.OrgID, &o.SiteID, &o.PriorityClassID, &o.Name,
		&o.ContentType, &o.SourceURL, &o.ContentHash,
		&o.SizeBytes, &o.Pinned, &o.Status, &o.Tags, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}
