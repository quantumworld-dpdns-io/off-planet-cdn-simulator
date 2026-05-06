package db

import (
	"context"

	"github.com/off-planet-cdn/control-api/internal/models"
)

func (c *Client) ListSites(ctx context.Context, orgID string) ([]models.Site, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, name, COALESCE(location,''), COALESCE(description,''), created_at, updated_at
         FROM sites WHERE org_id = $1 ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []models.Site
	for rows.Next() {
		var s models.Site
		if err := rows.Scan(&s.ID, &s.OrgID, &s.Name, &s.Location, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func (c *Client) CreateSite(ctx context.Context, orgID, name, location, description string) (*models.Site, error) {
	var s models.Site
	err := c.Pool.QueryRow(ctx,
		`INSERT INTO sites (org_id, name, location, description)
         VALUES ($1, $2, $3, $4)
         RETURNING id, org_id, name, COALESCE(location,''), COALESCE(description,''), created_at, updated_at`,
		orgID, name, location, description,
	).Scan(&s.ID, &s.OrgID, &s.Name, &s.Location, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) GetSite(ctx context.Context, orgID, siteID string) (*models.Site, error) {
	var s models.Site
	err := c.Pool.QueryRow(ctx,
		`SELECT id, org_id, name, COALESCE(location,''), COALESCE(description,''), created_at, updated_at
         FROM sites WHERE id = $1 AND org_id = $2`,
		siteID, orgID,
	).Scan(&s.ID, &s.OrgID, &s.Name, &s.Location, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
