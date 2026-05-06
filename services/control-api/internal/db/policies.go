package db

import (
	"context"

	"github.com/off-planet-cdn/control-api/internal/models"
)

func (c *Client) ListPolicies(ctx context.Context, orgID string) ([]models.CachePolicy, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, site_id, name, COALESCE(description,''), enabled, created_at, updated_at
         FROM cache_policies WHERE org_id = $1 ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var policies []models.CachePolicy
	for rows.Next() {
		var p models.CachePolicy
		if err := rows.Scan(&p.ID, &p.OrgID, &p.SiteID, &p.Name, &p.Description, &p.Enabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}
	return policies, rows.Err()
}

func (c *Client) CreatePolicy(ctx context.Context, orgID, siteID, name, description string, enabled bool) (*models.CachePolicy, error) {
	var p models.CachePolicy
	err := c.Pool.QueryRow(ctx,
		`INSERT INTO cache_policies (org_id, site_id, name, description, enabled)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, org_id, site_id, name, COALESCE(description,''), enabled, created_at, updated_at`,
		orgID, siteID, name, description, enabled,
	).Scan(&p.ID, &p.OrgID, &p.SiteID, &p.Name, &p.Description, &p.Enabled, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) UpdatePolicy(ctx context.Context, orgID, policyID, name, description string, enabled bool) (*models.CachePolicy, error) {
	var p models.CachePolicy
	err := c.Pool.QueryRow(ctx,
		`UPDATE cache_policies SET name = $1, description = $2, enabled = $3, updated_at = now()
         WHERE id = $4 AND org_id = $5
         RETURNING id, org_id, site_id, name, COALESCE(description,''), enabled, created_at, updated_at`,
		name, description, enabled, policyID, orgID,
	).Scan(&p.ID, &p.OrgID, &p.SiteID, &p.Name, &p.Description, &p.Enabled, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
